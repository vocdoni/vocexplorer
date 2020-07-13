package rpc

import (
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
	RecentBlocks       *[]coretypes.ResultBlock
	RecentBlockResults *[]coretypes.ResultBlockResults
	TxCount            int
	TxList             []*coretypes.ResultTx
	ChainID            int
	AppVersion         int
	MaxBlockSize       int
	NumValidators      int
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

//UpdateTendermintInfo updates the tendermint info
func UpdateTendermintInfo(c *http.HTTP, t *TendermintInfo) {
	t.GetHealth(c)
	t.GetAllTxs(c)
	t.GetGenesis(c)
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
	if t.RecentBlocks != nil {
		lastBlockHeight = (*t.RecentBlocks)[3].Block.Header.Height
	}
	numNew := t.ResultStatus.SyncInfo.LatestBlockHeight - lastBlockHeight
	*t.RecentBlocks = (*t.RecentBlocks)[numNew:3]
	for numNew < 0 {
		nextHeight := t.ResultStatus.SyncInfo.LatestBlockHeight - numNew
		result, err := c.Block(&nextHeight)
		if !util.ErrPrint(err) {
			*t.RecentBlocks = append(*t.RecentBlocks, *result)
		}
		numNew--
	}
}

// TODO complete this function so that we can get num txs per block
// // GetRecentBlockResults keeps a list of the four most recent blocks' results
// func (t *TendermintInfo) GetRecentBlockResults(c *http.HTTP) {
// 	var lastBlockHeight int64
// 	if t.RecentBlocks != nil {
// 		lastBlockHeight = (*t.RecentBlocks)[3].Block.Header.Height
// 	}
// 	numNew := t.ResultStatus.SyncInfo.LatestBlockHeight - lastBlockHeight
// 	*t.RecentBlocks = (*t.RecentBlocks)[numNew:4]
// 	for numNew < 0 {
// 		nextHeight := t.ResultStatus.SyncInfo.LatestBlockHeight - numNew
// 		result, err := c.Block(&nextHeight)
// 		if !util.ErrPrint(err) {
// 			*t.RecentBlocks = append(*t.RecentBlocks, *result)
// 		}
// 		numNew--
// 	}
// }
