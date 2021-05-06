package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
)

// ValidatorListView is the validator list component
type ValidatorListView struct {
	vecty.Core
}

// Render renders the validator list component
func (b *ValidatorListView) Render() vecty.ComponentOrHTML {
	if store.Validators.Count > 0 {
		return elem.Section(
			vecty.Markup(vecty.Class("list", "paginated")),
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading1(
						vecty.Text("Validators"),
					),
					renderValidators(),
				},
			}),
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain info..."))
}

func renderValidators() vecty.ComponentOrHTML {
	var validatorElems []vecty.MarkupOrChild

	for i := len(store.Validators.Validators) - 1; i >= 0; i-- {
		validatorElems = append(validatorElems, elem.Div(
			vecty.Markup(vecty.Class("paginated-card")),
			ValidatorCard(store.Validators.Validators[i]),
		))
	}
	if len(validatorElems) == 0 {
		return elem.Div(vecty.Text("Loading Validators..."))
	}
	validatorElems = append(validatorElems, vecty.Markup(vecty.Class("row")))
	return elem.Div(
		validatorElems...,
	)
}
