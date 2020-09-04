package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

func LoadingBar() vecty.ComponentOrHTML {
	return elem.Progress()
}
