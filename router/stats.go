package router

import (
	"encoding/json"
	"net/http"
	"time"

	prototypes "github.com/golang/protobuf/ptypes"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/db"
	ptypes "gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
	"google.golang.org/protobuf/proto"
)

//PingHandler responds to a ping
func PingHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}

// StatsHandler is the public api for all blockchain statistics & information
func StatsHandler(d *db.ExplorerDB, cfg *config.Cfg) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("Serving statistics ")
		stats := new(api.VochainStats)
		blockchainInfo := new(ptypes.BlockchainInfo)

		//get blockchainInfo
		rawBlockchainInfo, err := d.Db.Get([]byte(config.BlockchainInfoKey))
		if err != nil {
			log.Warn(err)
		}
		err = proto.Unmarshal(rawBlockchainInfo, blockchainInfo)
		if err != nil {
			log.Warn(err)
		}
		genesisTime, err := prototypes.Timestamp(blockchainInfo.GetGenesisTimeStamp())
		if err != nil {
			log.Warn(err)
			genesisTime = time.Unix(1, 0)
		}
		var blockTime [5]int32
		copy(blockTime[:], blockchainInfo.GetBlockTime())

		blockHeight := db.GetHeight(d.Db, config.LatestBlockHeightKey, 1)
		entityCount := db.GetHeight(d.Db, config.LatestEntityCountKey, 0)
		envelopeCount := db.GetHeight(d.Db, config.LatestEnvelopeCountKey, 0)
		processCount := db.GetHeight(d.Db, config.LatestProcessCountKey, 0)
		transactionHeight := db.GetHeight(d.Db, config.LatestTxHeightKey, 1)
		validatorCount := db.GetHeight(d.Db, config.LatestValidatorCountKey, 0)
		maxTxsPerBlock := db.GetInt64(d.Db, config.MaxTxsPerBlockKey)
		maxTxsPerMinute := db.GetInt64(d.Db, config.MaxTxsPerMinuteKey)
		maxTxsBlockHeight := db.GetInt64(d.Db, config.MaxTxsBlockHeightKey)
		rawMaxTxsBlockID, err := d.Db.Get([]byte(config.MaxTxsBlockIDKey))
		var maxTxsBlockID string
		if err != nil {
			maxTxsBlockID = ""
		} else {
			maxTxsBlockID = util.HexToString(rawMaxTxsBlockID)
		}
		rawMaxTxsMinute := db.GetInt64(d.Db, config.MaxTxsMinuteID)
		maxTxsMinute := time.Unix(rawMaxTxsMinute, 0)

		stats.Network = blockchainInfo.GetNetwork()
		stats.Version = blockchainInfo.GetVersion()
		stats.LatestBlockHeight = blockchainInfo.GetLatestBlockHeight()
		stats.GenesisTimeStamp = genesisTime
		stats.ChainID = blockchainInfo.GetChainID()
		stats.BlockTime = &blockTime
		stats.BlockTimeStamp = blockchainInfo.GetBlockTimeStamp()
		stats.Height = blockchainInfo.GetHeight()
		stats.Syncing = blockchainInfo.GetSyncing()
		stats.MaxBytes = blockchainInfo.GetMaxBytes()

		stats.BlockHeight = blockHeight.GetHeight()
		stats.EntityCount = entityCount.GetHeight()
		stats.EnvelopeCount = envelopeCount.GetHeight()
		stats.ProcessCount = processCount.GetHeight()
		stats.TransactionHeight = transactionHeight.GetHeight()
		stats.ValidatorCount = validatorCount.GetHeight()
		stats.MaxTxsPerBlock = maxTxsPerBlock
		stats.MaxTxsPerMinute = maxTxsPerMinute
		stats.MaxTxsBlockHash = maxTxsBlockID
		stats.MaxTxsBlockHeight = maxTxsBlockHeight
		stats.MaxTxsMinute = maxTxsMinute
		if transactionHeight.GetHeight() > 1 && blockHeight.GetHeight() > 1 {
			stats.AvgTxsPerBlock = float64(transactionHeight.GetHeight()-1) / float64(blockHeight.GetHeight()-1)
		}
		if transactionHeight.GetHeight() > 0 && int64(stats.BlockTimeStamp)-stats.GenesisTimeStamp.Unix() > 0 {
			stats.AvgTxsPerMinute = float64(transactionHeight.GetHeight()-1) / (float64(int64(stats.BlockTimeStamp)-stats.GenesisTimeStamp.Unix()) / float64(60))
		}

		msg, err := json.Marshal(stats)
		if err != nil {
			log.Warn(err)
			log.Debugf("Stats: %+v", stats)
			http.Error(w, "Unable to marshal stats", http.StatusInternalServerError)
			return
		}

		w.Write(msg)
	}
}

// HeightHandler writes the int64 height value corresponding to given key
func HeightHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			log.Warnf("Url Param 'key' is missing")
			http.Error(w, "Url Param 'key' missing", http.StatusBadRequest)
			return
		}
		val, err := d.Db.Get([]byte(keys[0]))
		if err != nil {
			log.Warn(err)
			http.Error(w, "Key not found", http.StatusInternalServerError)
			return
		}

		var height ptypes.Height
		err = proto.Unmarshal(val, &height)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Height not found", http.StatusInternalServerError)
			return
		}
		log.Debug("Sent height " + util.IntToString(height.GetHeight()) + " for key " + keys[0])
		json.NewEncoder(w).Encode(height.Mirror())
	}
}

// HeightMapHandler writes the string:int64 height map corresponding to given key
func HeightMapHandler(d *db.ExplorerDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, ok := r.URL.Query()["key"]
		if !ok || len(keys[0]) < 1 {
			log.Warnf("Url Param 'key' is missing")
			http.Error(w, "Url Param 'key' missing", http.StatusBadRequest)
			return
		}
		val, err := d.Db.Get([]byte(keys[0]))
		if err != nil {
			log.Warnf("Key %s not found: %s", keys[0], err.Error())
			http.Error(w, "Key not found", http.StatusInternalServerError)
			return
		}

		var heightMap ptypes.HeightMap
		err = proto.Unmarshal(val, &heightMap)
		if err != nil {
			log.Warn(err)
			http.Error(w, "Height map not found", http.StatusInternalServerError)
			return
		}
		log.Debug("Sent height map for key " + keys[0])
		json.NewEncoder(w).Encode(heightMap.Mirror())
	}
}
