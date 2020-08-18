package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	router "marwan.io/vecty-router"
)

// ValidatorsView renders the Validators page
type ValidatorsView struct {
	vecty.Core
	cfg *config.Cfg
}

// Render renders the ValidatorsView component
func (home *ValidatorsView) Render() vecty.ComponentOrHTML {
	v := new(ValidatorContents)
	address := router.GetNamedVar(home)["id"]
	validator := dbapi.GetValidator(address)
	if validator == nil {
		log.Errorf("Validator unavailable")
		return elem.Div(
			&Header{},
			elem.Main(vecty.Text("Validator not available")),
		)
	}
	return elem.Div(
		&Header{},
		initValidatorContentsView(v, validator, home.cfg),
	)
}
