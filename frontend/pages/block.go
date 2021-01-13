package pages

import (
	"strconv"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/vocdoni/vocexplorer/frontend/actions"
	"github.com/vocdoni/vocexplorer/frontend/components"
	"github.com/vocdoni/vocexplorer/frontend/dispatcher"
	"github.com/vocdoni/vocexplorer/frontend/store"
	"github.com/vocdoni/vocexplorer/logger"
	router "marwan.io/vecty-router"
)

// BlockView is the view for a single block
type BlockView struct {
	vecty.Core
}

// Render renders the Block component
func (home *BlockView) Render() vecty.ComponentOrHTML {
	dispatcher.Dispatch(&actions.SetCurrentPage{Page: "block"})

	height, err := strconv.ParseInt(router.GetNamedVar(home)["id"], 0, 64)
	if err != nil {
		logger.Error(err)
	}
	dispatcher.Dispatch(&actions.SetCurrentBlockHeight{Height: height})
	dash := new(components.BlockContents)
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateBlockContents(dash)
	return elem.Div(
		&components.Header{},
		dash,
	)
}
