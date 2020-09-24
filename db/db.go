package db

import (
	"errors"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"time"

	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/api/rpc"
	"gitlab.com/vocdoni/vocexplorer/config"
	voctypes "gitlab.com/vocdoni/vocexplorer/proto"
	"google.golang.org/protobuf/proto"
)

var exit chan struct{}

// NewDB initializes a badger database at the given path
func NewDB(path, chainID string) (*dvotedb.BadgerDB, error) {
	if chainID == "" {
		return nil, errors.New("chain ID empty, cannot initialize database. See --chainID config option if running in detached mode")
	}
	log.Infof("Initializing database at " + path + "/" + chainID)
	return dvotedb.NewBadgerDB(path + "/" + chainID)
}

// UpdateDB continuously updates the database by calling dvote & tendermint apis
func UpdateDB(d *dvotedb.BadgerDB, detached *bool, tmHost, gwHost string) {
	exit = make(chan struct{}, 100)

	// Init height keys
	batch := d.NewBatch()
	zeroHeight := voctypes.Height{Height: 0}
	encHeight, err := proto.Marshal(&zeroHeight)
	if err != nil {
		log.Error(err)
	}
	if ok, err := d.Has([]byte(config.LatestTxHeightKey)); !ok || err != nil {
		batch.Put([]byte(config.LatestTxHeightKey), encHeight)
	}
	if ok, err := d.Has([]byte(config.LatestBlockHeightKey)); !ok || err != nil {
		batch.Put([]byte(config.LatestBlockHeightKey), encHeight)
	}
	if err != nil {
		log.Error(err)
	}
	if ok, err := d.Has([]byte(config.LatestEntityCountKey)); !ok || err != nil {
		batch.Put([]byte(config.LatestEntityCountKey), encHeight)
	}
	if ok, err := d.Has([]byte(config.LatestEnvelopeCountKey)); !ok || err != nil {
		batch.Put([]byte(config.LatestEnvelopeCountKey), encHeight)
	}
	if ok, err := d.Has([]byte(config.LatestProcessCountKey)); !ok || err != nil {
		batch.Put([]byte(config.LatestProcessCountKey), encHeight)
	}
	if ok, err := d.Has([]byte(config.LatestValidatorCountKey)); !ok || err != nil {
		batch.Put([]byte(config.LatestValidatorCountKey), encHeight)
	}
	batch.Write()

	// Init tendermint client
	tClient, ok := api.StartTendermint(tmHost, 20)
	if !ok {
		log.Warn("Cannot connect to tendermint api. Running as detached database")
		return
	}
	log.Debugf("Connected to " + tmHost)

	// Init gateway client
	gwClient, cancel := api.InitGateway(gwHost)
	if gwClient == nil {
		log.Warn("Cannot connect to gateway api. Running as detached database")
		*detached = true
		return
	}
	defer (cancel)()
	log.Debugf("Connected to %s", gwHost)

	// Interrupt signal should close clients
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Infof("captured %v, stopping profiler and closing websockets connections...", sig)
			go func() {
				time.Sleep(30 * time.Second)
				os.Exit(1)
			}()
			if tClient != nil {
				tClient.Close()
			}
			gwClient.Close()
			pprof.StopCPUProfile()
			os.Exit(1)
		}
	}()

	i := 0
	hasSynced := false
	for {
		select {
		case <-exit:
			*detached = true
			log.Warnf("Gateway disconnected, converting to detached mode")
			return
		default:
			// If synced, wait. If first time synced, reduce connections to 2.
			waitSync(d, tClient, &hasSynced, tmHost)
			updateBlockList(d, tClient)
			// Update validators less frequently than blocks
			if i%40 == 0 {
				updateValidatorList(d, tClient)
			}
			updateEntityList(d, gwClient)
			updateProcessList(d, gwClient)
			time.Sleep(config.DBWaitTime * time.Millisecond)
			i++
		}
	}
}

func waitSync(d *dvotedb.BadgerDB, t *rpc.TendermintRPC, hasSynced *bool, host string) {
	sync := isSynced(d, t)
	for sync {
		if !*hasSynced {
			*hasSynced = true
			oldClient := *t
			newClient, ok := api.StartTendermint(host, 2)
			if !ok {
				log.Warn("Cannot connect to tendermint api. Running as detached database")
				return
			}
			*t = *newClient
			log.Infof("Blockchain storage is synced, reducing to 2 websockets connections")
			oldClient.Close()
		}
		time.Sleep(1 * time.Second)
		sync = isSynced(d, t)
	}
}

func isSynced(d *dvotedb.BadgerDB, t *rpc.TendermintRPC) bool {
	latestBlockHeight := GetHeight(d, config.LatestBlockHeightKey, 1)
	status, err := t.Status()
	// If error is returned, try the request more times, then fatal.
	if err != nil {
		if strings.Contains(err.Error(), "WebSocket closed") {
			exit <- struct{}{}
			return false
		}
		log.Error(err)
		for errs := 0; ; errs++ {
			if errs > 10 {
				log.Error("Gateway Disconnected")
				exit <- struct{}{}
				return false
			}
			status, err = t.Status()
			if err == nil {
				break
			}
		}
	}
	gwBlockHeight := status.SyncInfo.LatestBlockHeight
	return gwBlockHeight-latestBlockHeight.GetHeight() < 1
}
