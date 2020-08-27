package config

//Cfg is the global config to be served to pages
type Cfg struct {
	// GatewayHost is gateway websockets url
	GatewayHost string `json:"gatewayHost"`
	// GatewaySocket is gateway websockets socket label
	GatewaySocket string `json:"gatewaySocket"`
	// TendermintHost is tendermint api url
	TendermintHost string `json:"tendermintHost"`
	// RefreshTime is the number of seconds between page data refresh
	RefreshTime int `json:"refreshTime"`
}

//MainCfg includes backend and frontend config
type MainCfg struct {
	DataDir     string
	Detached    bool
	DisableGzip bool
	Global      Cfg
	HostURL     string
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
	ProcessIDPrefix = "00"
	//EntityIDPrefix is the key prefix for entity id's by height
	EntityIDPrefix = "01"
	//BlockHeightPrefix is the key prefix for block hashes by height
	BlockHeightPrefix = "02"
	//BlockHashPrefix is the key prefix for blocks by hash
	BlockHashPrefix = "03"
	//TxHeightPrefix is the key prefix for transaction hashes by height
	TxHeightPrefix = "04"
	//TxHashPrefix is the key prefix for transactions by hash
	TxHashPrefix = "05"
	//ValidatorPrefix is the key prefix for validators by address
	ValidatorPrefix = "06"
	//EnvPackagePrefix is the key prefix for envelope packages by height
	EnvPackagePrefix = "07"
	//EnvNullifierPrefix is the key prefix for envelope heights by nullifier
	EnvNullifierPrefix = "08"
	//EnvPIDPrefix is the key prefix for envelope heights by processID
	EnvPIDPrefix = "09"
	//BlockByValidatorPrefix is the key prefix for block hash by validator
	BlockByValidatorPrefix = "10"
	//ProcessByEntityPrefix is the key prefix for process heights by entity process height
	ProcessByEntityPrefix = "11"
	//ValidatorHeightPrefix is the key prefix for validator IDs by height
	ValidatorHeightPrefix = "12"
	//LatestBlockHeightKey is the key for the value of the latest block height stored
	LatestBlockHeightKey = "LatestBlockHeight"
	//LatestTxHeightKey is the key for the value of the latest tx height stored
	LatestTxHeightKey = "LatestTxHeight"
	//LatestValidatorHeightKey is the key for the value of the latest validator height stored
	LatestValidatorHeightKey = "LatestValHeight"
	//ValidatorHeightMapKey is the key for the map of validator block heights
	ValidatorHeightMapKey = "ValHeightMap"
	//ProcessEnvelopeHeightMapKey is the key for the map of process envelope heights
	ProcessEnvelopeHeightMapKey = "ProcEnvHeight"
	//LatestEnvelopeHeightKey is the key for the value of the latest envelope height stored
	LatestEnvelopeHeightKey = "LatestEnvHeight"
	//LatestEntityHeight is the key for the value of the latest entity height stored
	LatestEntityHeight = "LatestEntityHeight"
	//LatestProcessHeight is the key for the value of the latest process height stored
	LatestProcessHeight = "LatestProcessHeight"
	//EntityProcessHeightMapKey is the key for the map of entity process heights
	EntityProcessHeightMapKey = "EntityProcHeight"
)
