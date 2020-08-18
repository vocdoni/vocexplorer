package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ValidatorContents renders validator contents
type ValidatorContents struct {
	vecty.Core
	Validator *types.Validator
}

// Render renders the ValidatorContents component
func (contents *ValidatorContents) Render() vecty.ComponentOrHTML {
	return elem.Main(
		renderValidatorHeader(contents.Validator),
	)
}

func renderValidatorHeader(val *types.Validator) vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				elem.Heading2(
					vecty.Markup(vecty.Class("card-header")),
					vecty.Text("Address "+util.HexToString(val.GetAddress())),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Priority"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(util.IntToString(val.GetProposerPriority())),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Voting Power"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(util.IntToString(val.GetVotingPower())),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("PubKey"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(util.HexToString(val.GetPubKey())),
					),
				),
			),
		),
	)
}
