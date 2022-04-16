package tagger

import "reflect"

type TagSymbols struct {
	// key:value (: - symbol)
	KeyValue string
	// key|key2 (| - symbol)
	KeysSeparator string
}

// InHandlerF function handler for a tag
type InHandlerF func(data any, field *Field, in *reflect.Value) error

// InHandlerC struct which implementing contract for a tag
type InHandlerC interface {
	Handle(data any, field *Field, in *reflect.Value) error
}

// OutHandlerF function handler for a tag
type OutHandlerF func(data any, field *Field, in *reflect.Value) (any, error)

// OutHandlerC struct which implementing contract for a tag
type OutHandlerC interface {
	Handle(data any, field *Field, in *reflect.Value) (any, error)
}

// Tag custom tag with described name, symbols and handlers
// [In] & [Out] cannot be described twice, it's a useless
type Tag struct {
	Name       string
	TagSymbols TagSymbols

	InHandlerF InHandlerF
	InHandlerC InHandlerC

	OutHandlerF OutHandlerF
	OutHandlerC OutHandlerC
}

// InFunction pass handler which implementing the contract for Out operation
//
//    func(field *Field, in *reflect.Value) error {
//      return field.Set("foo")
//    }
//    New("json").InHandler(handler
func (t *Tag) InFunction(handler InHandlerF) *Tag {
	t.InHandlerF = handler

	return t
}

// InContract pass handler which implementing the contract for In operation
//
//   type JsonHandler struct {}
//   func (s *JsonHandler) Handle(field *Field, in *reflect.Value) error {
//     return field.Set("myValue")
//   }
//
//   New("json").InHandlerContract(JsonHandler{})
//
func (t *Tag) InContract(handler InHandlerC) *Tag {
	t.InHandlerC = handler

	return t
}

// OutFunction pass handler which implementing the contract for Out operation
//
//    func handler(field *Field, in *reflect.Value, data interface{}) (interface{}, error) {
//      return data, nil
//    }
//    New("json").OutHandler(handler)
func (t *Tag) OutFunction(handler OutHandlerF) *Tag {
	t.OutHandlerF = handler

	return t
}

// OutContract pass handler which implementing the contract for Out operation
//
//   type JsonHandler struct {
//     fields map[string]string
//   }
//   func (s *JsonHandler) Handle(field *Field, in *reflect.Value, data interface{}) (interface{}, error) {
//     s.fields[field.Name] = field.Type
//
//     return data, nil
//   }
//
//   New("json").OutHandlerContract(JsonHandler{})
//
func (t *Tag) OutContract(handler OutHandlerC) *Tag {
	t.OutHandlerC = handler

	return t
}

// Symbols set symbols for a tag
//   NewTag("name").Symbols(":", "|") // json:"a:b|c:d"
func (t *Tag) Symbols(keyValue, keysSeparator string) *Tag {
	t.TagSymbols = TagSymbols{
		KeyValue:      keyValue,
		KeysSeparator: keysSeparator,
	}

	return t
}

// New initialize of a tag
//  New("json")
func New(name string) *Tag {
	return &Tag{
		Name: name,
	}
}
