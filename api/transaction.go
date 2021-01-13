package api

import (
	"encoding/json"

	types "github.com/vocdoni/vocexplorer/api/dbtypes"
	"github.com/vocdoni/vocexplorer/config"
	"github.com/vocdoni/vocexplorer/logger"
	"github.com/vocdoni/vocexplorer/util"
)

//GetTxList returns a list of transactions from the database
func GetTxList(from int) ([config.ListSize]*types.Transaction, bool) {
	body, ok := requestBody("/api/listtxs/?from=" + util.IntToString(from))
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return [config.ListSize]*types.Transaction{}, false
	}
	var rawTxList types.ItemList
	err := json.NewDecoder(body).Decode(&rawTxList)
	if err != nil {
		logger.Error(err)
	}
	var txList [config.ListSize]*types.Transaction
	for i, rawTx := range rawTxList.Items {
		if len(rawTx) > 0 {
			var tx types.Transaction
			err = json.Unmarshal(rawTx, &tx)
			if err != nil {
				logger.Error(err)
			}
			txList[i] = &tx
		}
	}
	return txList, true
}

//GetTxByHash returns a transaction from the database
func GetTxByHash(hash string) (*types.Transaction, bool) {
	body, ok := requestBody("/api/txbyhash/?id=" + hash)
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return &types.Transaction{}, false
	}
	var tx types.Transaction
	err := json.NewDecoder(body).Decode(&tx)
	if err != nil {
		logger.Error(err)
	}
	return &tx, true
}

//GetTxByHeight returns a transaction from the database
func GetTxByHeight(height int64) (*types.Transaction, bool) {
	body, ok := requestBody("/api/txbyheight/?id=" + util.IntToString(height))
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return &types.Transaction{}, false
	}
	var tx types.Transaction
	err := json.NewDecoder(body).Decode(&tx)
	if err != nil {
		logger.Error(err)
	}
	return &tx, true
}

//GetTxHeightFromHash finds the height corresponding to a given tx hash
func GetTxHeightFromHash(hash string) (int64, bool) {
	body, ok := requestBody("/api/txhash/?hash=" + hash)
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return 0, false
	}
	var height types.Height
	err := json.NewDecoder(body).Decode(&height)
	if err != nil {
		logger.Error(err)
	}
	return height.Height, true
}

//GetTransactionSearch returns a list of transactions from the database according to the search term
func GetTransactionSearch(term string) ([config.ListSize]*types.Transaction, bool) {
	itemList, ok := getItemList(&types.Transaction{}, "/api/transactionsearch/?term="+term)
	if !ok {
		return [config.ListSize]*types.Transaction{}, false
	}
	list, ok := itemList.([config.ListSize]*types.Transaction)
	if !ok {
		return [config.ListSize]*types.Transaction{}, false
	}
	return list, true
}
