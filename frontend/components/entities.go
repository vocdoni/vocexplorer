package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	router "marwan.io/vecty-router"
)

// EntitiesView renders the Entities page
type EntitiesView struct {
	vecty.Core
	cfg *config.Cfg
}

// Render renders the EntitiesView component
func (home *EntitiesView) Render() vecty.ComponentOrHTML {
	var entity client.EntityInfo
	var dash EntitiesDashboardView
	return elem.Div(
		&Header{},
		elem.Main(
			initEntitiesDashboardView(&entity, &dash, router.GetNamedVar(home)["id"], home.cfg),
		),
	)
}
