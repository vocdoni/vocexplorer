package components

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// DashboardView renders the dashboard landing page
type DashboardView struct {
	vecty.Core
	t        *rpc.TendermintInfo
	vc       *client.VochainInfo
	gwClient *client.Client
	tClient  *http.HTTP
}

// Render renders the DashboardView component
func (dash *DashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && dash.gwClient != nil && dash.tClient != nil && dash.t != nil && dash.vc != nil {
		return elem.Div(
			elem.Main(
				vecty.Markup(vecty.Class("info-pane")),
				&StatsView{
					t:  dash.t,
					vc: dash.vc,
				},
			),
			vecty.Markup(
				event.BeforeUnload(func(i *vecty.Event) {
					js.Global().Get("alert").Invoke("Closing page")
					dash.gwClient.Close()
				},
				),
			),
		)
	}
	return vecty.Text("Connecting to blockchain clients")
}

func initDashboardView(t *rpc.TendermintInfo, vc *client.VochainInfo, DashboardView *DashboardView) *DashboardView {
	js.Global().Set("apiEnabled", true)
	// Init tendermint client
	fmt.Println("connecting to %s", config.TendermintHost)
	tClient, err := rpc.InitClient()
	if err != nil {
		if js.Global().Get("confirm").Invoke("Unable to connect to Tendermint client. Reload with client running").Bool() {
			js.Global().Get("location").Call("reload")
		}
		return nil
	}
	// Init Gateway client
	fmt.Println("connecting to %s", config.GatewayHost)
	gwClient, cancel, err := client.New(config.GatewayHost)
	defer cancel()
	if util.ErrPrint(err) {
		if js.Global().Get("confirm").Invoke("Unable to connect to Gateway client. Reload with client running").Bool() {
			js.Global().Get("location").Call("reload")
		}
		return nil
	}

	// var t *rpc.TendermintInfo
	DashboardView.tClient = tClient
	DashboardView.gwClient = gwClient
	DashboardView.t = t
	DashboardView.vc = vc
	go updateAndRenderDashboard(DashboardView)
	return DashboardView
}

func updateAndRenderDashboard(d *DashboardView) {
	defer d.gwClient.Close()
	first := true
	for js.Global().Get("apiEnabled").Bool() {
		rpc.UpdateTendermintInfo(d.tClient, d.t)
		client.UpdateDashboardInfo(d.gwClient, d.vc)
		if first {
			first = false
		} else {
			time.Sleep(config.RefreshTime * time.Second)
		}
		vecty.Rerender(d)
	}
	fmt.Println("Closing tendermint/gateway updater")
}
