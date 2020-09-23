package components

import (
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
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
			return elem.Div(renderEnvelopes(p, index)...)
		}
		return p
	}
	return elem.Div(vecty.Text("No envelopes available"))
}

func renderEnvelopes(p *Pagination, index int) []vecty.MarkupOrChild {
	empty := p.ListSize
	var elemList []vecty.MarkupOrChild
	for i := len(store.Envelopes.Envelopes) - 1; i >= len(store.Envelopes.Envelopes)-p.ListSize; i-- {
		if proto.EnvelopeIsEmpty(store.Envelopes.Envelopes[i]) {
			empty--
		} else {
			envelope := store.Envelopes.Envelopes[i]
			elemList = append(elemList, EnvelopeBlock(envelope))
		}
	}
	if empty == 0 || len(elemList) < 1 {
		if *p.Searching {
			return []vecty.MarkupOrChild{vecty.Text("No Envelopes Found With Given ID")}
		}
		return []vecty.MarkupOrChild{vecty.Text("Loading envelopes...")}
	}
	return elemList
}

// EnvelopeBlock renders a single envelope block
func EnvelopeBlock(envelope *proto.Envelope) vecty.ComponentOrHTML {
	processResults := store.Processes.ProcessResults[strings.ToLower(util.TrimHex(envelope.ProcessID))]
	processEnvelopeCount := store.Processes.EnvelopeHeights[strings.ToLower(util.TrimHex(envelope.ProcessID))]
	if processResults.EnvelopeCount < 1 && processResults.ProcessType == "" && processResults.State == "" {
		return elem.Div(
			vecty.Markup(vecty.Class("tile", "empty")),
			elem.Div(
				vecty.Markup(vecty.Class("tile-body")),
				elem.Div(
					vecty.Markup(vecty.Class("type")),
					elem.Div(
						elem.Span(
							vecty.Markup(vecty.Class("title")),
							vecty.Text("Loading envelope..."),
						),
					),
				),
			),
		)
	}
	return elem.Div(
		vecty.Markup(vecty.Class("tile", processResults.State)),
		elem.Div(
			vecty.Markup(vecty.Class("tile-body")),
			elem.Div(
				vecty.Markup(vecty.Class("type")),
				elem.Div(
					elem.Span(
						vecty.Markup(vecty.Class("title")),
						vecty.Text("#"+util.IntToString(envelope.GetGlobalHeight())),
					),
					vecty.If(
						processResults.ProcessType != "",
						elem.Span(
							vecty.Markup(vecty.Class("title")),
							vecty.Text(processResults.ProcessType),
						),
					),
					vecty.If(
						processResults.State != "",
						elem.Span(
							vecty.Markup(vecty.Class("status")),
							vecty.Text(processResults.State),
						),
					),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("contents")),
				elem.Div(
					elem.Div(
						elem.Div(
							Link(
								"/envelope/"+util.IntToString(envelope.GetGlobalHeight()),
								envelope.Nullifier,
								"hash",
							),
						),
						elem.Div(
							vecty.Markup(vecty.Class("text-truncate")),
							vecty.Text("packaged in transaction "),
							Link(
								"/transaction/"+util.IntToString(envelope.TxHeight),
								util.IntToString(envelope.TxHeight),
								"hash",
							),
						),
						elem.Div(
							vecty.Markup(vecty.Class("text-truncate")),
							vecty.If(
								processEnvelopeCount < 1,
								vecty.Text(humanize.Ordinal(int(envelope.GetProcessHeight()))+" envelope on process "),
							),
							vecty.If(
								processEnvelopeCount > 1,
								vecty.Text(humanize.Ordinal(int(envelope.GetProcessHeight()))+" of "+util.IntToString(processEnvelopeCount)+" envelopes on process "),
							),
							vecty.If(
								processEnvelopeCount == 1,
								vecty.Text("only envelope on process "),
							),
							Link(
								"/process/"+util.TrimHex(envelope.ProcessID),
								util.TrimHex(envelope.ProcessID),
								"hash",
							),
						),
					),
				),
			),
		),
	)
}
