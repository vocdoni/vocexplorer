package components

import (
	"encoding/hex"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

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
		bootstrap.Card(
			elem.Anchor(
				vecty.Markup(
					vecty.Attribute("href", "/blocks/"+util.IntToString(block.GetHeight())),
				),
				vecty.Text(util.IntToString(block.GetHeight())),
			),
			vecty.List{
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
			nil,
		),
	)
}
