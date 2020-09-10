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
		router.NewRoute("/processes", &pages.ProcessesView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/process/{id}", &pages.ProcessView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/entities", &pages.EntitiesView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/entity/{id}", &pages.EntityView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/envelopes", &pages.EnvelopesView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/envelope/{id}", &pages.EnvelopeView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/blocks", &pages.BlocksView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/block/{id}", &pages.BlockView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/transactions", &pages.TxsView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/transaction/{id}", &pages.TxView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/stats", &pages.Stats{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/validators", &pages.ValidatorsView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/validator/{id}", &pages.ValidatorView{}, router.NewRouteOpts{ExactMatch: true}),
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
