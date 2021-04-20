package config

//Cfg is the global config to be served to pages
type Cfg struct {
	// RefreshTime is the number of seconds between page data refresh
	RefreshTime int    `json:"refreshTime"`
	Environment string `json:"environment"`
	GatewayUrl  string `json:"gatewayUrl"`
}

//MainCfg includes backend and frontend config
type MainCfg struct {
	DataDir     string
	DisableGzip bool
	Global      Cfg
	HostURL     string
	LogLevel    string
}

const (
	//ListSize is the number of cards shown in list of blocks/processes/etc
	ListSize = 10
	//MaxListSize is the largest number of elements in a list
	MaxListSize = 1 << 32
	//HomeWidgetBlocksListSize is the number of blocks shown on the homepage
	HomeWidgetBlocksListSize = 4
	//DefaultNamespace is the default namespace value to get all processes
	DefaultNamespace = uint32(0)
)
