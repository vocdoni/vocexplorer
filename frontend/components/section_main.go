package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// SectionMain defines a normal page with its main tag defined
func SectionMain(markup ...vecty.MarkupOrChild) vecty.ComponentOrHTML {
	return elem.Body(vecty.List{
		&Header{},
		elem.Main(markup...),
		Footer(),
	})
}
