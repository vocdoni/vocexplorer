package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	router "marwan.io/vecty-router"
)

// ProcessesView renders the processes page
type ProcessesView struct {
	vecty.Core
	cfg *config.Cfg
}

// Render renders the ProcessesView component
func (home *ProcessesView) Render() vecty.ComponentOrHTML {
	process := new(client.FullProcessInfo)
	dash := new(ProcessesDashboardView)
	return elem.Div(
		&Header{},
		elem.Main(
			initProcessesDashboardView(process, dash, router.GetNamedVar(home)["id"], home.cfg),
		),
	)
}
