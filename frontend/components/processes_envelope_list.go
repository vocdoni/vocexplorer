package components

import (
	"fmt"

	humanize "github.com/dustin/go-humanize"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/proto"
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
			TotalPages:      int(store.Processes.CurrentProcess.EnvelopeCount) / config.ListSize,
			TotalItems:      &store.Processes.CurrentProcess.EnvelopeCount,
			CurrentPage:     &store.Processes.EnvelopesPage,
			RefreshCh:       store.Processes.Pagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Processes.Pagination.DisableUpdate,
			SearchCh:        store.Processes.Pagination.SearchChannel,
			Searching:       &store.Processes.Pagination.Search,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderProcessEnvelopes(p, store.Processes.CurrentProcess, index)
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
		if proto.EnvelopeIsEmpty(process.Envelopes[i]) {
			empty--
		} else {
			envelope := process.Envelopes[i]
			EnvelopeList = append(EnvelopeList, renderProcessEnvelope(envelope))
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

func renderProcessEnvelope(envelope *proto.Envelope) vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-header")),
				NavLink(
					"/envelope/"+util.IntToString(envelope.GetGlobalHeight()),
					util.IntToString(envelope.GetProcessCount()),
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
