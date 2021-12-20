package web

import "syscall/js"

// DOM convenience functions

func Global(key string) js.Value {
	return js.Global().Get(key)
}

func Document() js.Value {
	return Global("document")
}

func Body() js.Value {
	return Document().Get("body")
}

func CreateElement(tag string) js.Value {
	return Document().Call("createElement", tag)
}

func CreateSVGElement(tag string) js.Value {
	return Document().Call("createElementNS", "http://www.w3.org/2000/svg", tag)
}

func Alert(msg string) {
	Global("alert").Invoke(msg)
}

func MakeDownloadLink(data []byte, filename string, mime string) js.Value {
	// Create a blob from the data
	uint8Array := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(uint8Array, data)
	arrayBuffer := uint8Array.Get("buffer")
	blob := js.Global().Get("Blob").New([]interface{}{arrayBuffer}, map[string]interface{}{"type": mime})

	// Create the link element and set its url to the blob
	link := js.Global().Get("document").Call("createElement", "a")
	link.Set("href", js.Global().Get("URL").Call("createObjectURL", blob))
	link.Set("download", filename)
	return link
}
