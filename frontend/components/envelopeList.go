package components

import (
	"fmt"
	"strconv"

	humanize "github.com/dustin/go-humanize"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// EnvelopeListView renders the envelope list pane
type EnvelopeListView struct {
	vecty.Core
	currentPage   int
	disableUpdate *bool
	refreshCh     chan int
	vochain       *client.VochainInfo
}

// Render renders the EnvelopeListView component
func (b *EnvelopeListView) Render() vecty.ComponentOrHTML {
	if b.vochain != nil && b.vochain.EnvelopeHeight > 0 {
		p := &Pagination{
			TotalPages:      int(b.vochain.EnvelopeHeight) / config.ListSize,
			TotalItems:      &b.vochain.EnvelopeHeight,
			CurrentPage:     &b.currentPage,
			RefreshCh:       b.refreshCh,
			ListSize:        config.ListSize,
			DisableUpdate:   b.disableUpdate,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderEnvelopes(p, b.vochain, index)
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
	if b.vochain.EnvelopeHeight < 1 {
		return elem.Div(vecty.Text("No envelopes available"))
	}
	return elem.Div(vecty.Text("Waiting for envelopes..."))
}

func renderEnvelopes(p *Pagination, vochain *client.VochainInfo, index int) vecty.ComponentOrHTML {
	var EnvelopeList []vecty.MarkupOrChild

	empty := p.ListSize
	for i := len(vochain.EnvelopeList) - 1; i >= len(vochain.EnvelopeList)-p.ListSize; i-- {
		if types.EnvelopeIsEmpty(vochain.EnvelopeList[i]) {
			empty--
		} else {
			envelope := vochain.EnvelopeList[i]
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
				elem.Anchor(
					vecty.Markup(
						vecty.Class("nav-link"),
						vecty.Attribute("href", "/envelopes/"+util.IntToString(envelope.GetGlobalHeight())),
					),
					vecty.Text(util.IntToString(envelope.GetGlobalHeight())),
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
