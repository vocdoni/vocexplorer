package components

import (
	"fmt"
	"time"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/vocdoni/vocexplorer/api/dbtypes"

	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/frontend/update"
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
		elem.Section(
			vecty.Markup(vecty.Class("row", "paginated", "list")),
			elem.Div(
				vecty.Markup(vecty.Class("col-12")),
				bootstrap.Card(bootstrap.CardParams{
					Body: ValidatorDetails(),
				}),
			),
		),
	)
}

// UpdateValidatorContents keeps the validator contents page up to date
func (contents *ValidatorContents) UpdateValidatorContents() {
	dispatcher.Dispatch(&actions.SetCurrentValidator{Validator: nil})
	dispatcher.Dispatch(&actions.SetCurrentValidatorBlockCount{Count: 0})
	dispatcher.Dispatch(&actions.SetCurrentValidatorBlockList{BlockList: [config.ListSize]*dbtypes.StoreBlock{nil}})

	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: store.Client.GetGatewayInfo()})
	validator, ok := store.Client.GetValidator(store.Validators.CurrentValidatorID)
	if ok && validator != nil {
		contents.Unavailable = false
		dispatcher.Dispatch(&actions.SetCurrentValidator{Validator: validator})
	} else {
		contents.Unavailable = true
		dispatcher.Dispatch(&actions.SetCurrentValidator{Validator: nil})
		return
	}
	newVal, ok := store.Client.GetValidatorBlockHeight(util.HexToString(store.Validators.CurrentValidator.Address))
	if ok {
		dispatcher.Dispatch(&actions.SetCurrentValidatorBlockCount{Count: int(newVal)})
	}
	if !update.CheckCurrentPage("validator", ticker) {
		return
	}
	updateValidatorBlocks(contents, store.Validators.CurrentBlockCount-store.Validators.BlockPagination.Index)
	for {
		select {
		case <-store.RedirectChan:
			if !update.CheckCurrentPage("validator", ticker) {
				return
			}
		case <-ticker.C:
			if !update.CheckCurrentPage("validator", ticker) {
				return
			}
			if !store.Validators.BlockPagination.DisableUpdate {
				updateValidatorBlocks(contents, store.Validators.CurrentBlockCount-store.Validators.BlockPagination.Index)
			}
		case i := <-store.Validators.BlockPagination.PagChannel:
			if !update.CheckCurrentPage("validator", ticker) {
				return
			}
		blockloop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Validators.BlockPagination.PagChannel:
				default:
					break blockloop
				}
			}
			dispatcher.Dispatch(&actions.ValidatorBlocksIndexChange{Index: i})
			if i < 1 {
				newVal, ok := store.Client.GetValidatorBlockHeight(util.HexToString(store.Validators.CurrentValidator.Address))
				if ok {
					dispatcher.Dispatch(&actions.SetCurrentValidatorBlockCount{Count: int(newVal)})
				}
			}
			updateValidatorBlocks(contents, store.Validators.CurrentBlockCount-store.Validators.BlockPagination.Index)
		case search := <-store.Validators.BlockPagination.SearchChannel:
			if !update.CheckCurrentPage("validator", ticker) {
				return
			}
		blocksearch:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case search = <-store.Validators.BlockPagination.SearchChannel:
				default:
					break blocksearch
				}
			}
			logger.Info("search: " + search)
			dispatcher.Dispatch(&actions.ValidatorBlocksIndexChange{Index: 0})
			list, ok := store.Client.GetBlocksByValidatorSearch(search, store.Validators.CurrentValidatorID)
			if ok {
				reverseBlockList(&list)
				dispatcher.Dispatch(&actions.SetCurrentValidatorBlockList{BlockList: list})
			} else {
				dispatcher.Dispatch(&actions.SetCurrentValidatorBlockList{BlockList: [config.ListSize]*dbtypes.StoreBlock{nil}})
			}

		}
	}
}

func updateValidatorBlocks(contents *ValidatorContents, i int) {
	dispatcher.Dispatch(&actions.GatewayConnected{GatewayErr: store.Client.GetGatewayInfo()})
	newVal, ok := store.Client.GetValidatorBlockHeight(util.HexToString(store.Validators.CurrentValidator.Address))
	if ok {
		dispatcher.Dispatch(&actions.SetCurrentValidatorBlockCount{Count: int(newVal)})
	}
	if newVal > 0 {
		newList, ok := store.Client.GetBlockListByValidator(i, store.Validators.CurrentValidator.Address)
		if ok {
			reverseBlockList(&newList)
			dispatcher.Dispatch(&actions.SetCurrentValidatorBlockList{BlockList: newList})
		}
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
		elem.Div(
			vecty.Markup(vecty.Class("details")),
			elem.Span(vecty.Text(fmt.Sprintf(
				"Validated %d blocks",
				store.Validators.CurrentBlockCount,
			))),
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
			elem.DefinitionTerm(vecty.Text("Blocks")),
			elem.Description(vecty.Text(
				fmt.Sprintf("%d", store.Validators.CurrentBlockCount),
			)),
			elem.DefinitionTerm(vecty.Text("Proposing priority")),
			elem.Description(vecty.Text(
				fmt.Sprintf("%d", store.Validators.CurrentValidator.ProposerPriority),
			)),
			elem.DefinitionTerm(vecty.Text("Voting power")),
			elem.Description(vecty.Text(
				fmt.Sprintf("%d", store.Validators.CurrentValidator.VotingPower),
			)),
		),
	}
}

// ValidatorDetails renders the details of a validator contents
func ValidatorDetails() vecty.ComponentOrHTML {
	if store.Validators.CurrentBlockCount <= 0 {
		return elem.Preformatted(
			vecty.Markup(vecty.Class("empty")),
			elem.Code(
				vecty.Text("There aren't any blocks validated, yet"),
			),
		)
	}

	p := &Pagination{
		TotalPages:      int(store.Validators.CurrentBlockCount) / config.ListSize,
		TotalItems:      &store.Validators.CurrentBlockCount,
		CurrentPage:     &store.Validators.BlockPagination.CurrentPage,
		ListSize:        config.ListSize,
		DisableUpdate:   &store.Validators.BlockPagination.DisableUpdate,
		RefreshCh:       store.Validators.BlockPagination.PagChannel,
		SearchCh:        store.Validators.BlockPagination.SearchChannel,
		Searching:       &store.Validators.BlockPagination.Search,
		RenderSearchBar: true,
	}
	p.RenderFunc = func(index int) vecty.ComponentOrHTML {
		return renderValidatorBlocks(p, store.Validators.CurrentBlockList)
	}

	return elem.Div(
		elem.Heading2(
			vecty.Text("Blocks"),
		),
		p,
	)
}

func renderValidatorBlocks(p *Pagination, blocks [config.ListSize]*dbtypes.StoreBlock) vecty.ComponentOrHTML {
	var blockList []vecty.MarkupOrChild
	empty := config.ListSize
	for i := len(blocks) - 1; i >= len(blocks)-config.ListSize; i-- {
		if dbtypes.BlockIsEmpty(blocks[i]) {
			empty--
		} else {
			block := blocks[i]
			blockList = append(blockList, elem.Div(
				vecty.Markup(vecty.Class("paginated-card")),
				BlockCard(block),
			))
		}
	}
	if empty == 0 {
		if *p.Searching {
			return elem.Div(vecty.Text("No blocks found"))
		}
		return elem.Div(vecty.Text("Loading Blocks..."))
	}

	blockList = append(blockList, vecty.Markup(vecty.Class("row")))
	return elem.Div(
		blockList...,
	)
}
