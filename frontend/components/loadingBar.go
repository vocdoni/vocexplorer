package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// LoadingBar renders a basic loading bar
func LoadingBar() vecty.ComponentOrHTML {
	return elem.Progress()
}
