package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/frontend/components"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	router "marwan.io/vecty-router"
)

// ValidatorsView renders the Validators page
type ValidatorsView struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the ValidatorsView component
func (home *ValidatorsView) Render() vecty.ComponentOrHTML {
	address, ok := router.GetNamedVar(home)["id"]
	// If there is an ID to look for, render individual validator page
	if ok && address != "" {
		v := new(components.ValidatorContents)
		validator, ok := dbapi.GetValidator(address)
		if validator == nil || !ok {
			log.Errorf("Validator unavailable")
			return elem.Div(
				elem.Main(vecty.Text("Validator not available")),
			)
		}
		return elem.Div(
			components.InitValidatorContentsView(v, validator, home.Cfg),
		)
	}
	t := new(rpc.TendermintInfo)
	dash := new(components.ValidatorsDashboardView)
	return components.InitValidatorsDashboardView(t, dash, home.Cfg)
}
