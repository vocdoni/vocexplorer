package components

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
)

// GatewayView renders the gateway info component
type GatewayView struct {
	vecty.Core
	gwInfo *client.MetaResponse
	c      *client.Client
}

// Render renders the DashboardView component
func (gw *GatewayView) Render() vecty.ComponentOrHTML {
	// After rendering, run gouroutine to call api again
	defer func() { go gw.updateGatewayInfo() }()
	return elem.Div(
		renderGatewayInfo(gw.gwInfo),
	)
}

// Recursively calls gateway api, calls rerender, exits. Rerender starts new update routine
func (gw *GatewayView) updateGatewayInfo() {
	time.Sleep(5 * time.Second)
	fmt.Println("Getting info")
	resp, err := gw.c.GetGatewayInfo()
	if err != nil {
		fmt.Println("Unable to get gateway info")
		fmt.Println(err.Error())
	}
	gw.gwInfo = resp
	fmt.Println("body")
	vecty.Rerender(gw)
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
	gatewayHost := "ws://0.0.0.0:9090/dvote"
	// Establishing connection with gateway host
	fmt.Println("connecting to %s", gatewayHost)
	gw, cancel, err := client.New(gatewayHost)
	defer cancel()
	if err != nil {
		fmt.Println("Unable to connect to gateway")
	}
	// defer gw.Conn.Close(websocket.StatusNormalClosure, "")
	return &GatewayView{
		c:      gw,
		gwInfo: d,
	}
}
