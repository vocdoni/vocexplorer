package dbapi

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

func getHeight(url string) int64 {
	body, ok := request(url)
	if !ok {
		return 0
	}
	var height types.Height
	if len(body) > 0 {
		err := proto.Unmarshal(body, &height)
		util.ErrPrint(err)
	}
	return height.GetHeight()
}

func getHeightMap(url string) map[string]int64 {
	body, ok := request(url)
	if !ok {
		return map[string]int64{}
	}
	var heightMap types.HeightMap
	if len(body) > 0 {
		err := proto.Unmarshal(body, &heightMap)
		util.ErrPrint(err)
	}
	return heightMap.GetHeights()
}

//GetProcessEnvelopeHeight returns the height of envelopes belonging to given process stored by the database
func GetProcessEnvelopeHeight(process string) int64 {
	return getHeight("/db/envprocheight/?process=" + process)
}

//GetProcessEnvelopeHeightMap returns the entire map of process envelope heights
func GetProcessEnvelopeHeightMap() map[string]int64 {
	return getHeightMap("/db/heightmap/?key=" + config.ProcessEnvelopeHeightMapKey)
}

//GetEnvelopeHeight returns the latest envelope height stored by the database
func GetEnvelopeHeight() int64 {
	return getHeight("/db/height/?key=" + config.LatestEnvelopeHeightKey)
}

//GetProcessHeight returns the latest process height stored by the database
func GetProcessHeight() int64 {
	return getHeight("/db/height/?key=" + config.LatestProcessHeight)
}

//GetEntityHeight returns the latest envelope height stored by the database
func GetEntityHeight() int64 {
	return getHeight("/db/height/?key=" + config.LatestEntityHeight)
}

//GetEntityProcessHeight returns the number of processes belonging to a
func GetEntityProcessHeight(entity string) int64 {
	return getHeight("/db/entityprocheight/?entity=" + entity)
}

//GetEntityProcessHeightMap returns the entire map of entity process heights
func GetEntityProcessHeightMap() map[string]int64 {
	return getHeightMap("/db/heightmap/?key=" + config.EntityProcessHeightMapKey)
}

//GetBlockHeight returns the latest block height stored by the database
func GetBlockHeight() int64 {
	return getHeight("/db/height/?key=" + config.LatestBlockHeightKey)
}

//GetTxHeight returns the latest tx height stored by the database
func GetTxHeight() int64 {
	return getHeight("db/height/?key=" + config.LatestTxHeightKey)
}

//GetValidatorBlockHeight returns the height of blocks belonging to given validator stored by the database
func GetValidatorBlockHeight(proposer string) int64 {
	return getHeight("/db/numblocksvalidator/?proposer=" + proposer)
}

//GetBlockList returns a list of blocks from the database
func GetBlockList(i int) [config.ListSize]*types.StoreBlock {
	body, ok := request("/db/listblocks/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]*types.StoreBlock{}
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
	return blockList
}

//GetBlockListByValidator returns a list of blocks with given proposer from the database
func GetBlockListByValidator(i int, proposer []byte) [config.ListSize]*types.StoreBlock {
	if i < config.ListSize {
		i = config.ListSize
	}
	body, ok := request("/db/listblocksvalidator/?from=" + util.IntToString(i) + "&proposer=" + util.HexToString(proposer))
	if !ok {
		return [config.ListSize]*types.StoreBlock{}
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
	return blockList
}

//GetBlock returns a single block from the database
func GetBlock(i int64) *types.StoreBlock {
	body, ok := request("/db/block/?height=" + util.IntToString(i))
	if !ok {
		return &types.StoreBlock{}
	}
	var block types.StoreBlock
	err := proto.Unmarshal(body, &block)
	util.ErrPrint(err)
	return &block
}

//GetTxList returns a list of transactions from the database
func GetTxList(from int) [config.ListSize]*types.SendTx {
	body, ok := request("/db/listtxs/?from=" + util.IntToString(from))
	if !ok {
		return [config.ListSize]*types.SendTx{}
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
	return txList
}

//GetTx returns a transaction from the database
func GetTx(height int64) *types.SendTx {
	body, ok := request("/db/tx/?id=" + util.IntToString(height))
	if !ok {
		return &types.SendTx{}
	}
	var tx types.SendTx
	if len(body) > 0 {
		err := proto.Unmarshal(body, &tx)
		util.ErrPrint(err)
	}
	return &tx
}

//GetValidator returns a single validator from the database
func GetValidator(address string) *types.Validator {
	body, ok := request("/db/validator/?id=" + address)
	if !ok {
		return &types.Validator{}
	}
	var validator types.Validator
	err := proto.Unmarshal(body, &validator)
	util.ErrPrint(err)
	return &validator
}

//GetEnvelopeList returns a list of envelopes from the database
func GetEnvelopeList(i int) [config.ListSize]*types.Envelope {
	body, ok := request("/db/listenvelopes/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]*types.Envelope{}
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
	return envList
}

//GetEnvelope gets a single envelope by global height
func GetEnvelope(height int64) *types.Envelope {
	body, ok := request("/db/envelope/?height=" + util.IntToString(height))
	if !ok {
		return &types.Envelope{}
	}
	envelope := new(types.Envelope)
	err := proto.Unmarshal(body, envelope)
	util.ErrPrint(err)
	return envelope
}

//GetEnvelopeListByProcess returns a list of envelopes by process
func GetEnvelopeListByProcess(i int, process string) [config.ListSize]*types.Envelope {
	body, ok := request("/db/listenvelopesprocess/?from=" + util.IntToString(i) + "&process=" + process)
	if !ok {
		return [config.ListSize]*types.Envelope{}
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
	return envList
}

//GetEntityList returns a list of entities from the database
func GetEntityList(i int) [config.ListSize]string {
	body, ok := request("/db/listentities/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]string{}
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
	return entityList
}

//GetProcessList returns a list of entities from the database
func GetProcessList(i int) [config.ListSize]string {
	body, ok := request("/db/listprocesses/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]string{}
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
	return processList
}

//GetProcessListByEntity returns a list of processes by entity
func GetProcessListByEntity(i int, entity string) [config.ListSize]string {
	body, ok := request("/db/listprocessesbyentity/?from=" + util.IntToString(i) + "&entity=" + entity)
	if !ok {
		return [config.ListSize]string{}
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
	return envList
}
