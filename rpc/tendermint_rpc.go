package rpc

import (
	"github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
)

// Ping pings the tendermint client and returns true if ok
func Ping(c *http.HTTP) bool {
	status, err := c.Status()
	if err != nil || status == nil {
		return false
	}
	return true
}

//UpdateBlockchainStatus updates the tendermint info
func UpdateBlockchainStatus(c *http.HTTP) {
	GetHealth(c)
	GetGenesis(c)
}

// GetHealth calls the tendermint Health api
func GetHealth(c *http.HTTP) {
	status, err := c.Status()
	if err != nil {
		log.Error(err)
	} else {
		dispatcher.Dispatch(&actions.SetResultStatus{Status: status})
	}
}

// GetGenesis gets the first block
func GetGenesis(c *http.HTTP) {
	result, err := c.Genesis()
	if err != nil {
		log.Error(err)
	} else {
		dispatcher.Dispatch(&actions.SetGenesis{Genesis: result.Genesis})
	}
}

// GetBlock returns the contents of one block
func GetBlock(c *http.HTTP, height int64) *coretypes.ResultBlock {
	block, err := c.Block(&height)
	if err != nil {
		log.Error(err)
		return nil
	}
	return block
}

// GetTransaction gets a transaction by hash
func GetTransaction(c *http.HTTP, hash []byte) *coretypes.ResultTx {
	res, err := c.Tx(hash, false)
	if err != nil {
		log.Error(err)
		return nil
	}
	return res
}
