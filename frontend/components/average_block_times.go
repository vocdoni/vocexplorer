package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/bootstrap"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//AverageBlockTimes is the component to display avg block times
type AverageBlockTimes struct {
	vecty.Core
}

//Render renders the AverageBlockTimes component
func (a *AverageBlockTimes) Render() vecty.ComponentOrHTML {

	if store.Stats.BlockTime == nil {
		return &bootstrap.Alert{
			Type:     "warning",
			Contents: "Waiting for block times",
		}
	}

	var items vecty.List

	names := map[int]string{
		0: "1m",
		1: "10m",
		2: "1h",
		3: "6h",
		4: "24h",
	}

	for k, bt := range store.Stats.BlockTime {
		if bt <= 0 {
			continue
		}

		items = append(items, elem.TableRow(
			elem.TableData(vecty.Text(names[k])),
			elem.TableData(vecty.Text(util.MsToString(bt))),
		))
	}

	return elem.Section(
		bootstrap.Card(bootstrap.CardParams{
			Body: vecty.List{
				elem.Heading2(vecty.Text("Average block times: ")),
				elem.Table(
					vecty.Markup(
						vecty.Class("table"),
						vecty.Attribute("aria-label", "Table of average block times for given time periods."),
					),
					elem.TableHead(
						elem.TableRow(
							elem.TableHeader(vecty.Text("Time period")),
							elem.TableHeader(vecty.Text("Avg. time"))),
					),
					elem.TableBody(
						items,
					),
				),
			},
		}),
	)
}
