package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"go.vocdoni.io/dvote/types"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ProcessesEnvelopeListView renders the envelope list pane
type ProcessesEnvelopeListView struct {
	vecty.Core
}

// Render renders the EnvelopeListView component
func (b *ProcessesEnvelopeListView) Render() vecty.ComponentOrHTML {
	if store.Processes.CurrentProcess.EnvelopeCount > 0 {
		p := &Pagination{
			TotalPages:      store.Processes.CurrentProcess.EnvelopeCount / config.ListSize,
			TotalItems:      &store.Processes.CurrentProcess.EnvelopeCount,
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

	for _, envelope := range store.Processes.CurrentProcess.Envelopes {
		EnvelopeList = append(EnvelopeList, renderProcessEnvelope(envelope))
	}
	if len(EnvelopeList) == 0 {
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

func renderProcessEnvelope(envelope *types.EnvelopePackage) vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-header")),
				Link(
					"/envelope/"+util.HexToString(envelope.Nullifier),
					util.HexToString(envelope.Nullifier),
					"",
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					vecty.Markup(vecty.Class("block-card-heading")),
					// elem.Div(
					// vecty.Text(humanize.Ordinal(int(envelope.GlobalHeight))+" envelope on the blockchain"),
					// ),
					elem.Div(
						elem.Div(
							vecty.Markup(vecty.Class("dt")),
							vecty.Text("Block"),
						),
						elem.Div(
							vecty.Markup(vecty.Class("dd")),
							vecty.Text(util.IntToString(envelope.Height)),
						),
					),
					elem.Div(
						elem.Div(
							vecty.Markup(vecty.Class("dt")),
							vecty.Text("Index"),
						),
						elem.Div(
							vecty.Markup(vecty.Class("dd")),
							vecty.Text(util.IntToString(envelope.TxIndex)),
						),
					),
					elem.Div(
						elem.Div(
							vecty.Markup(vecty.Class("dt")),
							vecty.Text("Transaction"),
						),
						Link(
							"/transaction/"+util.IntToString(envelope.Height)+"/"+util.IntToString(envelope.TxIndex),
							util.HexToString(envelope.TxHash),
							"hash",
						),
					),
				),
			),
		),
	)
}
