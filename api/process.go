package api

import (
	"encoding/json"

	types "github.com/vocdoni/vocexplorer/api/dbtypes"
	"github.com/vocdoni/vocexplorer/config"
	"github.com/vocdoni/vocexplorer/logger"
	"github.com/vocdoni/vocexplorer/util"
)

// GetProcess fetches a process
func GetProcess(id string) (*types.Process, bool) {
	body, ok := requestBody("/api/process/?id=" + id)
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return &types.Process{}, false
	}
	process := new(types.Process)
	err := json.NewDecoder(body).Decode(&process)
	if err != nil {
		logger.Error(err)
		return process, false
	}
	return process, true
}

//GetProcessList returns a list of entities from the database
func GetProcessList(i int) ([config.ListSize]*types.Process, bool) {
	body, ok := requestBody("/api/listprocesses/?from=" + util.IntToString(i))
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return [config.ListSize]*types.Process{}, false
	}
	var rawProcessList types.ItemList
	err := json.NewDecoder(body).Decode(&rawProcessList)
	if err != nil {
		logger.Error(err)
	}
	var processList [config.ListSize]*types.Process
	for i, rawProcess := range rawProcessList.Items {
		if len(rawProcess) > 0 {
			process := new(types.Process)
			err := json.Unmarshal(rawProcess, &process)
			processList[i] = process
			if err != nil {
				logger.Error(err)
			}
		}
	}
	return processList, true
}

//GetProcessListByEntity returns a list of processes by entity
func GetProcessListByEntity(i int, entity string) ([config.ListSize]*types.Process, bool) {
	body, ok := requestBody("/api/listprocessesbyentity/?from=" + util.IntToString(i) + "&entity=" + entity)
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return [config.ListSize]*types.Process{}, false
	}
	var rawProcessList types.ItemList
	err := json.NewDecoder(body).Decode(&rawProcessList)
	if err != nil {
		logger.Error(err)
	}
	var processList [config.ListSize]*types.Process
	for i, rawProcess := range rawProcessList.Items {
		if len(rawProcess) > 0 {
			process := new(types.Process)
			err := json.Unmarshal(rawProcess, &process)
			processList[i] = process
			if err != nil {
				logger.Error(err)
			}
		}
	}
	return processList, true
}

//GetProcessSearch returns a list of processes from the database according to the search term
func GetProcessSearch(term string) ([config.ListSize]*types.Process, bool) {
	itemList, ok := getItemList(&types.Process{}, "/api/processsearch/?term="+term)
	if !ok {
		return [config.ListSize]*types.Process{}, false
	}
	list, ok := itemList.([config.ListSize]*types.Process)
	if !ok {
		return [config.ListSize]*types.Process{}, false
	}
	return list, true
}

// GetProcessResults fetches the results of a process
func GetProcessResults(id string) (*ProcessResults, bool) {
	body, ok := requestBody("/api/processresults/?id=" + id)
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return &ProcessResults{}, false
	}
	results := new(ProcessResults)
	err := json.NewDecoder(body).Decode(&results)
	if err != nil {
		logger.Error(err)
		return results, false
	}
	return results, true
}

// GetProcessKeys gets process keys
func GetProcessKeys(pid string) (*Pkeys, bool) {
	body, ok := requestBody("/api/processkeys/?id=" + pid)
	if body != nil {
		defer body.Close()
	}
	if !ok {
		return &Pkeys{}, false
	}
	keys := new(Pkeys)
	err := json.NewDecoder(body).Decode(&keys)
	if err != nil {
		logger.Error(err)
		return keys, false
	}
	return keys, true
}
