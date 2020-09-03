package components

import (
	"encoding/hex"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/ptypes"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/types"
	"gitlab.com/vocdoni/vocexplorer/util"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//BlockCard renders a single block card
func BlockCard(block *types.StoreBlock) vecty.ComponentOrHTML {
	var tm time.Time
	var err error
	if block.GetTime() != nil {
		tm, err = ptypes.Timestamp(block.GetTime())
		util.ErrPrint(err)
	}
	p := message.NewPrinter(language.English)
	return bootstrap.Card(bootstrap.CardParams{
		Header: Link(
			"/block/"+util.IntToString(block.GetHeight()),
			"#"+util.IntToString(block.GetHeight()),
			"",
		),
		Body: vecty.List{
			elem.Div(
				vecty.Markup(vecty.Class("block-card-heading")),
				elem.Span(
					vecty.Markup(vecty.Class("mr-2")),
					vecty.Text(p.Sprintf("%d transactions", block.GetNumTxs())),
				),
				elem.Span(
					vecty.Text(humanize.Time(tm)),
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
			elem.Div(
				elem.Div(
					vecty.Markup(vecty.Class("dt")),
					vecty.Text("Proposer Address"),
				),
				elem.Div(
					Link(
						"/validator/"+util.HexToString(block.GetProposer()),
						util.HexToString(block.GetProposer()),
						"",
					),
				),
			),
		},
	})
}
