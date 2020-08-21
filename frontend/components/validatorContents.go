package components

import (
	"fmt"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
	"gitlab.com/vocdoni/vocexplorer/dbapi"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// ValidatorContents renders validator contents
type ValidatorContents struct {
	vecty.Core
	Validator       *types.Validator
	BlockList       [config.ListSize]*types.StoreBlock
	TotalBlocks     int
	ValidatorBlocks int
	CurrentBlock    int
	blockRefresh    chan struct{}
	Cfg             *config.Cfg
}

// Render renders the ValidatorContents component
func (contents *ValidatorContents) Render() vecty.ComponentOrHTML {
	return elem.Main(
		contents.renderValidatorHeader(),
		contents.renderValidatorBlockList(),
	)
}

func InitValidatorContentsView(v *ValidatorContents, validator *types.Validator, cfg *config.Cfg) *ValidatorContents {
	v.Validator = validator
	v.Cfg = cfg
	v.TotalBlocks = int(dbapi.GetBlockHeight()) - 1
	v.ValidatorBlocks = int(dbapi.GetValidatorBlockHeight(util.HexToString(validator.Address)))
	go v.updateBlocks()
	return v
}

func (contents *ValidatorContents) updateBlocks() {
	contents.blockRefresh = make(chan struct{}, 50)
	ticker := time.NewTicker(time.Duration(contents.Cfg.RefreshTime) * time.Second)
	contents.BlockList = dbapi.GetBlockListByValidator(contents.TotalBlocks-contents.CurrentBlock, contents.Validator.GetAddress())
	reverseBlockList(&contents.BlockList)
	vecty.Rerender(contents)
	for {
		select {
		case <-contents.blockRefresh:
		blockloop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case <-contents.blockRefresh:
				default:
					break blockloop
				}
			}
			contents.BlockList = dbapi.GetBlockListByValidator(contents.TotalBlocks-contents.CurrentBlock, contents.Validator.GetAddress())
			reverseBlockList(&contents.BlockList)
			vecty.Rerender(contents)
		case <-ticker.C:
			contents.BlockList = dbapi.GetBlockListByValidator(contents.TotalBlocks-contents.CurrentBlock, contents.Validator.GetAddress())
			reverseBlockList(&contents.BlockList)
			vecty.Rerender(contents)
		}

	}

}

func (contents *ValidatorContents) renderValidatorHeader() vecty.ComponentOrHTML {
	return elem.Div(vecty.Markup(vecty.Class("card-deck-col")),
		elem.Div(vecty.Markup(vecty.Class("card")),
			elem.Div(
				elem.Heading2(
					vecty.Markup(vecty.Class("card-header")),
					vecty.Text("Validator Address "+util.HexToString(contents.Validator.GetAddress())),
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
						vecty.Text(util.IntToString(contents.ValidatorBlocks)),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Priority"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(util.IntToString(contents.Validator.GetProposerPriority())),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("Voting Power"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(util.IntToString(contents.Validator.GetVotingPower())),
					),
				),
				elem.Div(
					elem.Div(
						vecty.Markup(vecty.Class("dt")),
						vecty.Text("PubKey"),
					),
					elem.Div(
						vecty.Markup(vecty.Class("dd")),
						vecty.Text(util.HexToString(contents.Validator.GetPubKey())),
					),
				),
			),
		),
	)
}

func (contents *ValidatorContents) renderValidatorBlockList() vecty.ComponentOrHTML {
	p := &Pagination{
		TotalPages:      int(contents.TotalBlocks) / config.ListSize,
		TotalItems:      &contents.TotalBlocks,
		CurrentPage:     new(int),
		ListSize:        config.ListSize,
		RenderSearchBar: false,
	}
	p.RenderFunc = func(index int) vecty.ComponentOrHTML {
		return renderValidatorBlocks(contents.BlockList)
	}
	//TODO: keep track of pages with map
	p.PageLeft = func(e *vecty.Event) {
		*p.CurrentPage = util.Max(*p.CurrentPage-1, 0)
		contents.CurrentBlock = util.Max(contents.CurrentBlock-p.ListSize, 0)
		contents.blockRefresh <- struct{}{}
	}
	p.PageRight = func(e *vecty.Event) {
		*p.CurrentPage = util.Min(*p.CurrentPage+1, p.TotalPages)
		contents.CurrentBlock = util.Min(contents.CurrentBlock+p.ListSize, p.TotalPages)
		contents.blockRefresh <- struct{}{}
	}
	p.PageStart = func(e *vecty.Event) {
		*p.CurrentPage = 0
		contents.CurrentBlock = 0
		contents.blockRefresh <- struct{}{}
	}

	return elem.Div(
		vecty.Markup(vecty.Class("recent-blocks")),
		elem.Heading3(
			vecty.Text("Blocks"),
		),
		p,
	)
}

func renderValidatorBlocks(blocks [config.ListSize]*types.StoreBlock) vecty.ComponentOrHTML {
	var blockList []vecty.MarkupOrChild

	empty := config.ListSize
	for i := len(blocks) - 1; i >= len(blocks)-config.ListSize; i-- {
		if types.BlockIsEmpty(blocks[i]) {
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
