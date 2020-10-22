package api

import (
	"encoding/json"

	types "gitlab.com/vocdoni/vocexplorer/api/dbtypes"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//GetEnvelope gets a single envelope by global height
func GetEnvelope(height int64) (*types.Envelope, bool) {
	body, ok := requestBody("/api/envelope/?height=" + util.IntToString(height))
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return &types.Envelope{}, false
	}
	envelope := new(types.Envelope)
	err := json.NewDecoder(body).Decode(envelope)
	if err != nil {
		logger.Error(err)
	}
	return envelope, true
}

//GetEnvelopeHeightFromNullifier finds the height corresponding to a given envelope nullifier
func GetEnvelopeHeightFromNullifier(hash string) (int64, bool) {
	body, ok := requestBody("/api/envelopenullifier/?nullifier=" + hash)
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

//GetEnvelopeList returns a list of envelopes from the database
func GetEnvelopeList(i int) ([config.ListSize]*types.Envelope, bool) {
	body, ok := requestBody("/api/listenvelopes/?from=" + util.IntToString(i))
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return [config.ListSize]*types.Envelope{}, false
	}
	var rawEnvList types.ItemList
	err := json.NewDecoder(body).Decode(&rawEnvList)
	if err != nil {
		logger.Error(err)
	}
	var envList [config.ListSize]*types.Envelope
	for i, rawEnvelope := range rawEnvList.Items {
		if len(rawEnvelope) > 0 {
			envelope := new(types.Envelope)
			err = json.Unmarshal(rawEnvelope, envelope)
			envList[i] = envelope
			if err != nil {
				logger.Error(err)
			}
		}
	}
	return envList, true
}

//GetEnvelopeSearch returns a list of envelopes from the database according to the search term
func GetEnvelopeSearch(term string) ([config.ListSize]*types.Envelope, bool) {
	itemList, ok := getItemList(&types.Envelope{}, "/api/envelopesearch/?term="+term)
	if !ok {
		return [config.ListSize]*types.Envelope{}, false
	}
	list, ok := itemList.([config.ListSize]*types.Envelope)
	if !ok {
		return [config.ListSize]*types.Envelope{}, false
	}
	return list, true
}

//GetEnvelopeListByProcess returns a list of envelopes by process
func GetEnvelopeListByProcess(i int, process string) ([config.ListSize]*types.Envelope, bool) {
	body, ok := requestBody("/api/listenvelopesprocess/?from=" + util.IntToString(i) + "&process=" + process)
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return [config.ListSize]*types.Envelope{}, false
	}
	var rawEnvList types.ItemList
	err := json.NewDecoder(body).Decode(&rawEnvList)
	if err != nil {
		logger.Error(err)
	}
	var envList [config.ListSize]*types.Envelope
	for i, rawEnvelope := range rawEnvList.Items {
		if len(rawEnvelope) > 0 {
			envelope := new(types.Envelope)
			err = json.Unmarshal(rawEnvelope, envelope)
			envList[i] = envelope
			if err != nil {
				logger.Error(err)
			}
		}
	}
	return envList, true
}
