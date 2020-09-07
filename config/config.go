package config

//Cfg is the global config to be served to pages
type Cfg struct {
	// Detached is the detached state of the database
	Detached bool `json:"detached"`
	// GatewayHost is gateway websockets url
	GatewayHost string `json:"gatewayHost"`
	// TendermintHost is tendermint api url
	TendermintHost string `json:"tendermintHost"`
	// RefreshTime is the number of seconds between page data refresh
	RefreshTime int `json:"refreshTime"`
}

//MainCfg includes backend and frontend config
type MainCfg struct {
	DataDir     string
	DisableGzip bool
	Global      Cfg
	HostURL     string
	ChainID     string
	LogLevel    string
}

const (
	//ListSize is the number of cards shown in list of blocks/processes/etc
	ListSize = 10
	//HomeWidgetBlocksListSize is the number of blocks shown on the homepage
	HomeWidgetBlocksListSize = 4
	//NumBlockUpdates is the number of blocks updated per db batch
	NumBlockUpdates = 100
	//DBWaitTime is the number of seconds the backend waits before batching another set of blocks
	DBWaitTime = 0
	//ProcessIDPrefix is the key prefix for process id's by height
	ProcessIDPrefix = "ph_"
	//EntityIDPrefix is the key prefix for entity id's by height
	EntityIDPrefix = "eh_"
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
)
