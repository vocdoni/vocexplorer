package components

import (
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
	vc *client.VochainInfo
	c  *client.Client
}

// Render renders the DashboardView component
func (gw *GatewayView) Render() vecty.ComponentOrHTML {
	// defer func() { go gw.updateGatewayInfo() }()
	return elem.Div(
		renderVochainInfo(gw.vc),
	)
}

func renderVochainInfo(vc *client.VochainInfo) *vecty.HTML {
	if vc != nil && vc.Timestamp != 0 {
		return elem.Div(
			elem.Heading5(vecty.Text("Vochain Info")),
			elem.UnorderedList(
				vecty.If(vc.APIList != nil, elem.ListItem(vecty.Text("API list: "+strings.Join(vc.APIList, ", ")))),
				elem.ListItem(vecty.Text("Blockchain health: "+strconv.Itoa(int(vc.Health)))),
				elem.ListItem(vecty.Text("Ok? "+strconv.FormatBool(vc.Ok))),
				elem.ListItem(vecty.Text("Block height: "+strconv.Itoa(int(vc.Height)))),
				elem.ListItem(vecty.Text("Last block timestamp: "+util.SToTime(vc.BlockTimeStamp))),
				elem.ListItem(vecty.Text("Process ID list: (first 4 id's) "+strings.Join(vc.ProcessIDs[0:4], ", "))),
				elem.ListItem(vecty.Text("Entity ID list: (first 4 id's) "+strings.Join(vc.EntityIDs[0:4], ", "))),
			),
			vecty.If(vc.BlockTime != nil, elem.Table(
				elem.Caption(vecty.Text("Average block times")),
				elem.TableHead(
					elem.TableRow(elem.TableHeader(vecty.Text("Time period")), elem.TableHeader(vecty.Text("Avg time"))),
				),
				elem.TableBody(
					elem.TableRow(elem.TableData(vecty.Text("Last 1m")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[0])))),
					elem.TableRow(elem.TableData(vecty.Text("Last 10m")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[1])))),
					elem.TableRow(elem.TableData(vecty.Text("Last 1h")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[2])))),
					elem.TableRow(elem.TableData(vecty.Text("Last 6h")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[3])))),
					elem.TableRow(elem.TableData(vecty.Text("Last 24h")), elem.TableData(vecty.Text(util.MsToString(vc.BlockTime[4])))),
				),
			)),
			elem.Footer(vecty.Text("Last updated: "+util.SToTime(vc.Timestamp))),
		)
	}
	return vecty.Text("Waiting for blockchain info...")
}

// InitGatewayView connects to gateway websocket and returns a GatewayView component
func initGatewayView(vc *client.VochainInfo) *GatewayView {
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
	gwView.vc = vc
	go updateAndRender(&gwView)
	return &gwView
}

func updateAndRender(gw *GatewayView) {
	defer gw.c.Close()
	for js.Global().Get("gateway").Bool() {
		fmt.Println("Getting vochain info")
		client.UpdateVochainInfo(gw.c, gw.vc)
		time.Sleep(5 * time.Second)
		vecty.Rerender(gw)
	}
}
