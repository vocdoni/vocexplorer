package components

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ProcessesDashboardView renders the processes dashboard page
type ProcessesDashboardView struct {
	vecty.Core
	vc       *client.VochainInfo
	gwClient *client.Client
}

// Render renders the ProcessesDashboardView component
func (dash *ProcessesDashboardView) Render() vecty.ComponentOrHTML {
	if dash != nil && dash.gwClient != nil && dash.vc != nil {
		return elem.Div(
			elem.Main(
				vecty.Markup(vecty.Class("info-pane")),
				&VochainInfoView{
					vc: dash.vc,
				},
			),
			vecty.Markup(
				event.BeforeUnload(func(i *vecty.Event) {
					fmt.Println("Unloading page")
					js.Global().Get("alert").Invoke("Closing page")
					dash.gwClient.Close()
				}),
			),
		)
	}
	return vecty.Text("Connecting to blockchain clients")
}

func initProcessesDashboardView(vc *client.VochainInfo, ProcessesDashboardView *ProcessesDashboardView) *ProcessesDashboardView {
	js.Global().Set("apiEnabled", true)
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

	ProcessesDashboardView.gwClient = gwClient
	ProcessesDashboardView.vc = vc
	go updateAndRenderProcessesDashboard(ProcessesDashboardView)
	return ProcessesDashboardView
}

func updateAndRenderProcessesDashboard(d *ProcessesDashboardView) {
	defer d.gwClient.Close()
	first := true
	for js.Global().Get("apiEnabled").Bool() {
		client.UpdateProcessesDashboardInfo(d.gwClient, d.vc)
		if first {
			first = false
		} else {
			time.Sleep(config.RefreshTime * time.Second)
		}
		vecty.Rerender(d)
	}
	fmt.Println("Closing gateway updater")
}
