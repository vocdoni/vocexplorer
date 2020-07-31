package db

import (
	"bytes"

	"github.com/tendermint/go-amino"
	"gitlab.com/vocdoni/go-dvote/log"

	"time"

	"github.com/tendermint/tendermint/rpc/client/http"
	dvotedb "gitlab.com/vocdoni/go-dvote/db"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/types"
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

	// Init amino encoder
	var cdc = amino.NewCodec()
	cdc.RegisterConcrete(types.StoreBlock{}, "storeBlock", nil)
	for {
		updateBlockList(d, tClient, cdc)
		updateEntityList(d)
		updateProcessList(d)
		time.Sleep(2 * time.Second)
	}
}

func updateBlockList(d *dvotedb.BadgerDB, c *http.HTTP, cdc *amino.Codec) {
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
		latestHeight, _, err = amino.DecodeInt64(val)
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
	for i := 0; i < 2 && latestHeight < blockHeight; i++ {
		// encBuf := new(bytes.Buffer)
		// enc := gob.NewEncoder(encBuf)
		latestHeight++
		res, err := c.Block(&latestHeight)
		if err != nil {
			log.Error(err)
			break
		}
		var block types.StoreBlock
		block.Data = res.Block.Data
		block.Hash = res.BlockID.Hash
		block.Height = res.Block.Header.Height
		block.Time = res.Block.Header.Time
		// err = enc.Encode(block)

		value, err := cdc.MarshalBinaryLengthPrefixed(block)
		if err != nil {
			log.Fatal(err)
		}

		key := append([]byte(config.BlockPrefix), []byte(util.IntToString(block.Height))...)
		// value := encBuf.Bytes()
		batch.Put(key, value)

		key = append([]byte(config.BlockHashPrefix), block.Hash...)
		value = []byte(util.IntToString(block.Height))
		batch.Put(key, value)

	}
	log.Debugf("Setting block %d ", latestHeight)
	var buf bytes.Buffer
	err = amino.EncodeInt64(&buf, latestHeight)
	if err != nil {
		log.Error(err)
	}
	batch.Put([]byte(config.LatestBlockHeightKey), buf.Bytes())
	batch.Write()
}

func updateEntityList(d *dvotedb.BadgerDB) {

}

func updateProcessList(d *dvotedb.BadgerDB) {

}

// List returns a list of keys matching a given prefix
func list(d *dvotedb.BadgerDB, max, from int, prefix string) (list []string) {
	if max > 64 {
		max = 64
	}
	var keyList []string
	for i := 0; i < max; i++ {
		key := prefix + util.IntToString(from)
		has, err := d.Has([]byte(key))
		if err != nil {
			log.Error(err)
			break
		}
		if !has {
			break
		}
		log.Debugf("Getting key %s", key)
		keyList = append(keyList, string(key))
		from++
	}
	return keyList
}

// func list(d *dvotedb.BadgerDB, max int, from, prefix string) (list []string) {
// 	iter := d.NewIterator().(*dvotedb.BadgerIterator)
// 	if len(from) > 0 {
// 		// height, err := strconv.Atoi(from)
// 		// if err != nil {
// 		// 	log.Error(err)
// 		// }
// 		// iter.Seek(append([]byte(prefix), []byte(height)...))
// 		log.Debugf("Searching for key %s", prefix+from)
// 		iter.Seek(append([]byte(prefix), []byte(from)...))
// 	}
// 	var keyList []string
// 	for iter.Next() {
// 		if max < 1 {
// 			break
// 		}
// 		if strings.HasPrefix(string(iter.Key()), prefix) {
// 			log.Debugf("Found key %s", string(iter.Key()))
// 			keyList = append(keyList, string(iter.Key()))
// 			max--
// 		}
// 	}
// 	iter.Release()
// 	return keyList
// }
