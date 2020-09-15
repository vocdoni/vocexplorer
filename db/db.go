package db

import (
	"errors"
	"time"

	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/config"
	voctypes "gitlab.com/vocdoni/vocexplorer/proto"
	"google.golang.org/protobuf/proto"
)

var exit chan struct{}

// NewDB initializes a badger database at the given path
func NewDB(path, chainID string) (*dvotedb.BadgerDB, error) {
	if chainID == "" {
		return nil, errors.New("Chain ID empty, cannot initialize database. See --chainID config option if running in detached mode")
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
	tClient, ok := api.StartTendermint(tmHost, 100)
	if !ok {
		log.Warn("Cannot connect to tendermint api. Running as detached database")
		return
	}
	log.Debugf("Connected to " + tmHost)
	defer func() { tClient.Close() }()

	// Init gateway client
	gwClient, cancel, up := api.StartGateway(gwHost)
	if !up {
		log.Warn("Cannot connect to gateway api. Running as detached database")
		*detached = true
		return
	}
	defer func() { gwClient.Close() }()
	defer (*cancel)()
	log.Debugf("Connected to %s", gwHost)

	i := 0
	for {
		select {
		case <-exit:
			*detached = true
			log.Warnf("Gateway disconnected, converting to detached mode")
			return
		default:
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
