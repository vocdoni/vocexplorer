package components

import (
	"fmt"
	"strings"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
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
	if len(store.GatewayError) > 0 && (strings.Contains(store.GatewayError, "vote") || strings.Contains(store.GatewayError, "indexer") || strings.Contains(store.GatewayError, "client not connected")) {
		return elem.Div(
			&bootstrap.Alert{
				Type:     "warning",
				Contents: fmt.Sprintf("Cannot use %s %s: %s", b.connection, store.Config.GatewayUrl, store.GatewayError),
			},
		)
	} else {
		return elem.Div(
			&bootstrap.Alert{
				Type:     "warning",
				Contents: fmt.Sprintf("Connecting to %s %s", b.connection, store.Client.Address),
			},
		)
	}
}

func renderServerConnectionBanner() vecty.ComponentOrHTML {
	if !store.ServerConnected {
		return &ConnectedBanner{
			connection: "gateway",
		}
	}
	return nil
}
