package components

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/util"
	router "marwan.io/vecty-router"
)

// EntitiesView renders the Entities page
type EntitiesView struct {
	vecty.Core
	cfg *config.Cfg
}

// Render renders the EntitiesView component
func (home *EntitiesView) Render() vecty.ComponentOrHTML {
	id := router.GetNamedVar(home)["id"]
	js.Global().Set("page", "entity "+id[0:util.Min(8, len(id))]+"...")
	var entity client.EntityInfo
	var dash EntitiesDashboardView
	return elem.Div(
		&Header{},
		elem.Main(
			initEntitiesDashboardView(&entity, &dash, router.GetNamedVar(home)["id"], home.cfg),
		),
	)
}
