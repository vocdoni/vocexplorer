package main

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/pages"
	router "marwan.io/vecty-router"
)

// Body renders the <body> tag
type Body struct {
	vecty.Core
}

// Render body simply renders routes for application
func (b Body) Render() vecty.ComponentOrHTML {
	return components.SectionMain(
		router.NewRoute("/", &pages.HomeView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/vocdash", &pages.VocDashView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/process/{id}", &pages.ProcessesView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/entity/{id}", &pages.EntitiesView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/envelope/{id}", &pages.EnvelopesView{Rendered: false}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/blocktxs", &pages.BlockTxsView{}, router.NewRouteOpts{ExactMatch: true}),
		// router.NewRoute("/blocks", &pages.BlocksView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/block/{id}", &pages.BlocksView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/tx/{id}", &pages.TxsView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/stats", &pages.Stats{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/validator/{id}", &pages.ValidatorsView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/validators", &pages.ValidatorsView{}, router.NewRouteOpts{ExactMatch: true}),
		// Note that this handler only works for router.Link and router.Redirect accesses.
		// Directly accessing a non-existant route won't be handled by this.
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
