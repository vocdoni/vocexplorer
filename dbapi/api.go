package dbapi

import (
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
	var rawBlockList types.ItemList
	err = proto.Unmarshal(body, &rawBlockList)
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
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/db/listblocksvalidator/?from=" + util.IntToString(i) + "&proposer=" + util.HexToString(proposer))
	if util.ErrPrint(err) {
		return [config.ListSize]*types.StoreBlock{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return [config.ListSize]*types.StoreBlock{}
	}
	var rawBlockList types.ItemList
	err = proto.Unmarshal(body, &rawBlockList)
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

//GetValidatorBlockHeight returns the height of blocks belonging to given validator stored by the database
func GetValidatorBlockHeight(proposer string) int64 {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/db/numblocksvalidator/?proposer=" + proposer)
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
	if len(body) > 0 {
		err = proto.Unmarshal(body, &height)
		util.ErrPrint(err)
	}
	return height.GetHeight()
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
	var block types.StoreBlock
	err = proto.Unmarshal(body, &block)
	util.ErrPrint(err)
	return &block
}

//GetBlockHeight returns the latest block height stored by the database
func GetBlockHeight() int64 {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/db/height/?key=" + config.LatestBlockHeightKey)
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
	if len(body) > 0 {
		err = proto.Unmarshal(body, &height)
		util.ErrPrint(err)
	}
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
	var rawTxList types.ItemList
	err = proto.Unmarshal(body, &rawTxList)
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
	var tx types.SendTx
	if len(body) > 0 {
		err = proto.Unmarshal(body, &tx)
		util.ErrPrint(err)
	}
	return &tx
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
	if len(body) > 0 {
		err = proto.Unmarshal(body, &height)
		util.ErrPrint(err)
	}
	return height.GetHeight()
}

//GetValidator returns a single validator from the database
func GetValidator(address string) *types.Validator {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/db/validator/?id=" + address)
	if util.ErrPrint(err) {
		return &types.Validator{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return &types.Validator{}
	}
	var validator types.Validator
	err = proto.Unmarshal(body, &validator)
	util.ErrPrint(err)
	return &validator
}

//GetEnvelopeList returns a list of envelopes from the database
func GetEnvelopeList(i int) [config.ListSize]*types.Envelope {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/db/listenvelopes/?from=" + util.IntToString(i))
	if util.ErrPrint(err) {
		return [config.ListSize]*types.Envelope{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return [config.ListSize]*types.Envelope{}
	}
	var rawEnvList types.ItemList
	err = proto.Unmarshal(body, &rawEnvList)
	util.ErrPrint(err)
	var envList [config.ListSize]*types.Envelope
	for i, rawEnvelope := range rawEnvList.GetItems() {
		if len(rawEnvelope) > 0 {
			err = proto.Unmarshal(rawEnvelope, envList[i])
			util.ErrPrint(err)
		}
	}
	return envList
}

//GetEnvelopeListByProcess returns a list of envelopes by process
func GetEnvelopeListByProcess(i int, process string) [config.ListSize]*types.Envelope {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/db/listenvelopesprocess/?from=" + util.IntToString(i) + "&process=" + process)
	if util.ErrPrint(err) {
		return [config.ListSize]*types.Envelope{}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if util.ErrPrint(err) {
		return [config.ListSize]*types.Envelope{}
	}
	var rawEnvList types.ItemList
	err = proto.Unmarshal(body, &rawEnvList)
	util.ErrPrint(err)
	var envList [config.ListSize]*types.Envelope
	for i, rawEnvelope := range rawEnvList.GetItems() {
		if len(rawEnvelope) > 0 {
			err = proto.Unmarshal(rawEnvelope, envList[i])
			util.ErrPrint(err)
		}
	}
	return envList
}

//GetProcessEnvelopeHeight returns the height of envelopes belonging to given process stored by the database
func GetProcessEnvelopeHeight(process string) int64 {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := c.Get("/db/envprocheight/?process=" + process)
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
	if len(body) > 0 {
		err = proto.Unmarshal(body, &height)
		util.ErrPrint(err)
	}
	return height.GetHeight()
}
