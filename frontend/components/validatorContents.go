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
	serverConnected     bool
	Validator           *types.Validator
	BlockList           [config.ListSize]*types.StoreBlock
	ValidatorBlocks     int
	disableBlocksUpdate bool
	CurrentBlock        int
	CurrentPage         int
	quitCh              chan struct{}
	blockRefresh        chan int
	Cfg                 *config.Cfg
}

// Render renders the ValidatorContents component
func (contents *ValidatorContents) Render() vecty.ComponentOrHTML {
	return Container(
		renderServerConnectionBanner(contents.serverConnected),
		contents.renderValidatorHeader(),
		contents.renderValidatorBlockList(),
	)
}

//InitValidatorContentsView initializes the view
func InitValidatorContentsView(v *ValidatorContents, validator *types.Validator, cfg *config.Cfg) *ValidatorContents {
	v.Validator = validator
	v.Cfg = cfg
	newVal, ok := dbapi.GetValidatorBlockHeight(util.HexToString(validator.Address))
	if ok {
		v.ValidatorBlocks = int(newVal)
	}
	v.quitCh = make(chan struct{})
	v.blockRefresh = make(chan int, 50)
	v.disableBlocksUpdate = false
	v.CurrentBlock = 0
	v.serverConnected = true
	go v.updateBlocks()
	return v
}

func (contents *ValidatorContents) updateBlocks() {
	ticker := time.NewTicker(time.Duration(contents.Cfg.RefreshTime) * time.Second)
	updateValidatorBlocks(contents, contents.ValidatorBlocks-contents.CurrentBlock)
	vecty.Rerender(contents)
	for {
		select {
		case i := <-contents.blockRefresh:
		blockloop:
			for {
				// If many indices waiting in buffer, scan to last one.
				select {
				case i = <-contents.blockRefresh:
				default:
					break blockloop
				}
			}
			contents.CurrentBlock = i
			oldBlocks := contents.ValidatorBlocks
			newVal, ok := dbapi.GetValidatorBlockHeight(util.HexToString(contents.Validator.Address))
			if ok {
				contents.ValidatorBlocks = int(newVal) - 1
			}
			if i < 1 {
				oldBlocks = contents.ValidatorBlocks
			}
			updateValidatorBlocks(contents, oldBlocks-contents.CurrentBlock)
			vecty.Rerender(contents)
		case <-ticker.C:
			if !contents.disableBlocksUpdate {
				updateValidatorBlocks(contents, contents.ValidatorBlocks-contents.CurrentBlock)
			}
			vecty.Rerender(contents)
		}

	}

}

func updateValidatorBlocks(contents *ValidatorContents, i int) {
	if !dbapi.Ping() {
		contents.serverConnected = false
	} else {
		contents.serverConnected = true
	}
	newVal, ok := dbapi.GetValidatorBlockHeight(util.HexToString(contents.Validator.Address))
	if ok {
		contents.ValidatorBlocks = int(newVal) - 1
	}
	newList, ok := dbapi.GetBlockListByValidator(i, contents.Validator.GetAddress())
	if ok {
		contents.BlockList = newList
	}
	reverseBlockList(&contents.BlockList)
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
		TotalPages:      int(contents.ValidatorBlocks) / config.ListSize,
		TotalItems:      &contents.ValidatorBlocks,
		CurrentPage:     &contents.CurrentPage,
		ListSize:        config.ListSize,
		DisableUpdate:   &contents.disableBlocksUpdate,
		RefreshCh:       contents.blockRefresh,
		RenderSearchBar: false,
	}
	p.RenderFunc = func(index int) vecty.ComponentOrHTML {
		return renderValidatorBlocks(contents.BlockList)
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
