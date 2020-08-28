package main

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
	"github.com/tendermint/tendermint/rpc/client/http"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/frontend/pages"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	router "marwan.io/vecty-router"
)

// Body renders the <body> tag
type Body struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render body simply renders routes for application
func (b Body) Render() vecty.ComponentOrHTML {
	store.GatewayClient, store.TendermintClient = initClients(b.Cfg)
	return components.SectionMain(
		router.NewRoute("/", &pages.HomeView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/vocdash", &pages.VocDashView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/processes/{id}", &pages.ProcessesView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/entities/{id}", &pages.EntitiesView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/envelopes/{id}", &pages.EnvelopesView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/blocktxs", &pages.BlockTxsView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/blocks", &pages.BlocksView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/blocks/{id}", &pages.BlocksView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/txs/{id}", &pages.TxsView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/stats", &pages.Stats{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/validators/{id}", &pages.ValidatorsView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		router.NewRoute("/validators", &pages.ValidatorsView{Cfg: b.Cfg}, router.NewRouteOpts{ExactMatch: true}),
		// Note that this handler only works for router.Link and router.Redirect accesses.
		// Directly accessing a non-existant route won't be handled by this.
		router.NotFoundHandler(&notFound{}),
	)
}

func initClients(cfg *config.Cfg) (*client.Client, *http.HTTP) {
	// Init tendermint client
	tClient := rpc.StartClient(cfg.TendermintHost)
	// Init Gateway client
	gwClient, _ := client.InitGateway(cfg.GatewayHost)
	if gwClient == nil || tClient == nil {
		log.Error("Cannot connect to blockchain clients")
	}
	return gwClient, tClient
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
