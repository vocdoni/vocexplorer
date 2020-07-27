package rpc

import (
	"fmt"
	"syscall/js"

	"github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/rpc/client/http"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// TendermintInfo holds tendermint api results
type TendermintInfo struct {
	ResultStatus       *coretypes.ResultStatus
	Genesis            *types.GenesisDoc
	RecentBlocks       []coretypes.ResultBlock
	RecentBlockResults []coretypes.ResultBlockResults
	TxCount            int
	TxList             []*coretypes.ResultTx
	ChainID            int
	AppVersion         int
	MaxBlockSize       int
	NumValidators      int
}

// StartClient initializes an http tendermint api client on websockets
func StartClient() *http.HTTP {
	fmt.Println("connecting to %s", config.TendermintHost)
	tClient, err := initClient()
	if util.ErrPrint(err) {
		if js.Global().Get("confirm").Invoke("Unable to connect to Tendermint client. Reload with client running").Bool() {
			js.Global().Get("location").Call("reload")
		}
		return nil
	}
	return tClient
}

func initClient() (*http.HTTP, error) {
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

//UpdateTendermintInfo updates the tendermint info
func UpdateTendermintInfo(c *http.HTTP, t *TendermintInfo) {
	t.GetHealth(c)
	t.GetAllTxs(c)
	t.GetGenesis(c)
	t.GetRecentBlocks(c)
	t.GetRecentBlockResults(c)
}

// GetHealth calls the tendermint Health api
func (t *TendermintInfo) GetHealth(c *http.HTTP) {
	status, err := c.Status()
	if !util.ErrPrint(err) {
		t.ResultStatus = status
	}
}

// GetAllTxs gets all txs
func (t *TendermintInfo) GetAllTxs(c *http.HTTP) {
	fromHeight := t.TxCount
	query, err := query.New("tx.height>=" + util.IntToString(fromHeight))
	util.ErrPrint(err)
	txs, err := c.TxSearch(query.String(), false, 1, 100, "asc")
	if !util.ErrPrint(err) {
		t.TxCount += txs.TotalCount
		t.TxList = append(t.TxList, txs.Txs...)
	}
}

// GetGenesis gets the first block
func (t *TendermintInfo) GetGenesis(c *http.HTTP) {
	result, err := c.Genesis()
	if !util.ErrPrint(err) {
		t.Genesis = result.Genesis
	}
}

// GetRecentBlocks keeps a list of the four most recent blocks
func (t *TendermintInfo) GetRecentBlocks(c *http.HTTP) {
	var lastBlockHeight int64
	if t.RecentBlocks != nil && len(t.RecentBlocks) > 0 {
		lastBlockHeight = t.RecentBlocks[len(t.RecentBlocks)-1].Block.Header.Height
	}
	// Number of new blocks not already stored
	numNew := util.Min(int(t.ResultStatus.SyncInfo.LatestBlockHeight-lastBlockHeight)-1, 4)
	if numNew <= len(t.RecentBlocks) {
		t.RecentBlocks = t.RecentBlocks[numNew:]
	} else {
		t.RecentBlocks = []coretypes.ResultBlock{}
	}
	for numNew > 0 {
		nextHeight := t.ResultStatus.SyncInfo.LatestBlockHeight - int64(numNew)
		result, err := c.Block(&nextHeight)
		if !util.ErrPrint(err) {
			t.RecentBlocks = append(t.RecentBlocks, *result)
		}
		numNew--
	}
}

// GetRecentBlockResults keeps a list of the four most recent blocks
func (t *TendermintInfo) GetRecentBlockResults(c *http.HTTP) {
	var lastBlockHeight int64
	if t.RecentBlockResults != nil && len(t.RecentBlockResults) > 0 {
		lastBlockHeight = t.RecentBlockResults[len(t.RecentBlockResults)-1].Height
	}
	// Number of new blocks not already stored
	numNew := util.Min(int(t.RecentBlocks[len(t.RecentBlocks)-1].Block.Header.Height-lastBlockHeight)-1, 4)
	if numNew <= len(t.RecentBlockResults) {
		t.RecentBlockResults = t.RecentBlockResults[numNew:]
	} else {
		t.RecentBlockResults = []coretypes.ResultBlockResults{}
	}
	for numNew > 0 {
		nextHeight := t.RecentBlocks[len(t.RecentBlocks)-1].Block.Header.Height - int64(numNew)
		result, err := c.BlockResults(&nextHeight)
		if !util.ErrPrint(err) {
			t.RecentBlockResults = append(t.RecentBlockResults, *result)
		}
		numNew--
	}
}
