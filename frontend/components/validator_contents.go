package components

import (
	"fmt"
	"strings"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"go.vocdoni.io/proto/build/go/models"

	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/logger"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ValidatorContents renders validator contents
type ValidatorContents struct {
	vecty.Core
	vecty.Mounter
	Rendered    bool
	Unavailable bool
}

// Mount triggers when ValidatorContents renders
func (contents *ValidatorContents) Mount() {
	if !contents.Rendered {
		contents.Rendered = true
		vecty.Rerender(contents)
	}
}

// Render renders the ValidatorContents component
func (contents *ValidatorContents) Render() vecty.ComponentOrHTML {

	if !contents.Rendered {
		return LoadingBar()
	}
	if contents.Unavailable {
		return Unavailable("Validator unavailable", "")
	}
	if store.Validators.CurrentValidator == nil {
		return Unavailable("Loading validator...", "")
	}

	return Container(
		vecty.Markup(vecty.Attribute("id", "main")),
		renderServerConnectionBanner(),
		elem.Section(
			vecty.Markup(vecty.Class("details-view", "no-column")),
			elem.Div(
				vecty.Markup(vecty.Class("row")),
				elem.Div(
					vecty.Markup(vecty.Class("main-column")),
					bootstrap.Card(bootstrap.CardParams{
						Body: ValidatorView(),
					}),
				),
			),
		),
	)
}

// UpdateValidatorContents keeps the validator contents page up to date
func (contents *ValidatorContents) UpdateValidatorContents() {
	dispatcher.Dispatch(&actions.SetCurrentValidator{Validator: nil})

	dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: store.Client.GetGatewayInfo()})
	validators, err := store.Client.GetValidatorList()
	var currentValidator *models.Validator
	if err != nil {
		logger.Error(err)
	} else {
		dispatcher.Dispatch(&actions.SetValidatorList{List: validators})
		for _, validator := range validators {
			if strings.Contains(util.HexToString(validator.Address), store.Validators.CurrentValidatorID) {
				currentValidator = validator
			}
		}
	}
	if currentValidator != nil {
		contents.Unavailable = false
		dispatcher.Dispatch(&actions.SetCurrentValidator{Validator: currentValidator})
	} else {
		contents.Unavailable = true
		dispatcher.Dispatch(&actions.SetCurrentValidator{Validator: nil})
		return
	}
}

// ValidatorView renders a single validator
func ValidatorView() vecty.List {
	return vecty.List{
		elem.Heading1(
			vecty.Markup(vecty.Class("card-title")),
			vecty.Text("Validator details"),
		),
		elem.Heading2(
			vecty.Text(fmt.Sprintf(
				"Validator address: %x",
				store.Validators.CurrentValidator.Address,
			)),
		),
		elem.HorizontalRule(),
		elem.DescriptionList(
			elem.DefinitionTerm(vecty.Text("Address")),
			elem.Description(vecty.Text(
				fmt.Sprintf("%x", store.Validators.CurrentValidator.Address),
			)),
			elem.DefinitionTerm(vecty.Text("Public key")),
			elem.Description(vecty.Text(
				fmt.Sprintf("%x", store.Validators.CurrentValidator.PubKey),
			)),
			elem.DefinitionTerm(vecty.Text("Voting power")),
			elem.Description(vecty.Text(
				fmt.Sprintf("%d", store.Validators.CurrentValidator.Power),
			)),
		),
	}
}
