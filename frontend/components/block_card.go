package components

import (
	"encoding/hex"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/golang/protobuf/ptypes"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/go-dvote/log"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//BlockCard renders a single block card
func BlockCard(block *proto.StoreBlock) vecty.ComponentOrHTML {
	var tm time.Time
	var err error
	if block.GetTime() != nil {
		tm, err = ptypes.Timestamp(block.GetTime())
		if err != nil {
			log.Error(err)
		}
	}
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
					vecty.Text(humanize.Comma(block.GetNumTxs())+" transactions"),
				),
				elem.Span(
					vecty.Text(humanize.Time(tm)),
				),
			),
			elem.DescriptionList(
				elem.DefinitionTerm(
					vecty.Text("Hash"),
				),
				elem.Description(
					vecty.Markup(
						vecty.Attribute("title", hex.EncodeToString(block.GetHash())),
						vecty.Markup(vecty.Class("text-truncate")),
					),
					vecty.Text(hex.EncodeToString(block.GetHash())),
				),
				elem.DefinitionTerm(
					vecty.Text("Proposer Address"),
				),
				elem.Description(
					vecty.Markup(
						vecty.Class("text-truncate"),
						vecty.Attribute("title", util.HexToString(block.GetProposer())),
					),
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
