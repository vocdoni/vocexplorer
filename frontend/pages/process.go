package pages

import (
	"encoding/hex"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"go.vocdoni.io/proto/build/go/models"

	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/store/storeutil"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

// ProcessView renders the Process page
type ProcessView struct {
	vecty.Core
}

// Render renders the ProcessView component
func (home *ProcessView) Render() vecty.ComponentOrHTML {
	dispatcher.Dispatch(&actions.SetCurrentPage{Page: "process"})
	dash := new(components.ProcessContentsView)
	pid, err := hex.DecodeString(util.TrimHex(router.GetNamedVar(home)["id"]))
	if err != nil {
		logger.Error(err)
	}
	dispatcher.Dispatch(&actions.SetCurrentProcessStruct{Process: &storeutil.Process{Process: &models.Process{ProcessId: pid}}})
	dash.Rendered = false
	// Ensure component rerender is only triggered once component has been rendered
	if !store.Listeners.Has(dash) {
		store.Listeners.Add(dash, func() {
			if dash.Rendered {
				vecty.Rerender(dash)
			}
		})
	}
	go components.UpdateProcessContents(dash)
	return elem.Div(
		&components.Header{},
		dash,
	)
}
