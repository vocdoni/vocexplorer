package router

import (
	"encoding/json"
	"net/http"

	"github.com/vocdoni/vocexplorer/api"
	"github.com/vocdoni/vocexplorer/config"
	"github.com/vocdoni/vocexplorer/db"
	ptypes "github.com/vocdoni/vocexplorer/proto"
	"github.com/vocdoni/vocexplorer/util"
	"go.vocdoni.io/dvote/log"
	"google.golang.org/protobuf/proto"
)

// GetProcessHandler writes a single process
func GetProcessHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByIDHandler(d,
		"id",
		config.ProcessIDPrefix,
		func(key []byte) ([]byte, error) {
			var height ptypes.Height
			err := proto.Unmarshal(key, &height)
			if err != nil {
				return []byte{}, err
			}
			rawItem, err := d.Db.Get(append([]byte(config.ProcessHeightPrefix), util.EncodeInt(height.GetHeight())...))
			if err != nil {
				return []byte{}, err
			}
			return rawItem, nil
		},
		packProcess)
}

// ListProcessesHandler writes a list of processes from 'from'
func ListProcessesHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(d, config.ProcessHeightPrefix, nil, packProcess)
}

// ListProcessesByEntityHandler writes a list of processes belonging to 'entity'
func ListProcessesByEntityHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsByParent(d, "entity", config.EntityProcessCountMapKey, config.ProcessByEntityPrefix, config.ProcessHeightPrefix, true, packProcess)
}

// ProcessHeightByEntityHandler writes the number of processes which share the given entity
func ProcessHeightByEntityHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildHeightByParentHandler(d, "entity", config.EntityProcessCountMapKey)
}

// SearchProcessesHandler writes a list of processes by search term
func SearchProcessesHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildSearchHandler(d,
		config.ProcessIDPrefix,
		false,
		func(key []byte) ([]byte, error) {
			height := &ptypes.Height{}
			err := proto.Unmarshal(key, height)
			if err != nil {
				log.Warn("Unable to unmarshal process height")
			}
			return d.Db.Get(append([]byte(config.ProcessHeightPrefix), util.EncodeInt(height.GetHeight())...))
		},
		packProcess,
	)
}

// GetProcessResultsHandler writes the given process' results
func GetProcessResultsHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, ok := r.URL.Query()["id"]
		if !ok || len(ids[0]) < 1 {
			log.Warnf("Url Param 'id' is missing")
			http.Error(w, "Url Param 'id' missing", http.StatusBadRequest)
			return
		}
		id := ids[0]
		t, state, results, err := d.Vs.GetProcessResults(id)
		if err != nil {
			log.Warn(err)
			// http.Error(w, "Cannot get results for process "+id, http.StatusInternalServerError)
			// return
		}
		json.NewEncoder(w).Encode(&api.ProcessResults{
			Type:    t,
			State:   state,
			Results: results,
		})
	}
}

// GetProcessKeysHandler writes the given process' keys
func GetProcessKeysHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, ok := r.URL.Query()["id"]
		if !ok || len(ids[0]) < 1 {
			log.Warnf("Url Param 'id' is missing")
			http.Error(w, "Url Param 'id' missing", http.StatusBadRequest)
			return
		}
		id := ids[0]
		keys, err := d.Vs.GetProcessKeys(id)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Cannot get keys for process "+id, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(keys)
	}
}

func packProcess(raw []byte) []byte {
	var item ptypes.Process
	err := proto.Unmarshal(raw, &item)
	if err != nil {
		log.Error(err)
	}
	new, err := json.Marshal(item.Mirror())
	if err != nil {
		log.Error(err)
	}
	return new
}
