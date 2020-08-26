package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"gitlab.com/vocdoni/vocexplorer/frontend/dispatcher"
)

type TabAction interface {
	alias() string
	dispatch() interface{}
	store() string
	text() string
}

type Tab struct {
	Alias string
	Text  string
}

func (b *Tab) alias() string {
	return b.Alias
}
func (b *Tab) text() string {
	return b.Text
}

func TabLink(c vecty.Component, tab TabAction) vecty.ComponentOrHTML {
	return elem.ListItem(
		elem.Button(
			vecty.Markup(
				event.Click(func(e *vecty.Event) {
					dispatcher.Dispatch(tab.dispatch())
					vecty.Rerender(c)
				}),
			),
			vecty.Markup(vecty.ClassMap{
				"active": tab.store() == tab.alias(),
			}),
			vecty.Text(tab.text()),
		),
	)
}

func TabContents(tab TabAction, contents vecty.ComponentOrHTML) vecty.MarkupOrChild {
	return vecty.If(tab.alias() == tab.store(), elem.Div(
		contents,
	))
}
