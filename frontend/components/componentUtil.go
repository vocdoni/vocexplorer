package components

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/util"
)

// InitGateway initializes a connection with the gateway
func InitGateway(host string) (*client.Client, context.CancelFunc) {
	// Init Gateway client
	fmt.Println("connecting to %s", host)
	gwClient, cancel, err := client.New(host)
	if util.ErrPrint(err) {
		if js.Global().Get("confirm").Invoke("Unable to connect to Gateway client. Reload with client running").Bool() {
			js.Global().Get("location").Call("reload")
		}
		return nil, cancel
	}
	return gwClient, cancel
}

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
