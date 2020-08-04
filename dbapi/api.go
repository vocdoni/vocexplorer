package dbapi

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//GetBlockList returns a list of blocks from the database
func GetBlockList(i int) [config.ListSize]types.StoreBlock {
	resp, err := http.Get("/db/list/?prefix=" + config.BlockPrefix + "&from=" + util.IntToString(i))
	if util.ErrPrint(err) {
		return [config.ListSize]types.StoreBlock{}
	}
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return [config.ListSize]types.StoreBlock{}
	}
	var blockList [config.ListSize]types.StoreBlock
	err = json.Unmarshal(body, &blockList)
	util.ErrPrint(err)
	return blockList
}

//GetBlockHeight returns the latest block height stored by the database
func GetBlockHeight() int64 {
	resp, err := http.Get("db/height/?key=" + config.LatestBlockHeightKey)
	if util.ErrPrint(err) {
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
