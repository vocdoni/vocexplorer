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
	ResultStatus     *coretypes.ResultStatus
	Genesis          *types.GenesisDoc
	BlockList        [config.SearchPageSmall]coretypes.ResultBlock
	BlockListResults [config.SearchPageSmall]coretypes.ResultBlockResults
	TxCount          int
	TxList           []*coretypes.ResultTx
	ChainID          int
	AppVersion       int
	MaxBlockSize     int
	NumValidators    int
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
func UpdateTendermintInfo(c *http.HTTP, t *TendermintInfo, i int) {
	t.GetHealth(c)
	t.GetAllTxs(c)
	t.GetGenesis(c)
	UpdateBlockList(c, t, i)
}

//UpdateBlockList updates the latest blocks starting at index i
func UpdateBlockList(c *http.HTTP, t *TendermintInfo, i int) {
	t.GetBlockList(c, i)
	// t.GetBlockListResults(c, i)
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

// GetBlockList keeps a list of current blocks
func (t *TendermintInfo) GetBlockList(c *http.HTTP, index int) {
	lastBlockHeight := 0
	lastBlock := t.BlockList[config.SearchPageSmall-1].Block
	if lastBlock != nil {
		lastBlockHeight = int(lastBlock.Header.Height)
	}
	fmt.Println("Lash block height: " + util.IntToString(lastBlockHeight))
	// Offset from last index to new one, so we can recycle fetched blocks
	offset := int(t.ResultStatus.SyncInfo.LatestBlockHeight) - 1 - index - lastBlockHeight
	if offset == 0 {
		return
	}
	fmt.Println("initial offset: " + util.IntToString(offset))
	if offset > 0 {
		if offset < config.SearchPageSmall {
			for i := 0; i < config.SearchPageSmall-offset; i++ {
				t.BlockList[i] = t.BlockList[i+offset]
			}
		} else {
			offset = config.SearchPageSmall
		}
		for offset > 0 {
			nextHeight := t.ResultStatus.SyncInfo.LatestBlockHeight - int64(index+offset)
			fmt.Println("next height" + util.IntToString(nextHeight))
			result, err := c.Block(&nextHeight)
			if !util.ErrPrint(err) {
				t.BlockList[config.SearchPageSmall-offset] = *result
			}
			offset--
		}
	} else if offset < 0 {
		if offset > 0-config.SearchPageSmall {
			offset = 0 - offset
			for i := 0; i < config.SearchPageSmall-offset; i++ {
				t.BlockList[i+offset] = t.BlockList[i]
			}
		} else {
			offset = config.SearchPageSmall
		}
		for offset > 0 {
			nextHeight := t.ResultStatus.SyncInfo.LatestBlockHeight - int64(index+config.SearchPageSmall-offset+1)
			fmt.Println("next height" + util.IntToString(nextHeight))
			result, err := c.Block(&nextHeight)
			if !util.ErrPrint(err) {
				t.BlockList[offset-1] = *result
			}
			offset--
		}
	}
}

// GetBlockListResults keeps a list of current blocks
func (t *TendermintInfo) GetBlockListResults(c *http.HTTP, index int) {
	lastBlockHeight := 0
	lastBlockHeight = int(t.BlockListResults[config.SearchPageSmall-1].Height)
	// Offset from last index to new one, so we can recycle fetched blocks
	offset := int(t.ResultStatus.SyncInfo.LatestBlockHeight) - 1 - index - lastBlockHeight
	if offset == 0 {
		return
	}
	if offset > 0 {
		if offset < config.SearchPageSmall {
			for i := 0; i < config.SearchPageSmall-offset; i++ {
				t.BlockListResults[i] = t.BlockListResults[i+offset]
			}
		} else {
			offset = config.SearchPageSmall
		}
		for offset > 0 {
			nextHeight := t.ResultStatus.SyncInfo.LatestBlockHeight - int64(index+offset-1)
			result, err := c.BlockResults(&nextHeight)
			if !util.ErrPrint(err) {
				t.BlockListResults[config.SearchPageSmall-offset] = *result
			}
			offset--
		}
	} else if offset < 0 {
		if offset > 0-config.SearchPageSmall {
			offset = 0 - offset
			for i := 0; i < config.SearchPageSmall-offset; i++ {
				t.BlockListResults[i+offset] = t.BlockListResults[i]
			}
		} else {
			offset = config.SearchPageSmall
		}
		for offset > 0 {
			nextHeight := t.ResultStatus.SyncInfo.LatestBlockHeight - int64(index+config.SearchPageSmall-offset)
			result, err := c.BlockResults(&nextHeight)
			if !util.ErrPrint(err) {
				t.BlockListResults[offset-1] = *result
			}
			offset--
		}
	}
}
