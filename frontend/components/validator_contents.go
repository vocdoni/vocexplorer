package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/api"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ValidatorContents renders validator contents
type ValidatorContents struct {
	vecty.Core
	vecty.Mounter
	CurrentBlock int
	CurrentPage  int
	Rendered     bool
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
	return Container(
		renderServerConnectionBanner(),
		contents.renderValidatorHeader(),
		contents.renderValidatorBlockList(),
	)
}

// UpdateValidatorContents keeps the validator contents page up to date
func (contents *ValidatorContents) UpdateValidatorContents() {
	ticker := time.NewTicker(time.Duration(store.Config.RefreshTime) * time.Second)
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})
	validator, ok := api.GetValidator(store.Validators.CurrentValidatorID)
	if ok {
		dispatcher.Dispatch(&actions.SetCurrentValidator{Validator: validator})
	}
	newVal, ok := api.GetValidatorBlockHeight(util.HexToString(store.Validators.CurrentValidator.Address))
	if ok {
		dispatcher.Dispatch(&actions.SetCurrentValidatorBlockCount{Count: util.Max(int(newVal)-1, 0)})
	}
	updateValidatorBlocks(contents, store.Validators.CurrentBlockCount-contents.CurrentBlock)
	for {
		select {
		case i := <-store.Validators.Pagination.PagChannel:
		blockloop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-store.Validators.Pagination.PagChannel:
				default:
					break blockloop
				}
			}
			contents.CurrentBlock = i
			oldBlocks := store.Validators.CurrentBlockCount
			newVal, ok := api.GetValidatorBlockHeight(util.HexToString(store.Validators.CurrentValidator.Address))
			if ok {
				store.Validators.CurrentBlockCount = int(newVal) - 1
			}
			if i < 1 {
				oldBlocks = store.Validators.CurrentBlockCount
			}
			updateValidatorBlocks(contents, oldBlocks-contents.CurrentBlock)
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			if !store.Validators.Pagination.DisableUpdate {
				updateValidatorBlocks(contents, store.Validators.CurrentBlockCount-contents.CurrentBlock)
			}
		}

	}

}

func updateValidatorBlocks(contents *ValidatorContents, i int) {
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})
	newVal, ok := api.GetValidatorBlockHeight(util.HexToString(store.Validators.CurrentValidator.Address))
	if ok {
		dispatcher.Dispatch(&actions.SetCurrentValidatorBlockCount{Count: util.Max(int(newVal)-1, 0)})
	}
	if newVal > 0 {
		newList, ok := api.GetBlockListByValidator(i, store.Validators.CurrentValidator.GetAddress())
		if ok {
			reverseBlockList(&newList)
			dispatcher.Dispatch(&actions.SetCurrentValidatorBlockList{BlockList: newList})
		}
	}
}

func (contents *ValidatorContents) renderValidatorHeader() vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				elem.Heading2(
					vecty.Markup(vecty.Class("card-header")),
					vecty.Text("Validator Address "+util.HexToString(store.Validators.CurrentValidator.GetAddress())),
				),
			),
			elem.Div(
				vecty.Markup(vecty.Class("card-body")),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Blocks"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(util.IntToString(store.Validators.CurrentBlockCount)),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Priority"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(util.IntToString(store.Validators.CurrentValidator.GetProposerPriority())),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Voting Power"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(util.IntToString(store.Validators.CurrentValidator.GetVotingPower())),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("PubKey"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(util.HexToString(store.Validators.CurrentValidator.GetPubKey())),
					),
				),
			),
		),
	)
}

func (contents *ValidatorContents) renderValidatorBlockList() vecty.ComponentOrHTML {
	if store.Validators.CurrentBlockCount > 0 {
		p := &Pagination{
			TotalPages:      int(store.Validators.CurrentBlockCount) / config.ListSize,
			TotalItems:      &store.Validators.CurrentBlockCount,
			CurrentPage:     &contents.CurrentPage,
			ListSize:        config.ListSize,
			DisableUpdate:   &store.Validators.Pagination.DisableUpdate,
			RefreshCh:       store.Validators.Pagination.PagChannel,
			SearchCh:        store.Validators.Pagination.SearchChannel,
			Searching:       &store.Validators.Pagination.Search,
			RenderSearchBar: true,
		}
		p.RenderFunc = func(index int) vecty.ComponentOrHTML {
			return renderValidatorBlocks(store.Validators.CurrentBlockList)
		}

		return elem.Div(
			vecty.Markup(vecty.Class("recent-blocks")),
			elem.Heading3(
				vecty.Text("Blocks"),
			),
			p,
		)
	}
	return elem.Div(elem.Heading5(vecty.Text("No blocks validated")))
}

func renderValidatorBlocks(blocks [config.ListSize]*proto.StoreBlock) vecty.ComponentOrHTML {
	var blockList []vecty.MarkupOrChild

	empty := config.ListSize
	for i := len(blocks) - 1; i >= len(blocks)-config.ListSize; i-- {
		if proto.BlockIsEmpty(blocks[i]) {
			empty--
		} else {
			block := blocks[i]
			blockList = append(blockList, BlockCard(block))
		}
	}
	if empty == 0 {
		fmt.Println("No blocks available")
		return elem.Div(vecty.Text("Loading Blocks..."))
	}
	blockList = append(blockList, vecty.Markup(vecty.Class("responsive-card-deck")))
	return elem.Div(
		blockList...,
	)
}
