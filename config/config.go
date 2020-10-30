package config

import dvotecfg "gitlab.com/vocdoni/go-dvote/config"

//Cfg is the global config to be served to pages
type Cfg struct {
	// RefreshTime is the number of seconds between page data refresh
	RefreshTime int `json:"refreshTime"`
}

//MainCfg includes backend and frontend config
type MainCfg struct {
	Chain         string
	DataDir       string
	DisableGzip   bool
	Global        Cfg
	HostURL       string
	LogLevel      string
	VochainConfig *dvotecfg.VochainCfg
}

const (
	//ListSize is the number of cards shown in list of blocks/processes/etc
	ListSize = 10
	//MaxListSize is the largest number of elements in a list
	MaxListSize = 1 << 32
	//HomeWidgetBlocksListSize is the number of blocks shown on the homepage
	HomeWidgetBlocksListSize = 4
	//NumBlockUpdates is the number of blocks updated per db batch
	NumBlockUpdates = 100
	//DBWaitTime is the number of milliseconds the backend waits before batching another set of blocks
	DBWaitTime = 0
	//ProcessHeightPrefix is the key prefix for processes by height
	ProcessHeightPrefix = "ph_"
	//EntityHeightPrefix is the key prefix for entity id's by height
	EntityHeightPrefix = "eh_"
	//BlockHeightPrefix is the key prefix for block hashes by height
	BlockHeightPrefix = "bh_"
	//BlockHashPrefix is the key prefix for blocks by hash
	BlockHashPrefix = "bid"
	//TxHeightPrefix is the key prefix for transaction hashes by height
	TxHeightPrefix = "th_"
	//TxHashPrefix is the key prefix for transactions by hash
	TxHashPrefix = "tid"
	//ValidatorPrefix is the key prefix for validators by address
	ValidatorPrefix = "vid"
	//EnvPackagePrefix is the key prefix for envelope packages by height
	EnvPackagePrefix = "evh"
	//EnvNullifierPrefix is the key prefix for envelope heights by nullifier
	EnvNullifierPrefix = "evi"
	//EnvPIDPrefix is the key prefix for envelope heights by processID
	EnvPIDPrefix = "evp"
	//BlockByValidatorPrefix is the key prefix for block hash by validator
	BlockByValidatorPrefix = "bv_"
	//ProcessByEntityPrefix is the key prefix for process heights by entity process height
	ProcessByEntityPrefix = "pe_"
	//ValidatorHeightPrefix is the key prefix for validator IDs by height
	ValidatorHeightPrefix = "vh_"
	//ProcessIDPrefix is the key prefix for process heights by ID
	ProcessIDPrefix = "pid"
	//EntityIDPrefix is the key prefix for entity ID's as standalone keys (for iteration over keys)
	EntityIDPrefix = "eid"
	//BlockchainInfoKey is the key for the blockchain info storage
	BlockchainInfoKey = "BlockchainInfo"
	//LatestBlockHeightKey is the key for the value of the latest block height stored
	LatestBlockHeightKey = "LatestBlockHeight"
	//LatestTxHeightKey is the key for the value of the latest tx height stored
	LatestTxHeightKey = "LatestTxHeight"
	//LatestEnvelopeCountKey is the key for the value of the latest envelope count stored
	LatestEnvelopeCountKey = "LatestEnvHeight"
	//LatestEntityCountKey is the key for the value of the latest entity count stored
	LatestEntityCountKey = "LatestEntityCountKey"
	//LatestProcessCountKey is the key for the value of the latest process count stored
	LatestProcessCountKey = "LatestProcessCountKey"
	//LatestValidatorCountKey is the key for the value of the latest validator count stored
	LatestValidatorCountKey = "LatestValHeight"
	//EntityProcessCountMapKey is the key for the map of entity process counts
	EntityProcessCountMapKey = "EntityProcHeight"
	//ProcessEnvelopeCountMapKey is the key for the map of process envelope counts
	ProcessEnvelopeCountMapKey = "ProcEnvHeight"
	//ValidatorHeightMapKey is the key for the map of validator block heights
	ValidatorHeightMapKey = "ValHeightMap"
	//MaxTxsPerBlockKey is the key for the maximum number of txs on one block
	MaxTxsPerBlockKey = "MaxTxs"
	//MaxTxsBlockIDKey is the key for the block ID with the largest nubmer of txs
	MaxTxsBlockIDKey = "MaxBlock"
	//MaxTxsBlockHeightKey is the key for the block Height with the largest nubmer of txs
	MaxTxsBlockHeightKey = "MaxBlockHeight"
	// MaxTxsPerMinuteKey is the key for the maximum number of transactions in one minute of time
	MaxTxsPerMinuteKey = "MaxTxsMinute"
	// MaxTxsMinuteID is the unix code for the start of the minute with the maximum number of txs
	MaxTxsMinuteID = "MaxMinute"
	// GlobalProcessListKey is the key for the global list of processes
	GlobalProcessListKey = "ProcList"
)
