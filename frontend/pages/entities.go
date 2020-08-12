package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	router "marwan.io/vecty-router"
)

// EntitiesView renders the Entities page
type EntitiesView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the EntitiesView component
func (home *EntitiesView) Render() vecty.ComponentOrHTML {
	var entity client.EntityInfo
	var dash components.EntitiesDashboardView
	return elem.Div(
		&components.Header{},
		elem.Main(
			components.InitEntitiesDashboardView(&entity, &dash, router.GetNamedVar(home)["id"], home.Cfg),
		),
	)
}
