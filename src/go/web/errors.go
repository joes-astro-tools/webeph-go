//go:build js && wasm

package web

import "syscall/js"

var ErrMsg string

//export anyErrors
func AnyErrors() int {
	return len(ErrMsg)
}

//export printErr
func PrintErr() {
	console := js.Global().Get("console")
	console.Call("error", ErrMsg)
}
