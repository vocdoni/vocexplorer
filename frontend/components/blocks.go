package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// BlocksView renders the blocks page
type BlocksView struct {
	vecty.Core
}

// Render renders the BlocksView component
func (t *BlocksView) Render() vecty.ComponentOrHTML {
	return elem.Div{
		NavBar{},
	}
}
