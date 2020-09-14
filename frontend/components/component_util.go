package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

// RenderList renders a set of list elements from a slice of strings
func RenderList(slice []string) []vecty.MarkupOrChild {
	var elemList []vecty.MarkupOrChild
	for _, term := range slice {
		elemList = append(elemList, elem.ListItem(vecty.Text(term)))
	}
	return elemList
}
