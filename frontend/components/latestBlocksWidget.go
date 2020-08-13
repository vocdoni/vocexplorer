package components

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type LatestBlocksWidget struct {
	vecty.Core
	T *rpc.TendermintInfo
}

func (b *LatestBlocksWidget) Render() vecty.ComponentOrHTML {

	var blockList []vecty.MarkupOrChild

	empty := 4
	for i := len(b.T.BlockList) - 1; i >= len(b.T.BlockList)-4; i-- {
		if types.BlockIsEmpty(b.T.BlockList[i]) {
			empty--
		}
		block := b.T.BlockList[i]
		blockList = append(blockList, BlockCard(block))
	}
	if empty == 0 {
		fmt.Println("No blocks available")
		return elem.Div(vecty.Text("Loading Blocks..."))
	}
	blockList = append(blockList, vecty.Markup(vecty.Class("responsive-card-deck")))

	return elem.Section(
		vecty.Markup(vecty.Class("recent-blocks")),
		elem.Heading4(vecty.Text("Latest blocks")),
		elem.Div(
			blockList...,
		),
	)
}

func BlockCard(block *types.StoreBlock) vecty.ComponentOrHTML {
	var tm time.Time
	var err error
	if block.GetTime() != nil {
		tm, err = ptypes.Timestamp(block.GetTime())
		util.ErrPrint(err)
	}
	p := message.NewPrinter(language.English)
	return elem.Div(
		vecty.Markup(vecty.Class("card-deck-col")),
		&bootstrap.Card{
			Header: elem.Anchor(
				vecty.Markup(
					vecty.Attribute("href", "/blocks/"+util.IntToString(block.GetHeight())),
				),
				vecty.Text(util.IntToString(block.GetHeight())),
			),
			Body: vecty.List{
				elem.Div(
					vecty.Markup(vecty.Class("block-card-heading")),
					elem.Span(
						vecty.Markup(vecty.Class("mr-2")),
						vecty.Text(p.Sprintf("%d transactions", block.GetNumTxs())),
					),
					elem.Span(
						vecty.Text(util.GetTimeAgoFormatter().Format(tm)),
					),
				),
				elem.Div(
					vecty.Markup(vecty.Class("detail", "text-truncate", "mt-2")),
					elem.Span(
						vecty.Markup(vecty.Class("dt", "mr-2")),
						vecty.Text("Hash"),
					),
					elem.Span(
						vecty.Markup(vecty.Class("dd")),
						vecty.Markup(vecty.Attribute("title", hex.EncodeToString(block.GetHash()))),
						vecty.Text(hex.EncodeToString(block.GetHash())),
					),
				),
			},
		},
	)
}
