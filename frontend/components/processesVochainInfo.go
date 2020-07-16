package components

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// VochainInfoView renders the vochainInfo pane
type VochainInfoView struct {
	vecty.Core
	vc *client.VochainInfo
}

// Render renders the VochainInfoView component
func (b *VochainInfoView) Render() vecty.ComponentOrHTML {
	if b.vc != nil {
		// TODO: place this inside a pageload event?
		if js.Global().Get("search").IsUndefined() || !js.Global().Get("search").Bool() {
			// if !js.Global().Get("search").Bool() {
			b.vc.EntitySearchIDs = util.TrimSlice(b.vc.EntityIDs, 10, true)
			b.vc.ProcessSearchIDs = util.TrimSlice(b.vc.ProcessIDs, 10, true)
		}
		return elem.Section(
			elem.Input(vecty.Markup(
				event.Input(func(e *vecty.Event) {
					search := e.Target.Get("value").String()
					js.Global().Set("search", true)
					if search != "" {
						b.vc.EntitySearchIDs = util.TrimSlice(util.SearchSlice(b.vc.EntityIDs, search), 10, true)
						b.vc.ProcessSearchIDs = util.TrimSlice(util.SearchSlice(b.vc.ProcessIDs, search), 10, true)
					} else {
						b.vc.EntitySearchIDs = util.TrimSlice(b.vc.EntityIDs, 10, true)
						b.vc.ProcessSearchIDs = util.TrimSlice(b.vc.ProcessIDs, 10, true)
					}
					vecty.Rerender(b)
				}),
			)),
			renderEntityList(b.vc.EntitySearchIDs),
			renderProcessList(b.vc.ProcessSearchIDs, b.vc.EnvelopeHeights, b.vc.ProcessSearchList),
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain statistics..."))
}

func renderEntityList(entityIDs []string) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Heading4(vecty.Text("Entity ID list: ")),
		elem.UnorderedList(
			renderList(entityIDs)...,
		),
	)
}

func renderProcessList(processIDs []string, heights map[string]int64, procs map[string]client.ProcessInfo) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Heading4(vecty.Text("Process ID list: ")),
		elem.UnorderedList(
			renderProcessItems(processIDs, heights, procs)...,
		),
	)
}

func renderProcessItems(IDs []string, heights map[string]int64, procs map[string]client.ProcessInfo) []vecty.MarkupOrChild {
	var elemList []vecty.MarkupOrChild
	for _, ID := range IDs {
		height, hok := heights[ID]
		info, iok := procs[ID]
		elemList = append(
			elemList,
			elem.ListItem(
				vecty.Text(ID),
				vecty.If(!iok, vecty.Text(": loading process info...")),
				vecty.If(iok, vecty.Text(": type: "+info.ProcessType+", state: "+info.State)),
				vecty.If(hok, vecty.Text(", "+util.IntToString(height)+" envelopes")),
			))
	}
	return elemList
}

func renderList(slice []string) []vecty.MarkupOrChild {
	var elemList []vecty.MarkupOrChild
	for _, term := range slice {
		elemList = append(elemList, elem.ListItem(vecty.Text(term)))
	}
	return elemList
}
