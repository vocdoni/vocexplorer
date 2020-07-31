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
	// var strings []string
	var blockList [config.ListSize]types.StoreBlock
	err = json.Unmarshal(body, &blockList)
	// err = json.Unmarshal(body, &strings)
	util.ErrPrint(err)

	// var cdc = amino.NewCodec()
	// cdc.RegisterConcrete(types.StoreBlock{}, "storeBlock", nil)
	// var blockList [config.ListSize]types.StoreBlock
	// for i, val := range strings {
	// 	err := cdc.UnmarshalBinaryLengthPrefixed([]byte(val), &blockList[i])
	// 	util.ErrPrint(err)
	// }
	return blockList
}
