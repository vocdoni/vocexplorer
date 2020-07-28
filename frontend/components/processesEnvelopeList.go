package components

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// EnvelopeListView renders the envelope list pane
type EnvelopeListView struct {
	vecty.Core
	process        *client.FullProcessInfo
	envelopeIDs    []string
	envelopesIndex int
	numEnvelopes   int
}

// Render renders the EnvelopeListView component
func (b *EnvelopeListView) Render() vecty.ComponentOrHTML {
	if b.process != nil {
		if js.Global().Get("searchTerm").IsUndefined() || js.Global().Get("searchTerm").String() == "" {
			b.envelopeIDs = util.TrimSlice(b.process.Nullifiers, config.ListSize, &b.envelopesIndex)
			b.numEnvelopes = len(b.process.Nullifiers)
		} else {
			search := js.Global().Get("searchTerm").String()
			temp := util.SearchSlice(b.process.Nullifiers, search)
			b.envelopeIDs = util.TrimSlice(temp, config.ListSize, &b.envelopesIndex)
			b.numEnvelopes = len(temp)
		}
		return elem.Section(
			elem.Heading4(vecty.Text("Envelope ID list: ")),
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
			renderEnvelopeList(b),
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain statistics..."))
}

func renderEnvelopeList(b *EnvelopeListView) vecty.ComponentOrHTML {
	return elem.Div(
		elem.Button(
			vecty.Text("prev"),
			vecty.Markup(
				event.Click(func(e *vecty.Event) {
					b.envelopesIndex--
					vecty.Rerender(b)
				}),
				vecty.MarkupIf(
					b.envelopesIndex > 0,
					prop.Disabled(false),
				),
				vecty.MarkupIf(
					b.envelopesIndex < 1,
					prop.Disabled(true),
				),
			),
		),
		elem.Button(vecty.Text("next"),
			vecty.Markup(
				event.Click(func(e *vecty.Event) {
					b.envelopesIndex++
					vecty.Rerender(b)
				}),
				vecty.MarkupIf(
					(b.envelopesIndex+1)*config.ListSize < b.numEnvelopes,
					prop.Disabled(false),
				),
				vecty.MarkupIf(
					(b.envelopesIndex+1)*config.ListSize >= b.numEnvelopes,
					prop.Disabled(true),
				),
			),
		),
		elem.UnorderedList(
			renderEnvelopeItems(b.envelopeIDs)...,
		),
	)
}

func renderEnvelopeItems(slice []string) []vecty.MarkupOrChild {
	if len(slice) == 0 {
		return []vecty.MarkupOrChild{vecty.Text("No valid envelopes")}
	}
	var elemList []vecty.MarkupOrChild
	for _, term := range slice {
		elemList = append(elemList, elem.ListItem(vecty.Text(term)))
	}
	return elemList
}
