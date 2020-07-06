package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/nathanhack/vectyUI/nav"
)

// BlocksView renders the blocks page
type BlocksView struct {
	vecty.Core
}

// Render renders the BlocksView component
func (t *BlocksView) Render() vecty.ComponentOrHTML {
	return elem.Div{
		NavBar{}
	}
}

