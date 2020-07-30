package config

//Cfg is the global config to be served to pages
type Cfg struct {
	// GatewayHost is gateway websockets url
	GatewayHost string `json:"gatewayHost"`
	// TendermintHost is tendermint api url
	TendermintHost string `json:"tendermintHost"`
	// RefreshTime is the number of seconds between page data refresh
	RefreshTime int `json:"refreshTime"`
}

const (
	//ListSize is the number of cards shown in list of blocks/processes/etc
	ListSize = 4
	//ProcessIDPrefix is the key prefix for process id's
	ProcessIDPrefix = "00"
	//EntityIDPrefix is the key prefix for entity id's
	EntityIDPrefix = "01"
	//BlockPrefix is the key prefix for blocks
	BlockPrefix = "02"
	//LatestBlockHeightKey is the key for the value of the latest block height stored
	LatestBlockHeightKey = "LatestBlockHeight"
)
