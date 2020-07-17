package components

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

// ProcessesView renders the processes page
type ProcessesView struct {
	vecty.Core
}

// Render renders the ProcessesView component
func (home *ProcessesView) Render() vecty.ComponentOrHTML {
	id := router.GetNamedVar(home)["id"]
	js.Global().Set("page", "process "+id[0:util.Min(4, len(id))]+"...")
	var process client.FullProcessInfo
	var dash ProcessesDashboardView
	return elem.Div(
		&Header{},
		elem.Main(
			initProcessesDashboardView(&process, &dash, router.GetNamedVar(home)["id"]),
		),
	)
}
