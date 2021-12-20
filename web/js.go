package web

import (
	"fmt"
	"sync"
	"syscall/js"
)

// Cache of compiled functions
var funcRegistry *sync.Map

func init() {
	funcRegistry = &sync.Map{}
}

// Create a function from a javascript string. The function will be cached
// unless you pass `true` for nocache. Caching allows the browser to optimize
// frequently called functions, but every cached function stays in memory
// forever, so it is not appropriate for iifes.
func MakeFunction(code string, nocache ...bool) js.Value {
	useCache := len(nocache) == 0 || !nocache[0]
	if useCache {
		if existing, ok := funcRegistry.Load(code); ok {
			return existing.(js.Value)
		}
	}

	// Wrap in parens to make the function an expression
	wrapped := fmt.Sprintf("(%s)", code)
	fnValue := js.Global().Call("eval", wrapped)
	if useCache {
		funcRegistry.Store(code, fnValue)
	}

	return fnValue
}

// Call an immediate function and don't cache it
func JsExec(code string, args ...interface{}) js.Value {
	return MakeFunction(code, true).Invoke(args...)
}
