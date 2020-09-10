package api

import (
	gohttp "net/http"
	"time"

	"github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	jsonrpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	"github.com/tendermint/tendermint/types"
	"gitlab.com/vocdoni/go-dvote/log"
)

// StartTendermintClient initializes an http tendermint api client on websockets
func StartTendermintClient(host string) *http.HTTP {
	log.Infof("connecting to %s", host)
	tClient, err := initClient(host)
	if err != nil {
		log.Warn(err.Error())
		return nil
	}
	return tClient
}

func initClient(host string) (*http.HTTP, error) {
	httpClient, err := jsonrpcclient.DefaultHTTPClient(host)
	if err != nil {
		return nil, err
	}
	httpClient.Timeout = 2 * time.Second
	// Increase max idle connections. This fixes issue with too many concurrent requests, as described here: https://github.com/golang/go/issues/16012
	httpClient.Transport.(*gohttp.Transport).MaxIdleConnsPerHost = 10000
	c, err := http.NewWithClient(host, "/websocket", httpClient)
	// c, err := http.NewWithTimeout(host, "/websocket", 2)
	if err != nil {
		return nil, err
	}
	_, err = c.Status()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// PingTendermint pings the tendermint client and returns true if ok
func PingTendermint(c *http.HTTP) bool {
	if c == nil {
		return false
	}
	status, err := c.Status()
	if err != nil || status == nil {
		return false
	}
	return true
}

// GetHealth calls the tendermint Health api
func GetHealth(c *http.HTTP) *coretypes.ResultStatus {
	if c == nil {
		return nil
	}
	status, err := c.Status()
	if err != nil {
		log.Error(err)
		return nil
	}
	return status
}

// GetGenesis gets the first block
func GetGenesis(c *http.HTTP) *types.GenesisDoc {
	if c == nil {
		return nil
	}
	result, err := c.Genesis()
	if err != nil {
		log.Error(err)
		return nil
	}
	return result.Genesis
}

// GetBlock returns the contents of one block
func GetBlock(c *http.HTTP, height int64) *coretypes.ResultBlock {
	if c == nil {
		return nil
	}
	block, err := c.Block(&height)
	if err != nil {
		log.Error(err)
		return nil
	}
	return block
}

// GetTransaction gets a transaction by hash
func GetTransaction(c *http.HTTP, hash []byte) *coretypes.ResultTx {
	if c == nil {
		return nil
	}
	res, err := c.Tx(hash, false)
	if err != nil {
		log.Error(err)
		return nil
	}
	return res
}
