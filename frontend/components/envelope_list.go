package components

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/proto"
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
			SearchCh:        store.Envelopes.Pagination.SearchChannel,
			Searching:       &store.Envelopes.Pagination.Search,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderEnvelopes(p, index)
		}
		return p
	}
	return elem.Div(vecty.Text("No envelopes available"))
}

func renderEnvelopes(p *Pagination, index int) vecty.ComponentOrHTML {
	var EnvelopeList []vecty.MarkupOrChild

	empty := p.ListSize
	for i := len(store.Envelopes.Envelopes) - 1; i >= len(store.Envelopes.Envelopes)-p.ListSize; i-- {
		if proto.EnvelopeIsEmpty(store.Envelopes.Envelopes[i]) {
			empty--
		} else {
			envelope := store.Envelopes.Envelopes[i]
			EnvelopeList = append(EnvelopeList, renderEnvelope(envelope))
		}
	}
	if empty == 0 {
		if *p.Searching {
			return elem.Div(vecty.Text("No Envelopes Found With Given ID"))
		}
		fmt.Println("No envelopes available")
		return elem.Div(vecty.Text("Loading envelopes..."))
	}
	EnvelopeList = append(EnvelopeList, vecty.Markup(vecty.Class("responsive-card-deck")))
	return elem.Div(
		EnvelopeList...,
	)
}

func renderEnvelope(envelope *proto.Envelope) vecty.ComponentOrHTML {
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
						vecty.Text(humanize.Ordinal(int(envelope.GetProcessCount()))+" envelope on process "+util.StripHexString(envelope.ProcessID)),
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
