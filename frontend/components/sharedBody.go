package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	router "marwan.io/vecty-router"
)

// Body renders the <body> tag
type Body struct {
	vecty.Core
}

// Render body simply renders routes for application
func (b Body) Render() vecty.ComponentOrHTML {
	return elem.Body(
		router.NewRoute("/", &HomeView{}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/processes", &ProcessesView{}, router.NewRouteOpts{ExactMatch: true}),
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
