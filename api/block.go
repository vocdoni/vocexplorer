package api

import (
	"encoding/json"

	types "gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//GetBlock fetches a single block from the vochain node
func GetBlock(height int64) (*Block, bool) {
	body, ok := requestBody("/api/block/?height=" + util.IntToString(height))
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return &Block{}, false
	}
	block := new(Block)
	err := json.NewDecoder(body).Decode(&block)
	if err != nil {
		logger.Error(err)
		return block, false
	}
	return block, true
}

//GetStoreBlock returns a single block from the database
func GetStoreBlock(i int64) (*types.StoreBlock, bool) {
	body, ok := requestBody("/api/blockheader/?height=" + util.IntToString(i))
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return &types.StoreBlock{}, false
	}
	var block types.StoreBlock
	err := json.NewDecoder(body).Decode(&block)
	if err != nil {
		logger.Error(err)
	}
	return &block, true
}

//GetBlockList returns a list of blocks from the database
func GetBlockList(i int) ([config.ListSize]*types.StoreBlock, bool) {
	body, ok := requestBody("/api/listblocks/?from=" + util.IntToString(i))
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return [config.ListSize]*types.StoreBlock{}, false
	}
	var rawBlockList types.ItemList
	err := json.NewDecoder(body).Decode(&rawBlockList)
	if err != nil {
		logger.Error(err)
	}
	var blockList [config.ListSize]*types.StoreBlock
	for i, rawBlock := range rawBlockList.Items {
		if len(rawBlock) > 0 {
			var block types.StoreBlock
			err := json.Unmarshal(rawBlock, &block)
			blockList[i] = &block
			if err != nil {
				logger.Error(err)
			}
		}
	}
	return blockList, true
}

//GetBlockSearch returns a list of blocks from the database according to the search term
func GetBlockSearch(term string) ([config.ListSize]*types.StoreBlock, bool) {
	itemList, ok := getItemList(&types.StoreBlock{}, "/api/blocksearch/?term="+term)
	if !ok {
		return [config.ListSize]*types.StoreBlock{}, false
	}
	list, ok := itemList.([config.ListSize]*types.StoreBlock)
	if !ok {
		return [config.ListSize]*types.StoreBlock{}, false
	}
	return list, true
}

//GetBlockListByValidator returns a list of blocks with given proposer from the database
func GetBlockListByValidator(i int, proposer []byte) ([config.ListSize]*types.StoreBlock, bool) {
	body, ok := requestBody("/api/listblocksvalidator/?from=" + util.IntToString(i) + "&proposer=" + util.HexToString(proposer))
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return [config.ListSize]*types.StoreBlock{}, false
	}
	var rawBlockList types.ItemList
	err := json.NewDecoder(body).Decode(&rawBlockList)
	if err != nil {
		logger.Error(err)
	}
	var blockList [config.ListSize]*types.StoreBlock
	for i, rawBlock := range rawBlockList.Items {
		if len(rawBlock) > 0 {
			var block types.StoreBlock
			err := json.Unmarshal(rawBlock, &block)
			blockList[i] = &block
			if err != nil {
				logger.Error(err)
			}
		}
	}
	return blockList, true
}

//GetBlocksByValidatorSearch returns a list of blocks from the database according to the search term and given validator
func GetBlocksByValidatorSearch(term, validator string) ([config.ListSize]*types.StoreBlock, bool) {
	itemList, ok := getItemList(&types.StoreBlock{}, "/api/validatorblocksearch/?term="+term+"&validator="+validator)
	if !ok {
		return [config.ListSize]*types.StoreBlock{}, false
	}
	list, ok := itemList.([config.ListSize]*types.StoreBlock)
	if !ok {
		return [config.ListSize]*types.StoreBlock{}, false
	}
	return list, true
}
