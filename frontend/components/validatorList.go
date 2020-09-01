package components

import (
	"fmt"
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ValidatorListView is the validator list component
type ValidatorListView struct {
	vecty.Core
	currentPage int
}

// Render renders the validator list component
func (b *ValidatorListView) Render() vecty.ComponentOrHTML {
	if store.Validators.Count > 0 {
		p := &Pagination{
			TotalPages:      int(store.Validators.Count) / config.ListSize,
			TotalItems:      &store.Validators.Count,
			CurrentPage:     &b.currentPage,
			RefreshCh:       store.Validators.Pagination.PagChannel,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Validators.Pagination.DisableUpdate,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderValidators(p, index)
		}
		p.SearchBar = func(self *Pagination) vecty.ComponentOrHTML {
			return elem.Input(vecty.Markup(
				event.Input(func(e *vecty.Event) {
					search := e.Target.Get("value").String()
					index, err := strconv.Atoi(e.Target.Get("value").String())
					if err != nil || index < 0 || index > int(*self.TotalItems) || search == "" {
						*self.CurrentPage = 0
						dispatcher.Dispatch(&actions.DisableValidatorUpdate{Disabled: false})
						self.RefreshCh <- *self.CurrentPage * config.ListSize
					} else {
						*self.CurrentPage = util.Max(int(*self.TotalItems)-index-1, 0) / config.ListSize
						dispatcher.Dispatch(&actions.DisableValidatorUpdate{Disabled: true})
						self.RefreshCh <- int(*self.TotalItems) - index
					}
					vecty.Rerender(self)
				}),
				prop.Placeholder("search validators"),
			))
		}

		return elem.Section(
			vecty.Markup(vecty.Class("list", "paginated")),
			bootstrap.Card(bootstrap.CardParams{
				Body: vecty.List{
					elem.Heading3(
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

	for i := len(store.Validators.Validators) - 1; i >= len(store.Validators.Validators)-p.ListSize; i-- {
		if types.ValidatorIsEmpty(store.Validators.Validators[i]) {
			continue
		}
		validatorElems = append(validatorElems, elem.Div(
			vecty.Markup(vecty.Class("paginated-card")),
			ValidatorCard(store.Validators.Validators[i]),
		))
	}
	if len(validatorElems) == 0 {
		fmt.Println("No validators available")
		return elem.Div(vecty.Text("Loading Validators..."))
	}
	validatorElems = append(validatorElems, vecty.Markup(vecty.Class("row")))
	return elem.Div(
		validatorElems...,
	)
}
