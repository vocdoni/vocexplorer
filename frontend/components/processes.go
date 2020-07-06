package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// ProcsView renders the processes page
type ProcsView struct {
	vecty.Core
}

// Render renders the ProcsView component
func (t *ProcsView) Render() vecty.ComponentOrHTML {
	return elem.Div{
		NavBar{},
	}
}
