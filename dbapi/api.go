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
)

//GetBlockList returns a list of blocks from the database
func GetBlockList(i int) [config.ListSize]types.StoreBlock {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/db/listblocks/?from=" + util.IntToString(i))
	if util.ErrPrint(err) {
		return [config.ListSize]types.StoreBlock{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return [config.ListSize]types.StoreBlock{}
	}
	var blockList [config.ListSize]types.StoreBlock
	err = json.Unmarshal(body, &blockList)
	util.ErrPrint(err)
	return blockList
}

// //GetBlockHash returns the hash of the block with the given height
// func GetBlockHash(i int) string {
// 	// resp, err := http.Get("/db/hash/?key=" + config.BlockHeightPrefix + util.IntToString(i))
// 	c := &http.Client{
// 		Timeout: 10 * time.Second,
// 	}
// 	resp, err := c.Get("/db/hash/?key=0220")
// 	if util.ErrPrint(err) {
// 		return ""
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != 200 {
// 		log.Errorf("Request not valid")
// 		return ""
// 	}
// 	hash, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
// 	if util.ErrPrint(err) {
// 		return ""
// 	}
// 	fmt.Println("Got hash")

// 	return string(hash)
// }

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
	var height int64
	err = json.Unmarshal(body, &height)
	util.ErrPrint(err)
	return height
}

//GetTxList returns a list of transactions from the database
func GetTxList(from int) [config.ListSize]types.SendTx {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/db/listtxs/?from=" + util.IntToString(from))
	if util.ErrPrint(err) {
		return [config.ListSize]types.SendTx{}
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Errorf("Request not valid")
		return [config.ListSize]types.SendTx{}
	}
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return [config.ListSize]types.SendTx{}
	}
	var txList [config.ListSize]types.SendTx
	err = json.Unmarshal(body, &txList)
	util.ErrPrint(err)
	return txList
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
	var height int64
	err = json.Unmarshal(body, &height)
	util.ErrPrint(err)
	return height
}
