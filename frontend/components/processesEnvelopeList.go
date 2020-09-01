package components

import (
	"fmt"
	"strconv"

	humanize "github.com/dustin/go-humanize"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ProcessesEnvelopeListView renders the envelope list pane
type ProcessesEnvelopeListView struct {
	vecty.Core
	currentPage int
}

// Render renders the EnvelopeListView component
func (b *ProcessesEnvelopeListView) Render() vecty.ComponentOrHTML {
	if store.Processes.CurrentProcess.EnvelopeCount > 0 {
		p := &Pagination{
			TotalPages:      int(store.Processes.CurrentProcess.EnvelopeCount) / config.ListSize,
			TotalItems:      &store.Processes.CurrentProcess.EnvelopeCount,
			CurrentPage:     &b.currentPage,
			RefreshCh:       store.Processes.Pagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Processes.Pagination.DisableUpdate,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderProcessEnvelopes(p, store.Processes.CurrentProcess, index)
		}
		p.SearchBar = func(self *Pagination) vecty.ComponentOrHTML {
			return elem.Input(vecty.Markup(
				event.Input(func(e *vecty.Event) {
					search := e.Target.Get("value").String()
					index, err := strconv.Atoi(e.Target.Get("value").String())
					if err != nil || index < 0 || index > int(*self.TotalItems) || search == "" {
						*self.CurrentPage = 0
						dispatcher.Dispatch(&actions.DisableProcessUpdate{Disabled: false})
						self.RefreshCh <- *self.CurrentPage * config.ListSize
					} else {
						*self.CurrentPage = util.Max(int(*self.TotalItems)-index-1, 0) / config.ListSize
						dispatcher.Dispatch(&actions.DisableProcessUpdate{Disabled: true})
						self.RefreshCh <- int(*self.TotalItems) - index
					}
					vecty.Rerender(self)
				}),
				prop.Placeholder("search envelopes"),
			))
		}
		return elem.Div(
			vecty.Markup(vecty.Class("recent-envelopes")),
			elem.Heading3(
				vecty.Text("Envelopes"),
			),
			p,
		)
	}
	return elem.Preformatted(
		vecty.Markup(vecty.Class("empty")),
		vecty.Text("This process has no envelopes"),
	)
}

func renderProcessEnvelopes(p *Pagination, process storeutil.Process, index int) vecty.ComponentOrHTML {
	var EnvelopeList []vecty.MarkupOrChild

	empty := p.ListSize
	for i := len(process.Envelopes) - 1; i >= len(process.Envelopes)-p.ListSize; i-- {
		if types.EnvelopeIsEmpty(process.Envelopes[i]) {
			empty--
		} else {
			envelope := process.Envelopes[i]
			EnvelopeList = append(EnvelopeList, renderProcessEnvelope(envelope))
		}
	}
	if empty == 0 {
		fmt.Println("No envelopes available")
		return elem.Div(vecty.Text("Loading envelopes..."))
	}
	EnvelopeList = append(EnvelopeList, vecty.Markup(vecty.Class("responsive-card-deck")))
	return elem.Div(
		EnvelopeList...,
	)
}

func renderProcessEnvelope(envelope *types.Envelope) vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-header")),
				Link(
					"/envelopes/"+util.IntToString(envelope.GetGlobalHeight()),
					util.IntToString(envelope.GetProcessHeight()),
					"nav-link",
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					vecty.Markup(vecty.Class("block-card-heading")),
					elem.Div(
						vecty.Text(humanize.Ordinal(int(envelope.GetGlobalHeight()))+" envelope on the blockchain"),
					),
					elem.Div(
						elem.Div(
							vecty.Markup(vecty.Class("dt")),
							vecty.Text("Nullifier"),
						),
						elem.Div(
							vecty.Markup(vecty.Class("dd")),
							vecty.Text(envelope.GetNullifier()),
						),
					),
				),
			),
		),
	)
}

// func renderEnvelopeList(b *EnvelopeListView) vecty.ComponentOrHTML {
// 	return elem.Div(
// 		elem.Button(
// 			vecty.Text("prev"),
// 			vecty.Markup(
// 				event.Click(func(e *vecty.Event) {
// 					b.envelopesIndex--
// 					vecty.Rerender(b)
// 				}),
// 				vecty.MarkupIf(
// 					b.envelopesIndex > 0,
// 					prop.Disabled(false),
// 				),
// 				vecty.MarkupIf(
// 					b.envelopesIndex < 1,
// 					prop.Disabled(true),
// 				),
// 			),
// 		),
// 		elem.Button(vecty.Text("next"),
// 			vecty.Markup(
// 				event.Click(func(e *vecty.Event) {
// 					b.envelopesIndex++
// 					vecty.Rerender(b)
// 				}),
// 				vecty.MarkupIf(
// 					(b.envelopesIndex+1)*config.ListSize < b.numEnvelopes,
// 					prop.Disabled(false),
// 				),
// 				vecty.MarkupIf(
// 					(b.envelopesIndex+1)*config.ListSize >= b.numEnvelopes,
// 					prop.Disabled(true),
// 				),
// 			),
// 		),
// 		elem.UnorderedList(
// 			renderEnvelopeItems(b.envelopeIDs)...,
// 		),
// 	)
// }

// func renderEnvelopeItems(slice []string) []vecty.MarkupOrChild {
// 	if len(slice) == 0 {
// 		return []vecty.MarkupOrChild{vecty.Text("No valid envelopes")}
// 	}
// 	var elemList []vecty.MarkupOrChild
// 	for _, term := range slice {
// 		elemList = append(elemList, elem.ListItem(vecty.Text(term)))
// 	}
// 	return elemList
// }
