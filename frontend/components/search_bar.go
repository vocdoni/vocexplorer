package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/event"
	"gitlab.com/vocdoni/vocexplorer/frontend/actions"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
	router "marwan.io/vecty-router"
)

//SearchBar is a component for a search bar
type SearchBar struct {
	vecty.Core
}

//Render renders the SearchBar component
func (s *SearchBar) Render() vecty.ComponentOrHTML {
	return elem.Div(
		elem.Input(
			vecty.Markup(vecty.Class("form-control", "mr-sm-4")),
			vecty.Markup(vecty.Attribute("aria-label", "Search by vote, process, entity id")),
			vecty.Markup(vecty.Attribute("placeholder", "Search by vote, process, entity id")),
			vecty.Markup(vecty.Attribute("type", "search")),
			// Trigger when 'enter' is pressed
			vecty.Markup(event.Change(func(e *vecty.Event) {
				search := e.Target.Get("value").String()
				if len(search) == 0 {
					return
				}
				if len(search) > 1 && (search[:2] == "0x" || search[:2] == "0X") {
					search = search[2:]
				}
				dispatcher.Dispatch(&actions.SetCurrentPage{Page: ""})
				dispatcher.Dispatch(&actions.SignalRedirect{})
				router.Redirect("/search/" + search)
			}),
			),
		),
	)
}
