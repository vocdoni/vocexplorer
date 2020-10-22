package db

import (
	"os"
	"os/signal"
	"runtime/pprof"
	"time"

	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	voctypes "gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/vochain"
	"google.golang.org/protobuf/proto"
)

type ExplorerDB struct {
	Db *dvotedb.BadgerDB
	Vs *vochain.VochainService
}

// NewDB initializes a badger database at the given path
func NewDB(cfg *config.MainCfg) *ExplorerDB {
	log.Infof("Initializing database at " + cfg.DataDir + "/" + cfg.Chain)
	var err error
	db := new(ExplorerDB)
	db.Vs, err = vochain.InitVochain(cfg)
	if err != nil {
		log.Fatal(err)
	}
	db.Db, err = dvotedb.NewBadgerDB(cfg.DataDir + "/" + db.Vs.GetStatus().ChainID)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Close closes the explorer db
func (d *ExplorerDB) Close() {
	d.Vs.Close()
}

// UpdateDB continuously updates the database by calling dvote & tendermint apis
func (d *ExplorerDB) UpdateDB() {
	// Init height keys
	batch := d.Db.NewBatch()
	zeroHeight := voctypes.Height{Height: 0}
	encHeight, err := proto.Marshal(&zeroHeight)
	if err != nil {
		log.Error(err)
	}
	if ok, err := d.Db.Has([]byte(config.LatestTxHeightKey)); !ok || err != nil {
		batch.Put([]byte(config.LatestTxHeightKey), encHeight)
	}
	if ok, err := d.Db.Has([]byte(config.LatestBlockHeightKey)); !ok || err != nil {
		batch.Put([]byte(config.LatestBlockHeightKey), encHeight)
	}
	if err != nil {
		log.Error(err)
	}
	if ok, err := d.Db.Has([]byte(config.LatestEntityCountKey)); !ok || err != nil {
		batch.Put([]byte(config.LatestEntityCountKey), encHeight)
	}
	if ok, err := d.Db.Has([]byte(config.LatestEnvelopeCountKey)); !ok || err != nil {
		batch.Put([]byte(config.LatestEnvelopeCountKey), encHeight)
	}
	if ok, err := d.Db.Has([]byte(config.LatestProcessCountKey)); !ok || err != nil {
		batch.Put([]byte(config.LatestProcessCountKey), encHeight)
	}
	if ok, err := d.Db.Has([]byte(config.LatestValidatorCountKey)); !ok || err != nil {
		batch.Put([]byte(config.LatestValidatorCountKey), encHeight)
	}
	batch.Write()

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
			pprof.StopCPUProfile()
			d.Close()
			os.Exit(1)
		}
	}()

	waitTime := 0
	i := 0
	for {
		// If synced, wait.
		d.waitSync(&waitTime)
		d.updateBlockchainInfo()
		d.updateBlockList()
		// Update validators less frequently than blocks
		if i%40 == 0 {
			d.updateValidatorList()
		}
		d.updateEntityList()
		d.updateProcessList()
		time.Sleep(time.Duration(waitTime) * time.Millisecond)
		i++
	}
}

func (d *ExplorerDB) waitSync(waitTime *int) {
	sync := d.isSynced()
	for sync {
		*waitTime = config.DBWaitTime
		time.Sleep(1 * time.Second)
		sync = d.isSynced()
	}
}

func (d *ExplorerDB) isSynced() bool {
	localBlockHeight := GetHeight(d.Db, config.LatestBlockHeightKey, 1)
	globalBlockHeight := d.Vs.GetBlockHeight()
	return globalBlockHeight-localBlockHeight.GetHeight() < 1
}
