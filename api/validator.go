package api

import (
	"encoding/json"

	types "gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//GetValidatorSearch returns a list of validators from the database according to the search term
func GetValidatorSearch(term string) ([config.ListSize]*types.Validator, bool) {
	itemList, ok := getItemList(&types.Validator{}, "/api/validatorsearch/?term="+term)
	if !ok {
		return [config.ListSize]*types.Validator{}, false
	}
	list, ok := itemList.([config.ListSize]*types.Validator)
	if !ok {
		return [config.ListSize]*types.Validator{}, false
	}
	return list, true
}

//GetValidatorList returns a list of validators from the database
func GetValidatorList(i int) ([config.ListSize]*types.Validator, bool) {
	body, ok := requestBody("/api/listvalidators/?from=" + util.IntToString(i))
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return [config.ListSize]*types.Validator{}, false
	}
	var rawValidatorList types.ItemList
	err := json.NewDecoder(body).Decode(&rawValidatorList)
	if err != nil {
		logger.Error(err)
	}
	var validatorList [config.ListSize]*types.Validator
	for i, rawVal := range rawValidatorList.Items {
		if len(rawVal) > 0 {
			var validator types.Validator
			err = json.Unmarshal(rawVal, &validator)
			validatorList[i] = &validator
			if err != nil {
				logger.Error(err)
			}
		}
	}
	return validatorList, true
}

//GetValidator returns a single validator from the database
func GetValidator(address string) (*types.Validator, bool) {
	body, ok := requestBody("/api/validator/?id=" + address)
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return &types.Validator{}, false
	}
	var validator types.Validator
	err := json.NewDecoder(body).Decode(&validator)
	if err != nil {
		logger.Error(err)
	}
	return &validator, true
}
