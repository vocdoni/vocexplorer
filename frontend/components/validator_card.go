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
	name := validator.Name
	if name == "" {
		name = util.HexToString(validator.Address)
	}
	return bootstrap.Card(bootstrap.CardParams{
		Header: vLink(name),
		Body: vecty.List{
			elem.DescriptionList(
				elem.DefinitionTerm(
					vecty.Text("Address"),
				),
				elem.Description(
					vecty.Text(util.HexToString(validator.Address)),
				),
				elem.DefinitionTerm(
					vecty.Text("PubKey"),
				),
				elem.Description(
					vecty.Text(util.HexToString(validator.PubKey)),
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
