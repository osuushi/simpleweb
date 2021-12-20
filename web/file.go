package web

import (
	"strings"
	"syscall/js"
)

// Wrappers for file upload.

// Create an input element for single file upload which will push a slice of
// bytes into the provided channel.
func CreateInputElement(accept []string, ch chan<- []byte) js.Value {
	createFn := MakeFunction(`(types, cb) => {
		const input = document.createElement('input')
		input.type = "file"
		input.accept = types
		input.onchange = async function(e) {
			e.target.after(document.createElement('progress'))
			e.disabled = true
			const file = e.target.files[0]
			const arrayBuffer = await file.arrayBuffer()
			const bytes = new Uint8Array(arrayBuffer)
			cb(bytes)
		}
		return input
	}`)

	cbFn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			byteArray := args[0]
			outBytes := make([]byte, byteArray.Length())
			js.CopyBytesToGo(outBytes, byteArray)
			ch <- outBytes
		}()
		return nil
	})

	input := createFn.Invoke(strings.Join(accept, ","), cbFn)
	return input
}

// func onOpen(this js.Value, args []js.Value) interface{} {
// 	go func() {
// 		fileValue := args[0]
// 		log.Println("Received file", fileValue.Get("name"))

// 		doneCh := make(chan struct{})
// 		var data []byte
// 		onFileLoaded := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 			arr := args[0] // Uint8Array
// 			log.Println("File loaded")
// 			data = make([]byte, arr.Get("length").Int()) // Create a byte slice of the size of the array
// 			js.CopyBytesToGo(data, arr)
// 			doneCh <- struct{}{}
// 			return nil
// 		})
// 		defer onFileLoaded.Release()

// 		go func() {
// 			// Read the file
// 			web.JsExec(`(file, callback) => {
// 			file.arrayBuffer().then(buffer => {
// 				callback(new Uint8Array(buffer))
// 			})
// 		}`, fileValue, onFileLoaded)
// 		}()

// 		<-doneCh

// 		log.Println("Loaded file of size:", len(data))

// 		// Load the solid
// 		solid, err := stl.ReadAll(bytes.NewReader(data))
// 		if err != nil {
// 			log.Println("Error reading STL:", err)
// 			return
// 		}

// 		web.JsExec(`(t) => {
// 			const el = document.createElement('div')
// 			el.textContent = t
// 			document.body.append(el)
// 		}`, fmt.Sprintf("%s has %d triangles", fileValue.Get("name"), len(solid.Triangles)))
// 	}()
// 	return nil
// }
