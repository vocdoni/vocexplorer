package actions

//TendermintClientInit is the action to initialize the global tendermint rpc client
type TendermintClientInit struct {
}

//GatewayClientInit is the action to initialize the global gateway websockets client
type GatewayClientInit struct {
}

//GatewayConnected is the action to change the connection status of the gateway
type GatewayConnected struct {
	Connected bool
}

//ServerConnected is the action to change the connection status of the web server
type ServerConnected struct {
	Connected bool
}
