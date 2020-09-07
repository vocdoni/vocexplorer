package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/tendermint/tendermint/p2p"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	types "gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
	"google.golang.org/protobuf/proto"
)

// VochainStats is the type used by the public stats api
type VochainStats struct {
	BlockHeight       int64               `json:"block_height"`
	EntityCount       int64               `json:"entity_count"`
	EnvelopeCount     int64               `json:"envelope_count"`
	ProcessCount      int64               `json:"process_count"`
	TransactionHeight int64               `json:"transaction_height"`
	ValidatorCount    int64               `json:"validator_count"`
	BlockTime         *[5]int32           `json:"block_time"`
	BlockTimeStamp    int32               `json:"block_time_stamp"`
	ChainID           string              `json:"chain_id"`
	GatewayAPIList    []string            `json:"gateway_api_list"`
	GatewayHealth     int32               `json:"gateway_health"`
	GenesisTimeStamp  time.Time           `json:"genesis_time_stamp"`
	Height            int64               `json:"height"`
	NodeInfo          p2p.DefaultNodeInfo `json:"node_info"`
	SyncInfo          coretypes.SyncInfo  `json:"sync_info"`
}

//PingServer pings the web server
func PingServer() bool {
	c := &http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := c.Get("/ping")
	if err != nil || resp == nil {
		return false
	}
	return true
}

// For requests where we don't want to ReadAll the response body
func requestBody(url string) (io.ReadCloser, bool) {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get(url)
	if err != nil {
		log.Error(err)
		return resp.Body, false
	}
	if resp.StatusCode != http.StatusOK {
		return resp.Body, false
	}
	return resp.Body, true
}

func request(url string) ([]byte, bool) {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get(url)
	if err != nil {
		log.Error(err)
		return []byte{}, false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []byte{}, false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
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
		if err != nil {
			log.Error(err)
		}
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
		if err != nil {
			log.Error(err)
		}
	}
	return heightMap.GetHeights(), true
}

//GetProcessEnvelopeCount returns the height of envelopes belonging to given process stored by the database
func GetProcessEnvelopeCount(process string) (int64, bool) {
	return getHeight("/api/envprocheight/?process=" + process)
}

//GetProcessEnvelopeCountMap returns the entire map of process envelope heights
func GetProcessEnvelopeCountMap() (map[string]int64, bool) {
	return getHeightMap("/api/heightmap/?key=" + config.ProcessEnvelopeCountMapKey)
}

//GetEnvelopeCount returns the latest envelope height stored by the database
func GetEnvelopeCount() (int64, bool) {
	return getHeight("/api/height/?key=" + config.LatestEnvelopeCountKey)
}

//GetProcessCount returns the latest process height stored by the database
func GetProcessCount() (int64, bool) {
	return getHeight("/api/height/?key=" + config.LatestProcessCountKey)
}

//GetEntityCount returns the latest envelope height stored by the database
func GetEntityCount() (int64, bool) {
	return getHeight("/api/height/?key=" + config.LatestEntityCountKey)
}

//GetEntityProcessCount returns the number of processes belonging to a
func GetEntityProcessCount(entity string) (int64, bool) {
	return getHeight("/api/entityprocheight/?entity=" + entity)
}

//GetEntityProcessCountMap returns the entire map of entity process heights
func GetEntityProcessCountMap() (map[string]int64, bool) {
	return getHeightMap("/api/heightmap/?key=" + config.EntityProcessCountMapKey)
}

//GetBlockHeight returns the latest block height stored by the database
func GetBlockHeight() (int64, bool) {
	return getHeight("/api/height/?key=" + config.LatestBlockHeightKey)
}

//GetTxHeight returns the latest tx height stored by the database
func GetTxHeight() (int64, bool) {
	return getHeight("/api/height/?key=" + config.LatestTxHeightKey)
}

//GetValidatorBlockHeight returns the height of blocks belonging to given validator stored by the database
func GetValidatorBlockHeight(proposer string) (int64, bool) {
	return getHeight("/api/numblocksvalidator/?proposer=" + proposer)
}

//GetValidatorCount returns the latest validator count stored by the database
func GetValidatorCount() (int64, bool) {
	return getHeight("/api/height/?key=" + config.LatestValidatorCountKey)
}

//GetBlockList returns a list of blocks from the database
func GetBlockList(i int) ([config.ListSize]*types.StoreBlock, bool) {
	body, ok := request("/api/listblocks/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]*types.StoreBlock{}, false
	}
	var rawBlockList types.ItemList
	err := proto.Unmarshal(body, &rawBlockList)
	if err != nil {
		log.Error(err)
	}
	var blockList [config.ListSize]*types.StoreBlock
	for i, rawBlock := range rawBlockList.GetItems() {
		if len(rawBlock) > 0 {
			var block types.StoreBlock
			err = proto.Unmarshal(rawBlock, &block)
			blockList[i] = &block
			if err != nil {
				log.Error(err)
			}
		}
	}
	return blockList, true
}

//GetBlockListByValidator returns a list of blocks with given proposer from the database
func GetBlockListByValidator(i int, proposer []byte) ([config.ListSize]*types.StoreBlock, bool) {
	body, ok := request("/api/listblocksvalidator/?from=" + util.IntToString(i) + "&proposer=" + util.HexToString(proposer))
	if !ok {
		return [config.ListSize]*types.StoreBlock{}, false
	}
	var rawBlockList types.ItemList
	err := proto.Unmarshal(body, &rawBlockList)
	if err != nil {
		log.Error(err)
	}
	var blockList [config.ListSize]*types.StoreBlock
	for i, rawBlock := range rawBlockList.GetItems() {
		if len(rawBlock) > 0 {
			var block types.StoreBlock
			err = proto.Unmarshal(rawBlock, &block)
			blockList[i] = &block
			if err != nil {
				log.Error(err)
			}
		}
	}
	return blockList, true
}

//GetValidatorList returns a list of validators from the database
func GetValidatorList(i int) ([config.ListSize]*types.Validator, bool) {
	body, ok := request("/api/listvalidators/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]*types.Validator{}, false
	}
	var rawValidatorList types.ItemList
	err := proto.Unmarshal(body, &rawValidatorList)
	if err != nil {
		log.Error(err)
	}
	var validatorList [config.ListSize]*types.Validator
	for i, rawVal := range rawValidatorList.GetItems() {
		if len(rawVal) > 0 {
			var validator types.Validator
			err = proto.Unmarshal(rawVal, &validator)
			validatorList[i] = &validator
			if err != nil {
				log.Error(err)
			}
		}
	}
	return validatorList, true
}

//GetStoreBlock returns a single block from the database
func GetStoreBlock(i int64) (*types.StoreBlock, bool) {
	body, ok := request("/api/block/?height=" + util.IntToString(i))
	if !ok {
		return &types.StoreBlock{}, false
	}
	var block types.StoreBlock
	err := proto.Unmarshal(body, &block)
	if err != nil {
		log.Error(err)
	}
	return &block, true
}

//GetTxList returns a list of transactions from the database
func GetTxList(from int) ([config.ListSize]*types.SendTx, bool) {
	body, ok := request("/api/listtxs/?from=" + util.IntToString(from))
	if !ok {
		return [config.ListSize]*types.SendTx{}, false
	}
	var rawTxList types.ItemList
	err := proto.Unmarshal(body, &rawTxList)
	if err != nil {
		log.Error(err)
	}
	var txList [config.ListSize]*types.SendTx
	for i, rawTx := range rawTxList.GetItems() {
		if len(rawTx) > 0 {
			var tx types.SendTx
			err = proto.Unmarshal(rawTx, &tx)
			if err != nil {
				log.Error(err)
			}
			txList[i] = &tx
		}
	}
	return txList, true
}

//GetTx returns a transaction from the database
func GetTx(height int64) (*types.SendTx, bool) {
	body, ok := request("/api/tx/?id=" + util.IntToString(height))
	if !ok {
		return &types.SendTx{}, false
	}
	var tx types.SendTx
	if len(body) > 0 {
		err := proto.Unmarshal(body, &tx)
		if err != nil {
			log.Error(err)
		}
	}
	return &tx, true
}

//GetTxHeightFromHash finds the height corresponding to a given tx hash
func GetTxHeightFromHash(hash string) (int64, bool) {
	body, ok := request("/api/txhash/?hash=" + hash)
	if !ok {
		return 0, false
	}
	var height types.Height
	err := proto.Unmarshal(body, &height)
	if err != nil {
		log.Error(err)
	}
	return height.GetHeight(), true
}

//GetValidator returns a single validator from the database
func GetValidator(address string) (*types.Validator, bool) {
	body, ok := request("/api/validator/?id=" + address)
	if !ok {
		return &types.Validator{}, false
	}
	var validator types.Validator
	err := proto.Unmarshal(body, &validator)
	if err != nil {
		log.Error(err)
	}
	return &validator, true
}

//GetEnvelopeList returns a list of envelopes from the database
func GetEnvelopeList(i int) ([config.ListSize]*types.Envelope, bool) {
	body, ok := request("/api/listenvelopes/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]*types.Envelope{}, false
	}
	var rawEnvList types.ItemList
	err := proto.Unmarshal(body, &rawEnvList)
	if err != nil {
		log.Error(err)
	}
	var envList [config.ListSize]*types.Envelope
	for i, rawEnvelope := range rawEnvList.GetItems() {
		if len(rawEnvelope) > 0 {
			envelope := new(types.Envelope)
			err = proto.Unmarshal(rawEnvelope, envelope)
			envList[i] = envelope
			if err != nil {
				log.Error(err)
			}
		}
	}
	return envList, true
}

//GetEnvelope gets a single envelope by global height
func GetEnvelope(height int64) (*types.Envelope, bool) {
	body, ok := request("/api/envelope/?height=" + util.IntToString(height))
	if !ok {
		return &types.Envelope{}, false
	}
	envelope := new(types.Envelope)
	err := proto.Unmarshal(body, envelope)
	if err != nil {
		log.Error(err)
	}
	return envelope, true
}

//GetEnvelopeHeightFromNullifier finds the height corresponding to a given envelope nullifier
func GetEnvelopeHeightFromNullifier(hash string) (int64, bool) {
	body, ok := request("/api/envelopenullifier/?nullifier=" + hash)
	if !ok {
		return 0, false
	}
	var height types.Height
	err := proto.Unmarshal(body, &height)
	if err != nil {
		log.Error(err)
	}
	return height.GetHeight(), true
}

//GetEnvelopeListByProcess returns a list of envelopes by process
func GetEnvelopeListByProcess(i int, process string) ([config.ListSize]*types.Envelope, bool) {
	body, ok := request("/api/listenvelopesprocess/?from=" + util.IntToString(i) + "&process=" + process)
	if !ok {
		return [config.ListSize]*types.Envelope{}, false
	}
	var rawEnvList types.ItemList
	err := proto.Unmarshal(body, &rawEnvList)
	if err != nil {
		log.Error(err)
	}
	var envList [config.ListSize]*types.Envelope
	for i, rawEnvelope := range rawEnvList.GetItems() {
		if len(rawEnvelope) > 0 {
			envelope := new(types.Envelope)
			err = proto.Unmarshal(rawEnvelope, envelope)
			envList[i] = envelope
			if err != nil {
				log.Error(err)
			}
		}
	}
	return envList, true
}

//GetEntityList returns a list of entities from the database
func GetEntityList(i int) ([config.ListSize]string, bool) {
	body, ok := request("/api/listentities/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]string{}, false
	}
	var rawEntityList types.ItemList
	err := proto.Unmarshal(body, &rawEntityList)
	if err != nil {
		log.Error(err)
	}
	var entityList [config.ListSize]string
	for i, rawEntity := range rawEntityList.GetItems() {
		if len(rawEntity) > 0 {
			entity := strings.ToLower(util.HexToString(rawEntity))
			entityList[i] = entity
			if err != nil {
				log.Error(err)
			}
		}
	}
	return entityList, true
}

//GetProcessList returns a list of entities from the database
func GetProcessList(i int) ([config.ListSize]string, bool) {
	body, ok := request("/api/listprocesses/?from=" + util.IntToString(i))
	if !ok {
		return [config.ListSize]string{}, false
	}
	var rawProcessList types.ItemList
	err := proto.Unmarshal(body, &rawProcessList)
	if err != nil {
		log.Error(err)
	}
	var processList [config.ListSize]string
	for i, rawProcess := range rawProcessList.GetItems() {
		if len(rawProcess) > 0 {
			process := strings.ToLower(util.HexToString(rawProcess))
			processList[i] = process
			if err != nil {
				log.Error(err)
			}
		}
	}
	return processList, true
}

//GetProcessListByEntity returns a list of processes by entity
func GetProcessListByEntity(i int, entity string) ([config.ListSize]string, bool) {
	body, ok := request("/api/listprocessesbyentity/?from=" + util.IntToString(i) + "&entity=" + entity)
	if !ok {
		return [config.ListSize]string{}, false
	}
	var rawProcessList types.ItemList
	err := proto.Unmarshal(body, &rawProcessList)
	if err != nil {
		log.Error(err)
	}
	var envList [config.ListSize]string
	for i, rawProcess := range rawProcessList.GetItems() {
		if len(rawProcess) > 0 {
			envelope := strings.ToLower(util.HexToString(rawProcess))
			envList[i] = envelope
		}
	}
	return envList, true
}

//GetStats gets the latest blockchain statistics
func GetStats() (*VochainStats, bool) {
	body, ok := requestBody("/api/stats")
	defer body.Close()
	if !ok {
		return &VochainStats{}, false
	}
	stats := new(VochainStats)
	err := json.NewDecoder(body).Decode(&stats)
	if err != nil {
		log.Error(err)
		return stats, false
	}
	return stats, true
}