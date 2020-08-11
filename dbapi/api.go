package dbapi

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
	"google.golang.org/protobuf/proto"
)

//GetBlockList returns a list of blocks from the database
func GetBlockList(i int) [config.ListSize]*types.StoreBlock {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/db/listblocks/?from=" + util.IntToString(i))
	if util.ErrPrint(err) {
		return [config.ListSize]*types.StoreBlock{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return [config.ListSize]*types.StoreBlock{}
	}
	var rawBlockList [config.ListSize][]byte
	err = json.Unmarshal(body, &rawBlockList)
	util.ErrPrint(err)
	var blockList [config.ListSize]*types.StoreBlock
	for i, rawTx := range rawBlockList {
		err = proto.Unmarshal(rawTx, blockList[i])
		util.ErrPrint(err)
	}
	return blockList
}

//GetBlock returns a single block from the database
func GetBlock(i int64) *types.StoreBlock {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/db/block/?id=" + util.IntToString(i))
	if util.ErrPrint(err) {
		return &types.StoreBlock{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return &types.StoreBlock{}
	}
	var block *types.StoreBlock
	err = proto.Unmarshal(body, block)
	util.ErrPrint(err)
	return block
}

//GetBlockHeight returns the latest block height stored by the database
func GetBlockHeight() int64 {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("db/height/?key=" + config.LatestBlockHeightKey)
	if util.ErrPrint(err) {
		return 0
	}
	if resp.StatusCode != 200 {
		log.Errorf("Request not valid")
		return 0
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return 0
	}
	var height types.Height
	err = proto.Unmarshal(body, &height)
	util.ErrPrint(err)
	return height.GetHeight()
}

//GetTxList returns a list of transactions from the database
func GetTxList(from int) [config.ListSize]*types.SendTx {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/db/listtxs/?from=" + util.IntToString(from))
	if util.ErrPrint(err) {
		return [config.ListSize]*types.SendTx{}
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Errorf("Request not valid")
		return [config.ListSize]*types.SendTx{}
	}
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return [config.ListSize]*types.SendTx{}
	}
	var rawTxList [config.ListSize][]byte
	err = json.Unmarshal(body, &rawTxList)
	util.ErrPrint(err)
	var txList [config.ListSize]*types.SendTx
	for i, rawTx := range rawTxList {
		err = proto.Unmarshal(rawTx, txList[i])
		util.ErrPrint(err)
	}
	return txList
}

//GetTx returns a transaction from the database
func GetTx(height int64) *types.SendTx {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/db/tx/?id=" + util.IntToString(height))
	if util.ErrPrint(err) {
		return &types.SendTx{}
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Errorf("Request not valid")
		return &types.SendTx{}
	}
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return &types.SendTx{}
	}
	var tx *types.SendTx
	err = proto.Unmarshal(body, tx)
	util.ErrPrint(err)
	return tx
}

//GetTxHeight returns the latest tx height stored by the database
func GetTxHeight() int64 {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("db/height/?key=" + config.LatestTxHeightKey)
	if util.ErrPrint(err) {
		return 0
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Errorf("Request not valid")
		return 0
	}
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return 0
	}
	var height types.Height
	err = proto.Unmarshal(body, &height)
	util.ErrPrint(err)
	return height.GetHeight()
}
