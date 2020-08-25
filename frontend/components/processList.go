package components

import (
	"fmt"
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ProcessListView renders the process list pane
type ProcessListView struct {
	vecty.Core
	currentPage   int
	vochain       *client.VochainInfo
	disableUpdate *bool
	refreshCh     chan int
}

// Render renders the ProcessListView component
func (b *ProcessListView) Render() vecty.ComponentOrHTML {
	if b.vochain != nil && b.vochain.ProcessCount > 0 {
		p := &Pagination{
			TotalPages:      int(b.vochain.ProcessCount) / config.ListSize,
			TotalItems:      &b.vochain.ProcessCount,
			CurrentPage:     &b.currentPage,
			RefreshCh:       b.refreshCh,
			ListSize:        config.ListSize,
			DisableUpdate:   b.disableUpdate,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return elem.Div(renderProcessItems(b.vochain.ProcessIDs, b.vochain.EnvelopeHeights, b.vochain.ProcessResults)...)
		}
		p.SearchBar = func(self *Pagination) vecty.ComponentOrHTML {
			return elem.Input(vecty.Markup(
				event.Input(func(e *vecty.Event) {
					search := e.Target.Get("value").String()
					index, err := strconv.Atoi(e.Target.Get("value").String())
					if err != nil || index < 0 || index > int(*self.TotalItems) || search == "" {
						*self.CurrentPage = 0
						*b.disableUpdate = false
						self.RefreshCh <- *self.CurrentPage * config.ListSize
					} else {
						*self.CurrentPage = util.Max(int(*self.TotalItems)-index-1, 0) / config.ListSize
						*b.disableUpdate = true
						self.RefreshCh <- int(*self.TotalItems) - index
					}
					vecty.Rerender(self)
				}),
				prop.Placeholder("search processes"),
			))
		}
		return elem.Div(
			vecty.Markup(vecty.Class("recent-processes")),
			elem.Heading3(
				vecty.Text("Processes"),
			),
			p,
		)
	}
	if b.vochain.ProcessCount < 1 {
		return elem.Div(vecty.Text("No processes available"))
	}
	return elem.Div(vecty.Text("Waiting for processes..."))
}

func ProcessBlock(ID string, hok bool, height int64, info client.ProcessInfo) vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("tile", info.State)),
		elem.Div(
			vecty.Markup(vecty.Class("tile-body")),
			elem.Div(
				vecty.Markup(vecty.Class("type")),
				elem.Div(
					elem.Span(
						vecty.Markup(vecty.Class("title")),
						vecty.Text(info.ProcessType),
					),
					elem.Span(
						vecty.Markup(vecty.Class("status")),
						vecty.Text(info.State),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("contents")),
				elem.Div(
					elem.Div(
						elem.Anchor(
							vecty.Markup(vecty.Class("hash")),
							vecty.Markup(vecty.Attribute("href", "/processes/"+ID)),
							vecty.Text(ID),
						),
					),
					elem.Div(
						vecty.Markup(vecty.Class("envelopes")),
						vecty.Text(
							fmt.Sprintf("%d envelopes", height),
						),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("details")),
				elem.Div(
					vecty.Text("(date?)"),
				),
			),
		),
	)
}

func renderProcessItems(IDs [config.ListSize]string, heights map[string]int64, procs map[string]client.ProcessInfo) []vecty.MarkupOrChild {
	if len(IDs) == 0 {
		return []vecty.MarkupOrChild{vecty.Text("No valid processes")}
	}
	var elemList []vecty.MarkupOrChild
	for _, ID := range IDs {
		if ID != "" {
			height, hok := heights[ID]
			info, iok := procs[ID]

			if !iok {
				elemList = append(
					elemList,
					elem.Div(
						vecty.Markup(vecty.Class("loading")),
						vecty.Text("Loading process info..."),
					),
				)
				continue
			}

			elemList = append(
				elemList,
				ProcessBlock(ID, hok, height, info),
			)
		}
	}
	return elemList
}
