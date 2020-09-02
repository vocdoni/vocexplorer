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
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// EnvelopeList renders the envelope list pane
type EnvelopeList struct {
	vecty.Core
}

// Render renders the EnvelopeList component
func (b *EnvelopeList) Render() vecty.ComponentOrHTML {
	if store.Envelopes.Count > 0 {
		p := &Pagination{
			TotalPages:      int(store.Envelopes.Count) / config.ListSize,
			TotalItems:      &store.Envelopes.Count,
			CurrentPage:     &store.Envelopes.Pagination.CurrentPage,
			RefreshCh:       store.Envelopes.Pagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Envelopes.Pagination.DisableUpdate,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderEnvelopes(p, index)
		}
		p.SearchBar = func(self *Pagination) vecty.ComponentOrHTML {
			return elem.Input(vecty.Markup(
				event.Input(func(e *vecty.Event) {
					search := e.Target.Get("value").String()
					index, err := strconv.Atoi(e.Target.Get("value").String())
					if err != nil || index < 0 || index > int(*self.TotalItems) || search == "" {
						*self.CurrentPage = 0
						dispatcher.Dispatch(&actions.DisableEnvelopeUpdate{Disabled: false})
						self.RefreshCh <- *self.CurrentPage * config.ListSize
					} else {
						*self.CurrentPage = util.Max(int(*self.TotalItems)-index-1, 0) / config.ListSize
						dispatcher.Dispatch(&actions.DisableEnvelopeUpdate{Disabled: true})
						self.RefreshCh <- int(*self.TotalItems) - index
					}
					vecty.Rerender(self)
				}),
				prop.Placeholder("search envelopes"),
			))
		}
		return p
	}
	return elem.Div(vecty.Text("No envelopes available"))
}

func renderEnvelopes(p *Pagination, index int) vecty.ComponentOrHTML {
	var EnvelopeList []vecty.MarkupOrChild

	empty := p.ListSize
	for i := len(store.Envelopes.Envelopes) - 1; i >= len(store.Envelopes.Envelopes)-p.ListSize; i-- {
		if types.EnvelopeIsEmpty(store.Envelopes.Envelopes[i]) {
			empty--
		} else {
			envelope := store.Envelopes.Envelopes[i]
			EnvelopeList = append(EnvelopeList, renderEnvelope(envelope))
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

func renderEnvelope(envelope *types.Envelope) vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-header")),
				Link(
					"/envelope/"+util.IntToString(envelope.GetGlobalHeight()),
					util.IntToString(envelope.GetGlobalHeight()),
					"nav-link",
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					vecty.Markup(vecty.Class("block-card-heading")),
					elem.Div(
						vecty.Text(humanize.Ordinal(int(envelope.GetProcessHeight()))+" envelope on process "+util.StripHexString(envelope.ProcessID)),
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
