package config

const (
	// GatewayHost is gateway websockets url
	GatewayHost = "ws://0.0.0.0:9090/dvote"
	// TendermintHost is tendermint api url
	TendermintHost = "http://0.0.0.0:26657"
	// RefreshTime is the number of seconds between page data refresh
	RefreshTime = 5
	// SearchPageSmall is the of elements on a small search page
	SearchPageSmall = 10
)
