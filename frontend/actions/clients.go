package actions

import (
	"gitlab.com/vocdoni/vocexplorer/api"
	"nhooyr.io/websocket"
)

//TendermintClientInit is the action to initialize the global tendermint rpc client
type TendermintClientInit struct {
	Client *websocket.Conn
}

//GatewayClientInit is the action to initialize the global gateway websockets client
type GatewayClientInit struct {
	Client *api.GatewayClient
}

//GatewayConnected is the action to change the connection status of the gateway
type GatewayConnected struct {
	Connected bool
}

//ServerConnected is the action to change the connection status of the web server
type ServerConnected struct {
	Connected bool
}
