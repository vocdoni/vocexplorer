package components

import (
	"syscall/js"

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

// BeforeUnload packages the given func in an eventlistener function to be called before page unload
func BeforeUnload(close func()) {
	var unloadFunc js.Func
	unloadFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		close()
		unloadFunc.Release() // release the function if the button will not be clicked again
		return nil
	})
	js.Global().Call("addEventListener", "beforeunload", unloadFunc)
}

// OnLoad packages the given func in an eventlistener function to be called on page load
func OnLoad(close func()) {
	var loadFunc js.Func
	loadFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		close()
		loadFunc.Release() // release the function if the button will not be clicked again
		return nil
	})
	js.Global().Call("addEventListener", "load", loadFunc)
}

// <p>
//
//   <button class="btn btn-primary" type="button" data-toggle="collapse" data-target="#collapseExample" aria-expanded="false" aria-controls="collapseExample">
//     Button with data-target
//   </button>
// </p>
// <div class="collapse" id="collapseExample">
//   <div class="card card-body">
//     Anim pariatur cliche reprehenderit, enim eiusmod high life accusamus terry richardson ad squid. Nihil anim keffiyeh helvetica, craft beer labore wes anderson cred nesciunt sapiente ea proident.
//   </div>
// </div>

func renderCollapsible(head, body, accordionName, num string) vecty.ComponentOrHTML {
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
					elem.Preformatted(vecty.Text(body)),
				),
			),
		),
	)
}
