package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

//ConnectedBanner is the component to display a banner if server is disconnected
type ConnectedBanner struct {
	vecty.Core
	connection string
}

//Render renders the ConnectedBanner component
func (b *ConnectedBanner) Render() vecty.ComponentOrHTML {
	return elem.Div(
		&bootstrap.Alert{
			Type:     "warning",
			Contents: "Disconnected from " + b.connection,
		},
	)
}

func renderGatewayConnectionBanner() vecty.ComponentOrHTML {
	if !store.GatewayConnected {
		return &ConnectedBanner{
			connection: "blockchain Gateway",
		}
	}
	return nil
}
func renderServerConnectionBanner() vecty.ComponentOrHTML {
	if !store.ServerConnected {
		return &ConnectedBanner{
			connection: "web server",
		}
	}
	return nil
}
