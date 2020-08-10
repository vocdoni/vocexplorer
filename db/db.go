package db

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"

	amino "github.com/tendermint/go-amino"
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
	latestBlockHeight := int64(0)
	has, err := d.Has([]byte(config.LatestBlockHeightKey))
	if err != nil {
		log.Error(err)
	}
	if has {
		val, err := d.Get([]byte(config.LatestBlockHeightKey))
		if err != nil {
			log.Error(err)
		}
		latestBlockHeight, _, err = amino.DecodeInt64(val)
		if err != nil {
			log.Error(err)
		}
	}
	latestTxHeight := int64(0)
	has, err = d.Has([]byte(config.LatestTxHeightKey))
	if err != nil {
		log.Error(err)
	}
	if has {
		val, err := d.Get([]byte(config.LatestTxHeightKey))
		util.ErrPrint(err)
		num := 0
		latestTxHeight, num, err = amino.DecodeInt64(val)
		util.ErrPrint(err)
		if num <= 1 {
			log.Debug("Could not get height")
		}
	}
	status, err := c.Status()
	if err != nil {
		log.Error(err)
	}
	blockHeight := status.SyncInfo.LatestBlockHeight
	batch := d.NewBatch()

	i := int64(0)
	complete := make(chan struct{}, config.NumBlockUpdates)
	for ; i < config.NumBlockUpdates && i+latestBlockHeight < blockHeight; i++ {
		go fetchBlock(i+latestBlockHeight, &latestTxHeight, &batch, c, cdc, complete)
	}
	num := 0
	// Sync: wait here for all goroutines to complete
	for range complete {
		if num >= config.NumBlockUpdates-1 {
			break
		}
		num++
	}

	log.Debugf("Setting block %d ", latestBlockHeight+i)
	var buf bytes.Buffer
	err = amino.EncodeInt64(&buf, latestBlockHeight+i)
	if err != nil {
		log.Error(err)
	}
	batch.Put([]byte(config.LatestBlockHeightKey), buf.Bytes())
	batch.Write()
}

func updateTxs(startTxHeight *int64, block *tmtypes.Block, d *dvotedb.BadgerDB, cdc *amino.Codec, c *tmhttp.HTTP, batch dvotedb.Batch) {
	numTxs := int64(0)
	for i, tx := range block.Txs {
		numTxs = int64(i) + 1
		txRes := rpc.GetTransaction(c, tx.Hash())

		txHashKey := append([]byte(config.TxHashPrefix), tx.Hash()...)
		txStore := types.StoreTx{
			Height:   txRes.Height,
			TxHeight: *startTxHeight + numTxs,
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
		log.Debugf("Log tx %d", txStore.TxHeight)
	}
	if numTxs > 0 {
		var buf bytes.Buffer
		err := amino.EncodeInt64(&buf, *startTxHeight+numTxs)
		util.ErrPrint(err)
		batch.Put([]byte(config.LatestTxHeightKey), buf.Bytes())
		*startTxHeight += numTxs
		log.Debugf("%d transactions logged at block %d", numTxs+1, block.Height)
	}
}

func fetchBlock(height int64, latestTxHeight *int64, batch *dvotedb.Batch, c *tmhttp.HTTP, cdc *amino.Codec, complete chan<- struct{}) {
	// Signal
	defer func() {
		complete <- struct{}{}
	}()
	// Thread-safe api request
	res, err := c.Block(&height)
	// If error is returned, try the request more times, then fatal.
	if err != nil {
		for errs := 0; ; errs++ {
			if errs > 10 {
				log.Fatal("Blockchain RPC Disconnected")
				return
			}
			res, err = c.Block(&height)
			if err == nil {
				break
			}
		}
	}
	var block types.StoreBlock
	block.NumTxs = len(res.Block.Data.Txs)
	block.Hash = res.BlockID.Hash
	block.Height = res.Block.Header.Height
	block.Time = res.Block.Header.Time

	// updateTxs(&latestTxHeight, res.Block, d, cdc, c, batch)

	bodyValue, err := cdc.MarshalBinaryLengthPrefixed(block)
	if err != nil {
		log.Error(err)
	}

	blockHeightKey := append([]byte(config.BlockHeightPrefix), []byte(util.IntToString(block.Height))...)
	blockHashKey := append([]byte(config.BlockHashPrefix), block.Hash...)
	hashValue := block.Hash

	// Thread-safe batch operations
	(*batch).Put(blockHashKey, bodyValue)
	(*batch).Put(blockHeightKey, hashValue)
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
			time.Sleep(1 * time.Second)
			continue
		} else {
			return tmClient, true
		}
	}
}
