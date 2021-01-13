package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/vocdoni/vocexplorer/frontend/bootstrap"
	"github.com/vocdoni/vocexplorer/frontend/store"
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
			Contents: "Connecting to " + b.connection,
		},
	)
}

func renderServerConnectionBanner() vecty.ComponentOrHTML {
	if !store.ServerConnected {
		return &ConnectedBanner{
			connection: "web server",
		}
	}
	return nil
}
