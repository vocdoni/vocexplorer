package dbapi

import (
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"net/http"

	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//GetBlockList returns a list of blocks from the database
func GetBlockList(i int) [config.ListSize]coretypes.ResultBlock {
	resp, err := http.Get("/db/list/?prefix=" + config.BlockPrefix + "&from=" + util.IntToString(i))
	if util.ErrPrint(err) {
		return [config.ListSize]coretypes.ResultBlock{}
	}
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return [config.ListSize]coretypes.ResultBlock{}
	}
	decBuf := bytes.NewBuffer(body)
	blockList := [config.ListSize]coretypes.ResultBlock{}
	err = gob.NewDecoder(decBuf).Decode(&blockList)
	return blockList

}
