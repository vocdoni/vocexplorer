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
	for {
		updateBlockList(d, tClient, cdc)
		updateEntityList(d)
		updateProcessList(d)
		time.Sleep(config.DBWaitTime * time.Second)
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
		keyList = append(keyList, string(key))
		from++
	}
	return keyList
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
