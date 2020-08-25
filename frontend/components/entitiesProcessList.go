package components

import (
	"strings"
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// EntityProcessListView renders the process list pane
type EntityProcessListView struct {
	vecty.Core
	entity         *client.EntityInfo
	numProcesses   int
	processesIndex int
	refreshCh      chan bool
}

// Render renders the EntityProcessListView component
func (b *EntityProcessListView) Render() vecty.ComponentOrHTML {
	if b.entity != nil {
		if js.Global().Get("searchTerm").IsUndefined() || js.Global().Get("searchTerm").String() == "" {
			b.entity.ProcessSearchIDs = util.TrimSlice(b.entity.ProcessIDs, config.ListSize, &b.processesIndex)
			b.numProcesses = len(b.entity.ProcessIDs)
		} else {
			search := js.Global().Get("searchTerm").String()
			temp := util.SearchSlice(b.entity.ProcessIDs, search)
			if len(temp) <= 0 {
				temp = entitySearchStateType(b.entity, search)
			}
			b.entity.ProcessSearchIDs = util.TrimSlice(temp, config.ListSize, &b.processesIndex)
			b.numProcesses = len(temp)
		}
		return elem.Section(
			elem.Input(vecty.Markup(
				event.Input(func(e *vecty.Event) {
					search := e.Target.Get("value").String()
					if search != "" {
						js.Global().Set("searchTerm", search)
					} else {
						js.Global().Set("searchTerm", "")
					}
					vecty.Rerender(b)
				}),
				prop.Placeholder("search IDs"),
			)),
			entityRenderProcessList(b),
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain statistics..."))
}
func entityRenderProcessList(b *EntityProcessListView) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Button(
			vecty.Text("prev"),
			vecty.Markup(
				event.Click(func(e *vecty.Event) {
					b.processesIndex--
					b.refreshCh <- true
					vecty.Rerender(b)
				}),
				vecty.MarkupIf(
					b.processesIndex > 0,
					prop.Disabled(false),
				),
				vecty.MarkupIf(
					b.processesIndex < 1,
					prop.Disabled(true),
				),
			),
		),
		elem.Button(vecty.Text("next"),
			vecty.Markup(
				event.Click(func(e *vecty.Event) {
					b.processesIndex++
					b.refreshCh <- true
					vecty.Rerender(b)
				}),
				vecty.MarkupIf(
					(b.processesIndex+1)*config.ListSize < b.numProcesses,
					prop.Disabled(false),
				),
				vecty.MarkupIf(
					(b.processesIndex+1)*config.ListSize >= b.numProcesses,
					prop.Disabled(true),
				),
			),
		),
		elem.Heading4(vecty.Text("Process ID list: ")),
		vecty.If(len(b.entity.ProcessList) < b.numProcesses, vecty.Text("Loading process info...")),
		elem.UnorderedList(
			entityRenderProcessItems(b.entity.ProcessSearchIDs, b.entity.EnvelopeHeights, b.entity.ProcessList)...,
		),
	)
}

func entityRenderProcessItems(IDs []string, heights map[string]int64, procs map[string]client.ProcessInfo) []vecty.MarkupOrChild {
	if len(IDs) == 0 {
		return []vecty.MarkupOrChild{vecty.Text("No valid processes")}
	}
	var elemList []vecty.MarkupOrChild
	for _, ID := range IDs {
		height, hok := heights[ID]
		info, iok := procs[ID]
		elemList = append(
			elemList,
			elem.ListItem(
				elem.Anchor(vecty.Markup(vecty.Attribute("href", "/processes/"+ID)), vecty.Text(ID)),
				vecty.If(!iok, vecty.Text(": loading process info...")),
				vecty.If(iok, vecty.Text(": type: "+info.ProcessType+", state: "+info.State)),
				vecty.If(hok, vecty.Text(", "+util.IntToString(height)+" envelopes")),
			))
	}
	return elemList
}

func entitySearchStateType(entity *client.EntityInfo, search string) []string {
	var IDList []string
	for _, key := range entity.ProcessIDs {
		info, ok := entity.ProcessList[key]
		if ok {
			if strings.Contains(info.State, search) || strings.Contains(info.ProcessType, search) {
				IDList = append(IDList, key)
			}
		}
	}
	return IDList
}
