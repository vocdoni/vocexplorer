package components

import (
	"fmt"
	"log"
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
	Rendered bool
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
		dispatcher.Dispatch(&actions.SetCurrentValidatorBlockCount{Count: int(newVal)})
	}
	updateValidatorBlocks(contents, store.Validators.CurrentBlockCount-store.Validators.BlockPagination.Index)
	for {
		select {
		case i := <-store.Validators.BlockPagination.PagChannel:
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
			oldBlocks := store.Validators.CurrentBlockCount
			newVal, ok := api.GetValidatorBlockHeight(util.HexToString(store.Validators.CurrentValidator.Address))
			if ok {
				dispatcher.Dispatch(&actions.SetCurrentValidatorBlockCount{Count: int(newVal)})
			}
			if i < 1 {
				oldBlocks = store.Validators.CurrentBlockCount
			}
			updateValidatorBlocks(contents, oldBlocks-store.Validators.BlockPagination.Index)
		case search := <-store.Validators.BlockPagination.SearchChannel:
		blocksearch:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case search = <-store.Validators.BlockPagination.SearchChannel:
				default:
					break blocksearch
				}
			}
			log.Println("search: " + search)
			dispatcher.Dispatch(&actions.ValidatorBlocksIndexChange{Index: 0})
			list, ok := api.GetBlocksByValidatorSearch(search, store.Validators.CurrentValidatorID)
			if ok {
				reverseBlockList(&list)
				dispatcher.Dispatch(&actions.SetCurrentValidatorBlockList{BlockList: list})
			} else {
				dispatcher.Dispatch(&actions.SetCurrentValidatorBlockList{BlockList: [config.ListSize]*proto.StoreBlock{nil}})
			}
		case <-store.RedirectChan:
			fmt.Println("Redirecting...")
			ticker.Stop()
			return
		case <-ticker.C:
			if !store.Validators.BlockPagination.DisableUpdate {
				updateValidatorBlocks(contents, store.Validators.CurrentBlockCount-store.Validators.BlockPagination.Index)
			}
		}

	}

}

func updateValidatorBlocks(contents *ValidatorContents, i int) {
	dispatcher.Dispatch(&actions.ServerConnected{Connected: api.PingServer()})
	newVal, ok := api.GetValidatorBlockHeight(util.HexToString(store.Validators.CurrentValidator.Address))
	if ok {
		dispatcher.Dispatch(&actions.SetCurrentValidatorBlockCount{Count: int(newVal)})
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
			vecty.Markup(vecty.Class("recent-blocks")),
			elem.Heading3(
				vecty.Text("Blocks"),
			),
			p,
		)
	}
	return elem.Div(elem.Heading5(vecty.Text("No blocks validated")))
}

func renderValidatorBlocks(p *Pagination, blocks [config.ListSize]*proto.StoreBlock) vecty.ComponentOrHTML {
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
		if *p.Searching {
			return elem.Div(vecty.Text("No Blocks Found With Given ID"))
		}
		return elem.Div(vecty.Text("Loading Blocks..."))
	}
	blockList = append(blockList, vecty.Markup(vecty.Class("responsive-card-deck")))
	return elem.Div(
		blockList...,
	)
}
