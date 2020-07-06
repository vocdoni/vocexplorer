package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// DashboardView renders the dashboard landing page
type DashboardView struct {
	vecty.Core
}

// Render renders the DashboardView component
func (dash *DashboardView) Render() vecty.ComponentOrHTML {
	return elem.Div(
		&Header{currentPage: "dashboard"},
	)
}
