package db

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/tendermint/go-amino"
	"gitlab.com/vocdoni/go-dvote/log"

	"time"

	tmhttp "github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"
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
func UpdateDB(d *dvotedb.BadgerDB, gwHost, tmHost string) {
	ping := pingGateway(gwHost)
	if !ping {
		log.Warn("Gateway Client is not running. Running as detached database")
		return
	}

	// Init tendermint client
	tClient, up := startTendermint(tmHost)
	if !up {
		log.Warn("Cannot connect to tendermint client. Running as detached database")
		return
	}

	// Init Gateway client
	// gwClient, cancel, up := startGateway(cfg)
	// if !up {
	// 	log.Warn("Cannot connect to gateway client. Running as detached database")
	// 	return
	// }

	log.Debugf("Connected")
	// defer (*cancel)()

	// Init amino encoder
	var cdc = amino.NewCodec()
	cdc.RegisterConcrete(types.StoreBlock{}, "storeBlock", nil)
	cdc.RegisterConcrete(types.StoreTx{}, "storeTx", nil)
	for {
		updateBlockList(d, tClient, cdc)
		updateEntityList(d)
		updateProcessList(d)
		time.Sleep(config.DBWaitTime * time.Millisecond)
	}
}

func updateBlockList(d *dvotedb.BadgerDB, c *tmhttp.HTTP, cdc *amino.Codec) {
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
	errs := 0
	for i := 0; i < config.NumBlockUpdates && latestHeight < blockHeight; i++ {
		// encBuf := new(bytes.Buffer)
		// enc := gob.NewEncoder(encBuf)
		latestHeight++
		res, err := c.Block(&latestHeight)
		if err != nil {
			if errs > 2 {
				log.Fatal("Blockchain RPC Disconnected")
				return
			}
			errs++
			log.Error(err)
			i--
			continue
		}
		errs = 0
		var block types.StoreBlock
		block.NumTxs = len(res.Block.Data.Txs)
		block.Hash = res.BlockID.Hash
		block.Height = res.Block.Header.Height
		block.Time = res.Block.Header.Time

		updateTxs(res.Block, d, cdc, c, batch)

		bodyValue, err := cdc.MarshalBinaryLengthPrefixed(block)
		if err != nil {
			log.Error(err)
		}

		blockHeightKey := append([]byte(config.BlockHeightPrefix), []byte(util.IntToString(block.Height))...)
		blockHashKey := append([]byte(config.BlockHashPrefix), block.Hash...)
		hashValue := block.Hash

		batch.Put(blockHashKey, bodyValue)
		batch.Put(blockHeightKey, hashValue)

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

func updateTxs(block *tmtypes.Block, d *dvotedb.BadgerDB, cdc *amino.Codec, c *tmhttp.HTTP, batch dvotedb.Batch) {
	currentTxs := int64(0)
	numTxs := int64(0)
	has, err := d.Has([]byte(config.LatestTxHeightKey))
	if err != nil {
		log.Error(err)
	}
	if has {
		val, err := d.Get([]byte(config.LatestTxHeightKey))
		util.ErrPrint(err)
		num := 0
		currentTxs, num, err = amino.DecodeInt64(val)
		util.ErrPrint(err)
		if num <= 1 {
			log.Debug("Could not get height")
		}
	}
	for i, tx := range block.Txs {
		numTxs = int64(i) + 1
		txRes := rpc.GetTransaction(c, tx.Hash())

		txHashKey := append([]byte(config.TxHashPrefix), tx.Hash()...)
		txStore := types.StoreTx{
			Height:   txRes.Height,
			TxHeight: currentTxs + numTxs,
			Tx:       txRes.Tx,
			TxResult: txRes.TxResult,
			Index:    txRes.Index,
		}
		txVal, err := cdc.MarshalBinaryLengthPrefixed(txStore)
		if err != nil {
			log.Error(err)
		}
		batch.Put(txHashKey, txVal)
		txHeightKey := []byte(config.TxHeightPrefix + util.IntToString(txStore.TxHeight))
		batch.Put(txHeightKey, tx.Hash())
	}
	if numTxs > 0 {
		var buf bytes.Buffer
		err = amino.EncodeInt64(&buf, currentTxs+numTxs)
		batch.Put([]byte(config.LatestTxHeightKey), buf.Bytes())
		log.Debugf("%d transactions logged at block %d", numTxs+1, block.Height)
	}
}

func updateEntityList(d *dvotedb.BadgerDB) {

}

func updateProcessList(d *dvotedb.BadgerDB) {

}

// listKeysByHeight returns a list of hashes matching a given prefix, where key is an integer
func listKeysByHeight(d *dvotedb.BadgerDB, max, from int, prefix string) (list []string) {
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
		keyList = append(keyList, string(key))
		from++
	}
	return keyList
}

// listHashesByHeight returns a list of hashes given integer keys
func listHashesByHeight(d *dvotedb.BadgerDB, max, height int, prefix string) [][]byte {
	if max > 64 {
		max = 64
	}
	var hashList [][]byte
	for ; max > 0; max-- {
		key := []byte(prefix + util.IntToString(height))
		has, err := d.Has(key)
		if !has || util.ErrPrint(err) {
			break
		}
		val, err := d.Get(key)
		if err != nil {
			log.Error(err)
		}
		hashList = append(hashList, val)
		height--
	}
	log.Debugf("Found %d hashes", len(hashList))
	return hashList
}

func pingGateway(host string) bool {
	pingClient := http.Client{
		Timeout: 5 * time.Second,
	}
	for i := 0; ; i++ {
		if i > 10 {
			return false
		}
		resp, err := pingClient.Get("http://" + host + "/ping")
		if err != nil {
			log.Debug(err.Error())
			time.Sleep(2 * time.Second)
			continue
		}
		body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
		if err != nil {
			log.Debug(err.Error())
			time.Sleep(time.Second)
			continue
		}
		if string(body) != "pong" {
			log.Warn("Gateway ping not yet available")
		} else {
			return true
		}
	}
}

func startGateway(host string) (*client.Client, *context.CancelFunc, bool) {
	for i := 0; ; i++ {
		if i > 20 {
			return nil, nil, false
		}
		gwClient, cancel := client.InitGateway(host)
		if gwClient == nil {
			time.Sleep(5 * time.Second)
			continue
		} else {
			return gwClient, &cancel, true
		}
	}
}

func startTendermint(host string) (*tmhttp.HTTP, bool) {
	for i := 0; ; i++ {
		if i > 20 {
			return nil, false
		}
		tmClient := rpc.StartClient("http://" + host)
		if tmClient == nil {
			time.Sleep(5 * time.Second)
			continue
		} else {
			return tmClient, true
		}
	}
}
