package pages

import (
	"github.com/gopherjs/vecty"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
)

// HomeView renders the Home landing page
type HomeView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the HomeView component
func (home *HomeView) Render() vecty.ComponentOrHTML {
	dash := new(components.DashboardView)
	return dash
}
