package components

import (
	humanize "github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/vocdoni/vocexplorer/api/dbtypes"
	"github.com/vocdoni/vocexplorer/config"
	"github.com/vocdoni/vocexplorer/frontend/store"
	"github.com/vocdoni/vocexplorer/util"
)

// ProcessesEnvelopeListView renders the envelope list pane
type ProcessesEnvelopeListView struct {
	vecty.Core
}

// Render renders the EnvelopeListView component
func (b *ProcessesEnvelopeListView) Render() vecty.ComponentOrHTML {
	if store.Processes.CurrentProcessResults.EnvelopeCount > 0 {
		p := &Pagination{
			TotalPages:      int(store.Processes.CurrentProcessResults.EnvelopeCount) / config.ListSize,
			TotalItems:      &store.Processes.CurrentProcessResults.EnvelopeCount,
			CurrentPage:     &store.Processes.EnvelopePagination.CurrentPage,
			RefreshCh:       store.Processes.EnvelopePagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Processes.EnvelopePagination.DisableUpdate,
			SearchCh:        store.Processes.EnvelopePagination.SearchChannel,
			Searching:       &store.Processes.EnvelopePagination.Search,
			RenderSearchBar: false,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderProcessEnvelopes(p, index)
		}
		return elem.Div(
			p,
		)
	}
	return elem.Preformatted(
		vecty.Markup(vecty.Class("empty")),
		vecty.Text("This process has no envelopes"),
	)
}

func renderProcessEnvelopes(p *Pagination, index int) vecty.ComponentOrHTML {

	var EnvelopeList []vecty.MarkupOrChild

	empty := p.ListSize
	for i := len(store.Processes.CurrentProcessEnvelopes) - 1; i >= len(store.Processes.CurrentProcessEnvelopes)-p.ListSize; i-- {
		if dbtypes.EnvelopeIsEmpty(store.Processes.CurrentProcessEnvelopes[i]) {
			empty--
		} else {
			envelope := store.Processes.CurrentProcessEnvelopes[i]
			EnvelopeList = append(EnvelopeList, renderProcessEnvelope(envelope))
		}
	}
	if empty == 0 {
		if *p.Searching {
			return elem.Div(vecty.Text("No envelopes found"))
		}
		return elem.Div(vecty.Text("No envelopes available"))
	}
	EnvelopeList = append(EnvelopeList, vecty.Markup(vecty.Class("responsive-card-deck")))
	return elem.Div(
		EnvelopeList...,
	)
}

func renderProcessEnvelope(envelope *dbtypes.Envelope) vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-header")),
				Link(
					"/envelope/"+util.IntToString(envelope.GlobalHeight),
					util.IntToString(envelope.ProcessHeight),
					"",
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					vecty.Markup(vecty.Class("block-card-heading")),
					elem.Div(
						vecty.Text(humanize.Ordinal(int(envelope.GlobalHeight))+" envelope on the blockchain"),
					),
					elem.Div(
						elem.Div(
							vecty.Markup(vecty.Class("dt")),
							vecty.Text("Nullifier"),
						),
						elem.Div(
							vecty.Markup(vecty.Class("dd")),
							vecty.Text(envelope.Nullifier),
						),
					),
					elem.Div(
						elem.Div(
							vecty.Markup(vecty.Class("dt")),
							vecty.Text("Transaction"),
						),
						Link(
							"/transaction/"+util.IntToString(envelope.TxHeight),
							util.IntToString(envelope.TxHeight),
							"hash",
						),
					),
				),
			),
		),
	)
}
