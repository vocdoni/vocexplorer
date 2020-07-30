package db

import (
	"bytes"
	"encoding/gob"
	"strconv"

	"gitlab.com/vocdoni/go-dvote/log"

	"strings"
	"time"

	"github.com/tendermint/tendermint/rpc/client/http"
	ttypes "github.com/tendermint/tendermint/types"
	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// NewDB initializes a badger database at the given path
func NewDB(path string) (*dvotedb.BadgerDB, error) {
	log.Infof("Initializing database")
	return dvotedb.NewBadgerDB(path)
}

// UpdateDB continuously updates the database by calling dvote & tendermint apis
func UpdateDB(d *dvotedb.BadgerDB, cfg *config.Cfg) {
	log.Infof("Updating database")
	// Init tendermint client
	tClient := rpc.StartClient(cfg.TendermintHost)
	// Init Gateway client
	gwClient, cancel := client.InitGateway(cfg.GatewayHost)
	if gwClient == nil || tClient == nil {
		log.Fatal("Cannot connect to blockchain clients")
	}
	log.Debugf("Connected")
	defer cancel()
	for {
		updateBlockList(d, tClient)
		updateEntityList(d)
		updateProcessList(d)
		time.Sleep(2 * time.Second)
	}
}

func updateBlockList(d *dvotedb.BadgerDB, c *http.HTTP) {
	latestHeight := int64(0)
	has, err := d.Has([]byte(config.LatestBlockHeightKey))
	if err != nil {
		log.Error(err)
	}
	if has {
		val, err := d.Get([]byte(config.LatestBlockHeightKey))
		if err != nil {
			log.Error(err)
		}
		latestHeight, err = strconv.ParseInt(string(val), 0, 64)
		if err != nil {
			log.Error(err)
		}
	}
	status, err := c.Status()
	if err != nil {
		log.Error(err)
	}
	blockHeight := status.SyncInfo.LatestBlockHeight
	batch := d.NewBatch()
	// encBuf := new(bytes.Buffer)
	// enc := gob.NewEncoder(encBuf)
	for i := 0; i < 2 && latestHeight < blockHeight; i++ {
		encBuf := new(bytes.Buffer)
		enc := gob.NewEncoder(encBuf)
		latestHeight++
		res, err := c.Block(&latestHeight)
		if err != nil {
			log.Error(err)
			break
		}
		key := append([]byte(config.BlockPrefix), []byte(util.IntToString(latestHeight))...)
		res.Block.LastCommit.Signatures = []ttypes.CommitSig{}
		res.Block.Evidence = ttypes.EvidenceData{}
		err = enc.Encode(res)
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("Getting block %d with key %s", latestHeight, string(key))

		value := encBuf.Bytes()
		batch.Put(key, value)
	}
	batch.Put([]byte(config.LatestBlockHeightKey), []byte(util.IntToString(latestHeight)))
	batch.Write()
}

func updateEntityList(d *dvotedb.BadgerDB) {

}

func updateProcessList(d *dvotedb.BadgerDB) {

}

// List returns a list of keys matching a given prefix
func list(d *dvotedb.BadgerDB, max int, from, prefix string) (list []string) {
	iter := d.NewIterator().(*dvotedb.BadgerIterator)
	if len(from) > 0 {
		iter.Seek([]byte(prefix + from))
	}
	for iter.Next() {
		if max < 1 {
			break
		}
		if strings.HasPrefix(string(iter.Key()), prefix) {
			log.Debugf("Found key %s", string(iter.Key()))
			list = append(list, string(iter.Key()))
			max--
		}
	}
	iter.Release()
	return
}
