package rpc

import (
	"fmt"

	"github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// TendermintInfo holds tendermint api results
type TendermintInfo struct {
	ResultStatus *coretypes.ResultStatus
	Genesis      *tmtypes.GenesisDoc
	BlockList    [config.ListSize]types.StoreBlock
	// BlockList        [config.ListSize]coretypes.ResultBlock
	// BlockListResults [config.ListSize]coretypes.ResultBlockResults
	TxCount       int
	TxList        []*coretypes.ResultTx
	ChainID       int
	AppVersion    int
	MaxBlockSize  int
	NumValidators int
	TotalBlocks   int
}

// StartClient initializes an http tendermint api client on websockets
func StartClient(host string) *http.HTTP {
	fmt.Println("connecting to " + host)
	tClient, err := initClient(host)
	if util.ErrPrint(err) {
		return nil
	}
	return tClient
}

func initClient(host string) (*http.HTTP, error) {
	c, err := http.NewWithTimeout(host, "/websocket", 2)
	if err != nil {
		return nil, err
	}
	_, err = c.Status()
	if err != nil {
		return nil, err
	}
	return c, nil
}

//UpdateTendermintInfo updates the tendermint info
func UpdateTendermintInfo(c *http.HTTP, t *TendermintInfo, i int) {
	t.GetHealth(c)
	t.GetGenesis(c)
	// UpdateBlockList(c, t, i)
}

// GetHealth calls the tendermint Health api
func (t *TendermintInfo) GetHealth(c *http.HTTP) {
	status, err := c.Status()
	if !util.ErrPrint(err) {
		t.ResultStatus = status
	}
}

// GetGenesis gets the first block
func (t *TendermintInfo) GetGenesis(c *http.HTTP) {
	result, err := c.Genesis()
	if !util.ErrPrint(err) {
		t.Genesis = result.Genesis
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
