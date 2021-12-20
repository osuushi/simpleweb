package web

import "syscall/js"

// Pages are the top level of the web application. This is a very simple UI
// model which avoids getting bogged down in MVC, since everything is well
// suited to singletons.

type Page struct {
	Html    string
	Element js.Value
}

func NewPage(html string) *Page {
	page := &Page{Html: html}
	page.Init()
	return page
}

// Re-init the page's element. If you do this while the element is on the page,
// its behavior is undefined.
func (p *Page) Init() {
	element := CreateElement("div")
	element.Set("innerHTML", p.Html)
	p.Element = element
}

// Replace whatever is on the screen with the page's element
func (p *Page) Show() {
	body := Body()
	body.Set("innerHTML", "")
	body.Call("appendChild", p.Element)
}

// Replace a placeholder element (by class name) with a new element. If element
// is a string instead of an element, a span element will be created instead.
func (p *Page) Replace(className string, element interface{}) {
	if e, ok := element.(string); ok {
		elementVal := CreateElement("span")
		elementVal.Set("innerText", e)
		element = elementVal
	}

	oldElement := p.Element.Call("querySelector", "."+className)
	oldElement.Call("replaceWith", element)
}

func (p *Page) Find(selector string) js.Value {
	return p.Element.Call("querySelector", selector)
}
