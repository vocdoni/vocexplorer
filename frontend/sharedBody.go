package main

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/pages"
	router "marwan.io/vecty-router"
)

// Body renders the <body> tag
type Body struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render body simply renders routes for application
func (b Body) Render() vecty.ComponentOrHTML {
	return components.SectionMain(
		router.NewRoute("/", &pages.HomeView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/vocdash", &pages.VocDashView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/processes/{id}", &pages.ProcessesView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/entities/{id}", &pages.EntitiesView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/blocktxs", &pages.BlockTxsView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/blocks/{id}", &pages.BlocksView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/txs/{id}", &pages.TxsView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/stats", &pages.Stats{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NotFoundHandler(&notFound{}),
	)
}

type notFound struct {
	vecty.Core
}

func (nf *notFound) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(prop.ID("home-view")),
		elem.Div(
			vecty.Markup(prop.ID("home-top")),
			elem.Heading1(
				vecty.Text("page not found!"),
			),
		),
	)
}
