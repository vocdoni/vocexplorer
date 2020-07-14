package components

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
)

// ProcessesView renders the processes page
type ProcessesView struct {
	vecty.Core
}

// Render renders the ProcessesView component
func (home *ProcessesView) Render() vecty.ComponentOrHTML {
	js.Global().Set("page", "processes")
	js.Global().Set("apiEnabled", false)
	var vc client.VochainInfo
	var dash ProcessesDashboardView
	return elem.Div(
		&Header{},
		elem.Main(
			initProcessesDashboardView(&vc, &dash),
		),
	)
}
