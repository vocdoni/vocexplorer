package rpc

import (
	"github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// Ping pings the tendermint client and returns true if ok
func Ping(c *http.HTTP) bool {
	status, err := c.Status()
	if err != nil || status == nil {
		return false
	}
	return true
}

//UpdateTendermintInfo updates the tendermint info
func UpdateTendermintInfo(c *http.HTTP) {
	GetHealth(c)
	GetGenesis(c)
}

// GetHealth calls the tendermint Health api
func GetHealth(c *http.HTTP) {
	status, err := c.Status()
	if !util.ErrPrint(err) {
		store.Stats.ResultStatus = status
	}
}

// GetGenesis gets the first block
func GetGenesis(c *http.HTTP) {
	result, err := c.Genesis()
	if !util.ErrPrint(err) {
		store.Stats.Genesis = result.Genesis
	}
}

// GetBlock returns the contents of one block
func GetBlock(c *http.HTTP, height int64) *coretypes.ResultBlock {
	block, err := c.Block(&height)
	if util.ErrPrint(err) {
		return nil
	}
	return block
}

// GetTransaction gets a transaction by hash
func GetTransaction(c *http.HTTP, hash []byte) *coretypes.ResultTx {
	res, err := c.Tx(hash, false)
	if util.ErrPrint(err) {
		return nil
	}
	return res
}
