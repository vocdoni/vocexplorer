package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/prop"
)

// RenderList renders a set of list elements from a slice of strings
func RenderList(slice []string) []vecty.MarkupOrChild {
	var elemList []vecty.MarkupOrChild
	for _, term := range slice {
		elemList = append(elemList, elem.ListItem(vecty.Text(term)))
	}
	return elemList
}

func renderCollapsible(head, accordionName, num string, body vecty.ComponentOrHTML) vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("card", "z-depth-0", "bordered")),
		elem.Paragraph(
			elem.Button(
				vecty.Markup(
					vecty.Class("btn", "btn-link"),
					prop.Type("button"),
					vecty.Attribute("data-toggle", "collapse"),
					vecty.Attribute("data-target", "#collapse"+accordionName+num),
					vecty.Attribute("aria-expanded", "false"),
					vecty.Attribute("aria-controls", "collapse"+accordionName+num),
				),
				elem.Heading5(
					vecty.Text(head),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("collapse"),
					prop.ID("collapse"+accordionName+num),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("card", "card-body"),
					),
					body,
				),
			),
		),
	)
}
