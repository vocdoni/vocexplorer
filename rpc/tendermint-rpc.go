package rpc

import (
	"github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// TendermintInfo holds tendermint api results
type TendermintInfo struct {
	Status *coretypes.ResultStatus
}

// InitClient initializes an http tendermint api client on websockets
func InitClient() (*http.HTTP, error) {
	c, err := http.NewWithTimeout(config.TendermintHost, "/websocket", 1)
	if err != nil {
		return nil, err
	}
	_, err = c.Status()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// GetHealth calls the tendermint Health api
func (t *TendermintInfo) GetHealth(c *http.HTTP) {
	status, err := c.Status()
	if !util.ErrPrint(err) {
		t.Status = status
	}
}

//UpdateBlockInfo updates the block list
func UpdateBlockInfo(c *http.HTTP, t *TendermintInfo) {
	t.GetHealth(c)
}
