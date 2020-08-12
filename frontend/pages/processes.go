package pages

import (
	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	router "marwan.io/vecty-router"
)

// ProcessesView renders the processes page
type ProcessesView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the ProcessesView component
func (home *ProcessesView) Render() vecty.ComponentOrHTML {
	var process client.FullProcessInfo
	var dash components.ProcessesDashboardView
	return components.InitProcessesDashboardView(&process, &dash, router.GetNamedVar(home)["id"], home.Cfg)
}
