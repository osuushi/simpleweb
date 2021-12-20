package web

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"syscall/js"
)

// Create labeled form elements with short and long descriptions
//
// backing must be a pointer to a value that will be set by the form
func CreateParamInput(
	ctx context.Context,
	label string,
	backing interface{},
	shortDesc string,
	longDesc string,
	changeCh chan struct{},
	attrs map[string]string,
) js.Value {
	if attrs == nil {
		attrs = make(map[string]string)
	}
	value := reflect.ValueOf(backing)
	if value.Kind() != reflect.Ptr {
		panic("backing must be a pointer")
	}
	value = value.Elem()

	// Value to set is a cleaned up value for numeric inputs, so we don't have
	// ugly precision issues in the display.
	var valueToSet string

	var inputType string
	var stepSize float64
	var min float64 = -1e9
	switch value.Kind() {
	case reflect.String:
		inputType = "text"
		valueToSet = value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		inputType = "number"
		stepSize = 1
		valueToSet = fmt.Sprintf("%d", value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		inputType = "number"
		stepSize = 1
		min = 0
		valueToSet = fmt.Sprintf("%d", value.Int())
	case reflect.Float32, reflect.Float64:
		valueToSet = strings.TrimRight(fmt.Sprintf("%.5f", value.Float()), "0.")
		inputType = "number"
		stepSize = 0.1
	default:
		panic("unsupported type for input")
	}

	// Handle overrides
	if p, ok := attrs["step"]; ok {
		stepSize, _ = strconv.ParseFloat(p, 64)
	}
	if p, ok := attrs["min"]; ok {
		min, _ = strconv.ParseFloat(p, 64)
	}

	// Create the container
	container := CreateElement("div")
	// Create the label
	labelEl := CreateElement("label")
	labelInnerEl := CreateElement("b")
	labelInnerEl.Set("textContent", label)
	labelEl.Call("appendChild", labelInnerEl)

	// Create the input and set attributes
	inputEl := CreateElement("input")
	inputEl.Set("type", inputType)
	if inputType == "number" {
		inputEl.Set("step", stepSize)
		inputEl.Set("min", min)
		if p, ok := attrs["max"]; ok {
			inputEl.Set("max", p)
		}
	}
	inputEl.Set("value", valueToSet)
	// Add the listener
	onChangeFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		target := args[0].Get("target")
		switch value.Kind() {
		case reflect.String:
			value.SetString(target.Get("value").String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value.SetInt(int64(target.Get("valueAsNumber").Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value.SetUint(uint64(target.Get("valueAsNumber").Int()))
		case reflect.Float32, reflect.Float64:
			value.SetFloat(target.Get("valueAsNumber").Float())
		default:
			panic("unsupported type for input")
		}

		changeCh <- struct{}{}
		return nil
	})

	inputEl.Call("addEventListener", "change", onChangeFunc)
	// Add the input to the label
	labelEl.Call("appendChild", inputEl)
	// Add the label to the container
	container.Call("appendChild", labelEl)

	// Create the long description disclosure
	longDescEl := CreateElement("details")
	longDescEl.Set("open", false)
	// A little custom styling to make these less emphasized in a parameter list
	longDescEl.Set("style", `
		margin-top: -45px;
		margin-left: 205px;
	`)
	londDescElSummary := CreateElement("summary")
	londDescElSummary.Set("textContent", shortDesc)
	londDescElSummary.Set("style", `
		background: white;
		padding: 5px;
	`)
	longDescEl.Call("appendChild", londDescElSummary)
	// Create the long description paragraphs by splitting on double newlines
	longDescParagraphs := strings.Split(longDesc, "\n\n")
	for _, paragraph := range longDescParagraphs {
		pEl := CreateElement("p")
		pEl.Set("innerHTML", paragraph)
		longDescEl.Call("appendChild", pEl)
	}
	container.Call("appendChild", longDescEl)

	// Clean up on cancel
	go func() {
		<-ctx.Done()
		inputEl.Call("removeEventListener", "change", onChangeFunc)
		onChangeFunc.Release()
	}()

	return container
}
