package actions

import (
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/vocdoni/vocexplorer/frontend/api"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

//TendermintClientInit is the action to initialize the global tendermint rpc client
type TendermintClientInit struct {
	Client *http.HTTP
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

// On initialization, register actions
func init() {
	dispatcher.Register(clientActions)
}

// clientActions is the handler for all connection-related store actions
func clientActions(action interface{}) {
	switch a := action.(type) {
	case *TendermintClientInit:
		store.TendermintClient = a.Client

	case *GatewayClientInit:
		store.GatewayClient = a.Client

	case *GatewayConnected:
		store.GatewayConnected = a.Connected

	case *ServerConnected:
		store.ServerConnected = a.Connected

	default:
		return // don't fire listeners
	}

	store.Listeners.Fire()
}
