package components

import (
	"encoding/json"

	humanize "github.com/dustin/go-humanize"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	dvotetypes "gitlab.com/vocdoni/go-dvote/types"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// EnvelopeContents renders envelope contents
type EnvelopeContents struct {
	vecty.Core
	DecryptionStatus string
	DisplayPackage   bool
	VotePackage      *dvotetypes.VotePackage
}

// Render renders the EnvelopeContents component
func (contents *EnvelopeContents) Render() vecty.ComponentOrHTML {
	return elem.Main(
		renderEnvelopeHeader(),
		vecty.Text(contents.DecryptionStatus),
		contents.renderVotePackage(),
	)
}

func renderEnvelopeHeader() vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				vecty.Markup(vecty.Class("card-header")),
				Link(
					"/envelope/"+util.IntToString(store.Envelopes.CurrentEnvelope.GetGlobalHeight()),
					util.IntToString(store.Envelopes.CurrentEnvelope.GetGlobalHeight()),
					"nav-link",
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					vecty.Markup(vecty.Class("block-card-heading")),
					elem.Div(
						vecty.Text(humanize.Ordinal(int(store.Envelopes.CurrentEnvelope.GetProcessHeight()))+" envelope on process "),
						Link(
							"/process/"+util.StripHexString(store.Envelopes.CurrentEnvelope.GetProcessID()),
							util.StripHexString(store.Envelopes.CurrentEnvelope.GetProcessID()),
							"hash",
						),
					),
					elem.Div(
						elem.Div(
							vecty.Markup(vecty.Class("dt")),
							vecty.Text("Nullifier"),
						),
						elem.Div(
							vecty.Markup(vecty.Class("dd")),
							vecty.Text(store.Envelopes.CurrentEnvelope.GetNullifier()),
						),
					),
				),
			),
		),
	)
}

func (contents *EnvelopeContents) renderVotePackage() vecty.ComponentOrHTML {
	if contents.DisplayPackage {
		voteString, err := json.MarshalIndent(contents.VotePackage, "", "\t")
		util.ErrPrint(err)
		accordionName := "accordionEnv"
		return elem.Div(
			vecty.Markup(vecty.Class("accordion"), prop.ID(accordionName)),
			renderCollapsible("Envelope Contents", accordionName, "One", elem.Preformatted(vecty.Text(string(voteString)))),
			// renderCollapsible("Data", accordionName, "Two", transactions),
			// renderCollapsible("Evidence", accordionName, "Three", elem.Preformatted(vecty.Text(string(evidence)))),
			// renderCollapsible("Last Commit", accordionName, "Four", elem.Preformatted(vecty.Text(string(commit)))),
		)
	}
	return nil
}
