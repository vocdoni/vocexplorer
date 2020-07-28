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
	ListSize = 10
)
