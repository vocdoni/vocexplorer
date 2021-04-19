package api

import (
	"encoding/json"
	"strings"

	types "gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//GetEntityList returns a list of entities from the database
func GetEntityList(i int) ([config.ListSize]string, bool) {
	body, ok := requestBody("/api/listentities/?from=" + util.IntToString(i))
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return [config.ListSize]string{}, false
	}
	var rawEntityList types.ItemList
	err := json.NewDecoder(body).Decode(&rawEntityList)
	if err != nil {
		logger.Error(err)
	}
	var entityList [config.ListSize]string
	for i, rawEntity := range rawEntityList.Items {
		if len(rawEntity) > 0 {
			entity := strings.ToLower(util.HexToString(rawEntity))
			entityList[i] = entity
		}
	}
	return entityList, true
}

//GetEntitySearch returns a list of entities from the database according to the search term
func GetEntitySearch(term string) ([config.ListSize]string, bool) {
	itemList, ok := getItemList("", "/api/entitysearch/?term="+term)
	if !ok {
		return [config.ListSize]string{}, false
	}
	list, ok := itemList.([config.ListSize]string)
	if !ok {
		return [config.ListSize]string{}, false
	}
	return list, true
}
