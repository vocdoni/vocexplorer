package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"

	"gitlab.com/vocdoni/vocexplorer/config"
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
		p := &Pagination{
			TotalPages:      int(store.Validators.Count) / config.ListSize,
			TotalItems:      &store.Validators.Count,
			CurrentPage:     &store.Validators.Pagination.CurrentPage,
			RefreshCh:       store.Validators.Pagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Validators.Pagination.DisableUpdate,
			SearchCh:        store.Validators.Pagination.SearchChannel,
			Searching:       &store.Validators.Pagination.Search,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderValidators(p, index)
		}

		return elem.Section(
			vecty.Markup(vecty.Class("list", "paginated")),
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading1(
						vecty.Text("Validators"),
					),
					p,
				},
			}),
		)
	}
	return elem.Div(vecty.Text("Waiting for blockchain info..."))
}

func renderValidators(p *Pagination, index int) vecty.ComponentOrHTML {
	var validatorElems []vecty.MarkupOrChild

	for i := len(store.Validators.Validators) - 1; i >= 0; i-- {
		if dbtypes.ValidatorIsEmpty(store.Validators.Validators[i]) {
			continue
		}
		validatorElems = append(validatorElems, elem.Div(
			vecty.Markup(vecty.Class("paginated-card")),
			ValidatorCard(store.Validators.Validators[i]),
		))
	}
	if len(validatorElems) == 0 {
		if *p.Searching {
			return elem.Div(vecty.Text("No validators found"))
		}
		return elem.Div(vecty.Text("Loading Validators..."))
	}
	validatorElems = append(validatorElems, vecty.Markup(vecty.Class("row")))
	return elem.Div(
		validatorElems...,
	)
}
