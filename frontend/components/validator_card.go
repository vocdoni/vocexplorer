package components

import (
	"fmt"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//ValidatorCard renders a single validator card
func ValidatorCard(validator *dbtypes.Validator) vecty.ComponentOrHTML {
	blocks := "none"
	numBlocks, ok := store.Validators.BlockHeights[util.HexToString(validator.Address)]
	if ok || numBlocks > 0 {
		blocks = util.IntToString(numBlocks)
	}
	vLink := func(text string) vecty.ComponentOrHTML {
		return Link(
			fmt.Sprintf("/validator/%x", validator.Address),
			text,
			"",
		)
	}
	return bootstrap.Card(bootstrap.CardParams{
		Header: vLink(fmt.Sprintf("#%d", validator.Height.Height)),
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
					vecty.Text("Blocks proposed: "),
				),
				elem.Description(
					vecty.Text(blocks),
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
				elem.Description(
					vecty.Text(util.IntToString(validator.ProposerPriority)),
				),
				elem.DefinitionTerm(
					vecty.Text("Voting Power"),
				),
				elem.Description(
					vecty.Text(util.IntToString(validator.VotingPower)),
				),
			),
		},
	})
}
