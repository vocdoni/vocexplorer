package components

import (
	"encoding/hex"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/proto"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//BlockCard renders a single block card
func BlockCard(block *proto.StoreBlock) vecty.ComponentOrHTML {
	var tm time.Time
	if block.GetTime() != nil {
		tm = time.Unix(block.GetTime().Seconds, int64(block.GetTime().Nanos)).UTC()
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
