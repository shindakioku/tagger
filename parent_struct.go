package tagger

import "reflect"

// ParentStruct struct can contain nested struct.
// You must handle only fields, but you may need access to parent struct
type ParentStruct struct {
	Value       reflect.Value
	ParentField *Field
}
