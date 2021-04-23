package components

import (
	"fmt"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"go.vocdoni.io/proto/build/go/models"

	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//ValidatorCard renders a single validator card
func ValidatorCard(validator *models.Validator) vecty.ComponentOrHTML {
	vLink := func(text string) vecty.ComponentOrHTML {
		return Link(
			fmt.Sprintf("/validator/%x", validator.Address),
			text,
			"",
		)
	}
	return bootstrap.Card(bootstrap.CardParams{
		Header: vLink(validator.Name),
		Body: vecty.List{
			elem.DescriptionList(
				elem.DefinitionTerm(
					vecty.Text("Address"),
				),
				elem.Description(
					vecty.Markup(vecty.Attribute("title", util.HexToString(validator.Address))),
					vLink(util.HexToString(validator.Address)),
				),
				elem.DefinitionTerm(
					vecty.Text("PubKey"),
				),
				elem.Description(
					vecty.Text(util.HexToString(validator.PubKey)),
				),
				elem.DefinitionTerm(
					vecty.Text("Priority"),
				),
				elem.DefinitionTerm(
					vecty.Text("Voting Power"),
				),
				elem.Description(
					vecty.Text(util.IntToString(validator.Power)),
				),
			),
		},
	})
}
