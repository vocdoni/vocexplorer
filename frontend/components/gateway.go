package components

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// GatewayView renders the gateway info component
type GatewayView struct {
	vecty.Core
	gwInfo *client.MetaResponse
	c      *client.Client
}

// Render renders the DashboardView component
func (gw *GatewayView) Render() vecty.ComponentOrHTML {
	// defer func() { go gw.updateGatewayInfo() }()
	return elem.Div(
		renderGatewayInfo(gw.gwInfo),
	)
}

// Iteratively calls gateway api until "gateway" env variable is set to false.
func (gw *GatewayView) updateGatewayInfo(cancel context.CancelFunc) {
	defer gw.c.Close()
	for js.Global().Get("gateway").Bool() {
		fmt.Println("Getting info")
		resp, err := gw.c.GetGatewayInfo()
		if util.ErrPrint(err) {
			fmt.Println("Unable to get gateway info")
		}
		gw.gwInfo = resp
		fmt.Println("body")
		vecty.Rerender(gw)
		time.Sleep(5 * time.Second)
	}
}

func renderGatewayInfo(info *client.MetaResponse) *vecty.HTML {
	if info != nil && info.Timestamp != 0 {
		return elem.Div(
			elem.Heading5(vecty.Text("Gateway Info")),
			elem.UnorderedList(
				vecty.If(info.APIList != nil, elem.ListItem(vecty.Text("API list: "+strings.Join(info.APIList, ", ")))),
				elem.ListItem(vecty.Text("Blockchain health: "+strconv.Itoa(int(info.Health)))),
				elem.ListItem(vecty.Text("Ok? "+strconv.FormatBool(info.Ok))),
				elem.ListItem(vecty.Text("Request: "+info.Request)),
				elem.ListItem(vecty.Text("Timestamp: "+strconv.Itoa(int(info.Timestamp)))),
			),
		)
	}
	return vecty.Text("Waiting for blockchain info...")
}

// InitGatewayView connects to gateway websocket and returns a GatewayView component
func initGatewayView(d *client.MetaResponse) *GatewayView {
	js.Global().Set("gateway", true)
	// Establishing connection with gateway host
	fmt.Println("connecting to %s", config.GatewayHost)
	gw, cancel, err := client.New(config.GatewayHost)
	defer cancel()
	if err != nil {
		js.Global().Get("alert").Invoke("Unable to connect to Gateway client. Please see readme file")
		return nil
	}
	var gwView GatewayView
	gwView.c = gw
	gwView.gwInfo = d
	go (&gwView).updateGatewayInfo(cancel)
	return &gwView
}
