package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

// SectionMain defines a normal page with its main tag defined
func SectionMain(markup ...vecty.MarkupOrChild) vecty.ComponentOrHTML {
	return elem.Body(vecty.List{
		elem.Main(markup...),
		Footer(),
	})
}
