package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//ValidatorCard renders a single validator card
func ValidatorCard(validator *types.Validator) vecty.ComponentOrHTML {
	return bootstrap.Card(bootstrap.CardParams{
		Header: elem.Anchor(
			vecty.Markup(
				vecty.Attribute("href", "/validators/"+util.HexToString(validator.GetAddress())),
			),
			vecty.Text("#"+util.IntToString(validator.GetHeight().GetHeight())),
		),
		Body: vecty.List{
			elem.Div(
				vecty.Markup(vecty.Class("detail", "text-truncate", "mt-2")),
				elem.Span(
					vecty.Markup(vecty.Class("dt", "mr-2")),
					vecty.Text("Address"),
				),
				elem.Span(
					vecty.Markup(vecty.Class("dd")),
					vecty.Markup(vecty.Attribute("title", util.HexToString(validator.GetAddress()))),
					vecty.Text(util.HexToString(validator.GetAddress())),
				),
			),
			elem.Div(
				elem.Div(
					vecty.Markup(vecty.Class("dt")),
					vecty.Text("PubKey"),
				),
				elem.Span(
					vecty.Markup(vecty.Class("dd")),
					vecty.Text(util.HexToString(validator.GetPubKey())),
				),
			),
			elem.Div(
				elem.Div(
					vecty.Markup(vecty.Class("dt")),
					vecty.Text("Priority"),
				),
				elem.Span(
					vecty.Markup(vecty.Class("dd")),
					vecty.Text(util.IntToString(validator.GetProposerPriority())),
				),
			),
			elem.Div(
				elem.Div(
					vecty.Markup(vecty.Class("dt")),
					vecty.Text("Voting Power"),
				),
				elem.Span(
					vecty.Markup(vecty.Class("dd")),
					vecty.Text(util.IntToString(validator.GetVotingPower())),
				),
			),
		},
	})
}