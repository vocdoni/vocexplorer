package router

import (
	"encoding/json"
	"net/http"

	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/db"
	ptypes "gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
	"google.golang.org/protobuf/proto"
)

// GetEnvelopeHandler writes a single envelope
func GetEnvelopeHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildItemByHeightHandler(d, config.LatestEnvelopeCountKey, config.EnvPackagePrefix, nil, packEnvelope)
}

// ListEnvelopesByProcessHandler writes a list of envelopes which share the given process
func ListEnvelopesByProcessHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsByParent(d, "process", config.ProcessEnvelopeCountMapKey, config.EnvPIDPrefix, config.EnvPackagePrefix, true, packEnvelope, false)
}

// ListEnvelopesHandler writes a list of envelopes
func ListEnvelopesHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildListItemsHandler(d, config.EnvPackagePrefix, nil, packEnvelope)
}

// EnvelopeHeightByProcessHandler writes the number of envelopes which share the given processID
func EnvelopeHeightByProcessHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildHeightByParentHandler(d, "process", config.ProcessEnvelopeCountMapKey)
}

// EnvelopeHeightFromNullifierHandler returns the height of the corresponding envelope nullifier
func EnvelopeHeightFromNullifierHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, ok := r.URL.Query()["nullifier"]
		if !ok || len(ids[0]) < 1 {
			log.Warnf("Url Param 'nullifier' is missing")
			http.Error(w, "Url Param 'nullifier' missing", http.StatusNotFound)
			return
		}
		id := ids[0]
		hash := []byte(id)
		key := append([]byte(config.EnvNullifierPrefix), hash...)
		has, err := d.Db.Has(key)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to get envelope height", http.StatusInternalServerError)
			return
		}
		if !has {
			log.Warnf("Envelope nullifier does not exist")
			http.Error(w, "Envelope nullifier does not exist", http.StatusInternalServerError)
			return
		}
		raw, err := d.Db.Get(key)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to get envelope key", http.StatusInternalServerError)
			return
		}

		var height ptypes.Height
		err = proto.Unmarshal(raw, &height)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Unable to unmarshal height", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(height.Mirror())
	}
}

// SearchEnvelopesHandler writes a list of envelopes by search term
func SearchEnvelopesHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return buildSearchHandler(d,
		config.EnvNullifierPrefix,
		false,
		func(key []byte) ([]byte, error) {
			height := &ptypes.Height{}
			err := proto.Unmarshal(key, height)
			if err != nil {
				log.Warn("Unable to unmarshal envelope height")
			}
			return d.Db.Get(append([]byte(config.EnvPackagePrefix), util.EncodeInt(height.GetHeight())...))
		},
		packEnvelope, false,
	)
}

func packEnvelope(raw []byte) []byte {
	var item ptypes.Envelope
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
