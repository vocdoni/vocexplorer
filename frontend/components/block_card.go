package components

import (
	"encoding/hex"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"go.vocdoni.io/proto/build/go/models"

	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//BlockCard renders a single block card
func BlockCard(block *models.BlockHeader) vecty.ComponentOrHTML {
	if block == nil {
		return nil
	}
	return bootstrap.Card(bootstrap.CardParams{
		Header: Link(
			"/block/"+util.IntToString(block.Height),
			"#"+util.IntToString(block.Height),
			"",
		),
		Body: vecty.List{
			elem.Div(
				vecty.Markup(vecty.Class("block-card-heading")),
				elem.Span(
					vecty.Markup(vecty.Class("mr-2")),
					vecty.Text(humanize.Comma(int64(block.NumTxs))+" transactions"),
				),
				elem.Span(
					vecty.Text(humanize.Time(time.Unix(block.Timestamp, 0))),
				),
			),
			elem.DescriptionList(
				elem.DefinitionTerm(
					vecty.Text("Hash"),
				),
				elem.Description(
					vecty.Markup(
						vecty.Attribute("title", hex.EncodeToString(block.BlockHash)),
						vecty.Markup(vecty.Class("text-truncate")),
					),
					vecty.Text(hex.EncodeToString(block.BlockHash)),
				),
				elem.DefinitionTerm(
					vecty.Text("Proposer Address"),
				),
				elem.Description(
					vecty.Markup(
						vecty.Class("text-truncate"),
						vecty.Attribute("title", util.HexToString(block.ProposerAddress)),
					),
					Link(
						"/validator/"+util.HexToString(block.ProposerAddress),
						util.HexToString(block.ProposerAddress),
						"",
					),
				),
			),
		},
	})
}
