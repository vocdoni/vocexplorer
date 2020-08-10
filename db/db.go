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
	latestBlockHeight := int64(1)
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
	latestTxHeight := int64(1)
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

	// Wait for new blocks to be available
	for blockHeight-latestBlockHeight < 1 {
		time.Sleep(500 * time.Millisecond)
		status, err := c.Status()
		if err != nil {
			log.Error(err)
		}
		blockHeight = status.SyncInfo.LatestBlockHeight
	}

	batch := d.NewBatch()

	i := int64(0)
	numNewBlocks := util.Min(config.NumBlockUpdates, int(blockHeight-latestBlockHeight))
	// Array of new tx id's. Each goroutine can only access its assigned index, making this array thread-safe as long as all goroutines exit before read access
	txsList := make([]tmtypes.Txs, numNewBlocks)
	complete := make(chan struct{}, config.NumBlockUpdates)
	for ; int(i) < numNewBlocks; i++ {
		go fetchBlock(i+latestBlockHeight, &batch, c, cdc, complete, &txsList[i])
	}
	num := 0
	if i > 0 {
		// Sync: wait here for all goroutines to complete
		for range complete {
			if num >= numNewBlocks-1 {
				break
			}
			num++
		}
		log.Debugf("Setting block %d ", latestBlockHeight+i)

		//TODO: don't create a goroutine for empty tx lists
		complete = make(chan struct{}, len(txsList))
		for _, txs := range txsList {
			go updateTxs(latestTxHeight, txs, cdc, c, batch, complete)
			latestTxHeight += int64(len(txs))
		}

		// Sync: wait here for all goroutines to complete
		num = 0
		for range complete {
			if num >= len(txsList)-1 {
				break
			}
			num++
		}
		var txbuf bytes.Buffer
		err := amino.EncodeInt64(&txbuf, latestTxHeight)
		util.ErrPrint(err)
		var blockbuf bytes.Buffer
		err = amino.EncodeInt64(&blockbuf, latestBlockHeight+i)
		util.ErrPrint(err)
		batch.Put([]byte(config.LatestTxHeightKey), txbuf.Bytes())
		batch.Put([]byte(config.LatestBlockHeightKey), blockbuf.Bytes())
		batch.Write()
	}

}

func updateTxs(startTxHeight int64, txs tmtypes.Txs, cdc *amino.Codec, c *tmhttp.HTTP, batch dvotedb.Batch, complete chan<- struct{}) {
	numTxs := int64(0)
	var height int64
	for i, tx := range txs {
		numTxs = int64(i)
		txRes := rpc.GetTransaction(c, tx.Hash())

		txHashKey := append([]byte(config.TxHashPrefix), tx.Hash()...)
		txStore := types.StoreTx{
			Height:   txRes.Height,
			TxHeight: startTxHeight,
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
		if i == 0 {
			height = txRes.Height
		}
		startTxHeight++
	}
	if numTxs > 0 {
		log.Debugf("%d transactions logged at block %d, height %d", numTxs+1, height, startTxHeight)
	}
	complete <- struct{}{}
}

func fetchBlock(height int64, batch *dvotedb.Batch, c *tmhttp.HTTP, cdc *amino.Codec, complete chan<- struct{}, txs *tmtypes.Txs) {
	// Signal
	defer func() {
		complete <- struct{}{}
	}()
	// Thread-safe api request
	res, err := c.Block(&height)
	// If error is returned, try the request more times, then fatal.
	if err != nil {
		log.Error(err)
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

	*txs = res.Block.Data.Txs

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
