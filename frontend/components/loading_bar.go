package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

// LoadingBar renders a basic loading bar
func LoadingBar() vecty.ComponentOrHTML {
	return elem.Progress()
}
