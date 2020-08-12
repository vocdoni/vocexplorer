package pages

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/config"
)

// Stats is a pretty page for all our blockchain statistics
type Stats struct {
	vecty.Core
	Cfg *config.Cfg
}

// Render renders the Stats component
func (stats *Stats) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Text("hi :)"),
	)
}
