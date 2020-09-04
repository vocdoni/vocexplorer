package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
	"google.golang.org/protobuf/proto"
)

//Ping pings the web server
func Ping() bool {
	c := &http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := c.Get("/ping")
	if err != nil || resp == nil {
		return false
	}
	return true
}

func request(url string) ([]byte, bool) {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get(url)
	if util.ErrPrint(err) {
		return []byte{}, false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) || resp.StatusCode != http.StatusOK {
		fmt.Println(string(body))
		return []byte{}, false
	}
	return body, true
}

func getHeight(url string) (int64, bool) {
	body, ok := request(url)
	if !ok {
		return 0, false
	}
	var height types.Height
	if len(body) > 0 {
		err := proto.Unmarshal(body, &height)
		util.ErrPrint(err)
	}
	return height.GetHeight(), true
}

func getHeightMap(url string) (map[string]int64, bool) {
	body, ok := request(url)
	if !ok {
		return map[string]int64{}, false
	}
	var heightMap types.HeightMap
	if len(body) > 0 {
		err := proto.Unmarshal(body, &heightMap)
		util.ErrPrint(err)
	}
	return heightMap.GetHeights(), true
}

//GetProcessEnvelopeCount returns the height of envelopes belonging to given process stored by the database
func GetProcessEnvelopeCount(process string) (int64, bool) {
	return getHeight("/db/envprocheight/?process=" + process)
}

//GetProcessEnvelopeCountMap returns the entire map of process envelope heights
func GetProcessEnvelopeCountMap() (map[string]int64, bool) {
	return getHeightMap("/db/heightmap/?key=" + config.ProcessEnvelopeCountMapKey)
}

//GetEnvelopeCount returns the latest envelope height stored by the database
func GetEnvelopeCount() (int64, bool) {
	return getHeight("/db/height/?key=" + config.LatestEnvelopeCountKey)
}

//GetProcessCount returns the latest process height stored by the database
func GetProcessCount() (int64, bool) {
	return getHeight("/db/height/?key=" + config.LatestProcessCountKey)
}

//GetEntityCount returns the latest envelope height stored by the database
func GetEntityCount() (int64, bool) {
	return getHeight("/db/height/?key=" + config.LatestEntityCountKey)
}

//GetEntityProcessCount returns the number of processes belonging to a
func GetEntityProcessCount(entity string) (int64, bool) {
	return getHeight("/db/entityprocheight/?entity=" + entity)
}

//GetEntityProcessCountMap returns the entire map of entity process heights
func GetEntityProcessCountMap() (map[string]int64, bool) {
	return getHeightMap("/db/heightmap/?key=" + config.EntityProcessCountMapKey)
}

//GetBlockHeight returns the latest block height stored by the database
func GetBlockHeight() (int64, bool) {
	return getHeight("/db/height/?key=" + config.LatestBlockHeightKey)
}

//GetTxHeight returns the latest tx height stored by the database
func GetTxHeight() (int64, bool) {
	return getHeight("db/height/?key=" + config.LatestTxHeightKey)
}

//GetValidatorBlockHeight returns the height of blocks belonging to given validator stored by the database
func GetValidatorBlockHeight(proposer string) (int64, bool) {
	return getHeight("/db/numblocksvalidator/?proposer=" + proposer)
}

//GetValidatorCount returns the latest validator count stored by the database
func GetValidatorCount() (int64, bool) {
	return getHeight("/db/height/?key=" + config.LatestValidatorCountKey)
}

//GetBlockList returns a list of blocks from the database
func GetBlockList(i int) ([config.ListSize]*types.StoreBlock, bool) {
	body, ok := request("/db/listblocks/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]*types.StoreBlock{}, false
	}
	var rawBlockList types.ItemList
	err := proto.Unmarshal(body, &rawBlockList)
	util.ErrPrint(err)
	var blockList [config.ListSize]*types.StoreBlock
	for i, rawBlock := range rawBlockList.GetItems() {
		if len(rawBlock) > 0 {
			var block types.StoreBlock
			err = proto.Unmarshal(rawBlock, &block)
			blockList[i] = &block
			util.ErrPrint(err)
		}
	}
	return blockList, true
}

//GetBlockListByValidator returns a list of blocks with given proposer from the database
func GetBlockListByValidator(i int, proposer []byte) ([config.ListSize]*types.StoreBlock, bool) {
	body, ok := request("/db/listblocksvalidator/?from=" + util.IntToString(i) + "&proposer=" + util.HexToString(proposer))
	if !ok {
		return [config.ListSize]*types.StoreBlock{}, false
	}
	var rawBlockList types.ItemList
	err := proto.Unmarshal(body, &rawBlockList)
	util.ErrPrint(err)
	var blockList [config.ListSize]*types.StoreBlock
	for i, rawBlock := range rawBlockList.GetItems() {
		if len(rawBlock) > 0 {
			var block types.StoreBlock
			err = proto.Unmarshal(rawBlock, &block)
			blockList[i] = &block
			util.ErrPrint(err)
		}
	}
	return blockList, true
}

//GetValidatorList returns a list of validators from the database
func GetValidatorList(i int) ([config.ListSize]*types.Validator, bool) {
	body, ok := request("/db/listvalidators/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]*types.Validator{}, false
	}
	var rawValidatorList types.ItemList
	err := proto.Unmarshal(body, &rawValidatorList)
	util.ErrPrint(err)
	var validatorList [config.ListSize]*types.Validator
	for i, rawVal := range rawValidatorList.GetItems() {
		if len(rawVal) > 0 {
			var validator types.Validator
			err = proto.Unmarshal(rawVal, &validator)
			validatorList[i] = &validator
			util.ErrPrint(err)
		}
	}
	return validatorList, true
}

//GetBlock returns a single block from the database
func GetBlock(i int64) (*types.StoreBlock, bool) {
	body, ok := request("/db/block/?height=" + util.IntToString(i))
	if !ok {
		return &types.StoreBlock{}, false
	}
	var block types.StoreBlock
	err := proto.Unmarshal(body, &block)
	util.ErrPrint(err)
	return &block, true
}

//GetTxList returns a list of transactions from the database
func GetTxList(from int) ([config.ListSize]*types.SendTx, bool) {
	body, ok := request("/db/listtxs/?from=" + util.IntToString(from))
	if !ok {
		return [config.ListSize]*types.SendTx{}, false
	}
	var rawTxList types.ItemList
	err := proto.Unmarshal(body, &rawTxList)
	util.ErrPrint(err)
	var txList [config.ListSize]*types.SendTx
	for i, rawTx := range rawTxList.GetItems() {
		if len(rawTx) > 0 {
			var tx types.SendTx
			err = proto.Unmarshal(rawTx, &tx)
			util.ErrPrint(err)
			txList[i] = &tx
		}
	}
	return txList, true
}

//GetTx returns a transaction from the database
func GetTx(height int64) (*types.SendTx, bool) {
	body, ok := request("/db/tx/?id=" + util.IntToString(height))
	if !ok {
		return &types.SendTx{}, false
	}
	var tx types.SendTx
	if len(body) > 0 {
		err := proto.Unmarshal(body, &tx)
		util.ErrPrint(err)
	}
	return &tx, true
}

//GetTxHeightFromHash finds the height corresponding to a given tx hash
func GetTxHeightFromHash(hash string) (int64, bool) {
	body, ok := request("/db/txhash/?hash=" + hash)
	if !ok {
		return 0, false
	}
	var height types.Height
	err := proto.Unmarshal(body, &height)
	util.ErrPrint(err)
	return height.GetHeight(), true
}

//GetValidator returns a single validator from the database
func GetValidator(address string) (*types.Validator, bool) {
	body, ok := request("/db/validator/?id=" + address)
	if !ok {
		return &types.Validator{}, false
	}
	var validator types.Validator
	err := proto.Unmarshal(body, &validator)
	util.ErrPrint(err)
	return &validator, true
}

//GetEnvelopeList returns a list of envelopes from the database
func GetEnvelopeList(i int) ([config.ListSize]*types.Envelope, bool) {
	body, ok := request("/db/listenvelopes/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]*types.Envelope{}, false
	}
	var rawEnvList types.ItemList
	err := proto.Unmarshal(body, &rawEnvList)
	util.ErrPrint(err)
	var envList [config.ListSize]*types.Envelope
	for i, rawEnvelope := range rawEnvList.GetItems() {
		if len(rawEnvelope) > 0 {
			envelope := new(types.Envelope)
			err = proto.Unmarshal(rawEnvelope, envelope)
			envList[i] = envelope
			util.ErrPrint(err)
		}
	}
	return envList, true
}

//GetEnvelope gets a single envelope by global height
func GetEnvelope(height int64) (*types.Envelope, bool) {
	body, ok := request("/db/envelope/?height=" + util.IntToString(height))
	if !ok {
		return &types.Envelope{}, false
	}
	envelope := new(types.Envelope)
	err := proto.Unmarshal(body, envelope)
	util.ErrPrint(err)
	return envelope, true
}

//GetEnvelopeHeightFromNullifier finds the height corresponding to a given envelope nullifier
func GetEnvelopeHeightFromNullifier(hash string) (int64, bool) {
	body, ok := request("/db/envelopenullifier/?nullifier=" + hash)
	if !ok {
		return 0, false
	}
	var height types.Height
	err := proto.Unmarshal(body, &height)
	util.ErrPrint(err)
	return height.GetHeight(), true
}

//GetEnvelopeListByProcess returns a list of envelopes by process
func GetEnvelopeListByProcess(i int, process string) ([config.ListSize]*types.Envelope, bool) {
	body, ok := request("/db/listenvelopesprocess/?from=" + util.IntToString(i) + "&process=" + process)
	if !ok {
		return [config.ListSize]*types.Envelope{}, false
	}
	var rawEnvList types.ItemList
	err := proto.Unmarshal(body, &rawEnvList)
	util.ErrPrint(err)
	var envList [config.ListSize]*types.Envelope
	for i, rawEnvelope := range rawEnvList.GetItems() {
		if len(rawEnvelope) > 0 {
			envelope := new(types.Envelope)
			err = proto.Unmarshal(rawEnvelope, envelope)
			envList[i] = envelope
			util.ErrPrint(err)
		}
	}
	return envList, true
}

//GetEntityList returns a list of entities from the database
func GetEntityList(i int) ([config.ListSize]string, bool) {
	body, ok := request("/db/listentities/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]string{}, false
	}
	var rawEntityList types.ItemList
	err := proto.Unmarshal(body, &rawEntityList)
	util.ErrPrint(err)
	var entityList [config.ListSize]string
	for i, rawEntity := range rawEntityList.GetItems() {
		if len(rawEntity) > 0 {
			entity := strings.ToLower(util.HexToString(rawEntity))
			entityList[i] = entity
			util.ErrPrint(err)
		}
	}
	return entityList, true
}

//GetProcessList returns a list of entities from the database
func GetProcessList(i int) ([config.ListSize]string, bool) {
	body, ok := request("/db/listprocesses/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]string{}, false
	}
	var rawProcessList types.ItemList
	err := proto.Unmarshal(body, &rawProcessList)
	util.ErrPrint(err)
	var processList [config.ListSize]string
	for i, rawProcess := range rawProcessList.GetItems() {
		if len(rawProcess) > 0 {
			process := strings.ToLower(util.HexToString(rawProcess))
			processList[i] = process
			util.ErrPrint(err)
		}
	}
	return processList, true
}

//GetProcessListByEntity returns a list of processes by entity
func GetProcessListByEntity(i int, entity string) ([config.ListSize]string, bool) {
	body, ok := request("/db/listprocessesbyentity/?from=" + util.IntToString(i) + "&entity=" + entity)
	if !ok {
		return [config.ListSize]string{}, false
	}
	var rawProcessList types.ItemList
	err := proto.Unmarshal(body, &rawProcessList)
	util.ErrPrint(err)
	var envList [config.ListSize]string
	for i, rawProcess := range rawProcessList.GetItems() {
		if len(rawProcess) > 0 {
			envelope := strings.ToLower(util.HexToString(rawProcess))
			envList[i] = envelope
		}
	}
	return envList, true
}
