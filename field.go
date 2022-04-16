package tagger

import (
	"errors"
	"fmt"
	"reflect"
)

// Field concrete field of the struct.
// Provided easier api than reflection.
type Field struct {
	// reflection.ValueOf(field)
	Value reflect.Value
	// reflection.TypeOf(field)
	StructField reflect.StructField
	// Index of the field
	Index int
	// If is nested struct
	IsStruct     bool
	ParentStruct *ParentStruct

	Tag FieldTag

	// For call
	tags Tags
}

func (f Field) Get() any {
	return f.Value.Interface()
}

func (f Field) GetFromPointer() any {
	return f.Value.Elem().Interface()
}

func (f Field) IsPtr() bool {
	return f.Value.Kind() == reflect.Ptr
}

func (f Field) Name() string {
	return f.StructField.Name
}

func (f Field) Type() reflect.Kind {
	return f.StructField.Type.Kind()
}

// Set value to the struct field
func (f *Field) Set(value interface{}) error {
	if !f.Value.CanSet() {
		return errors.New(fmt.Sprintf(
			"Can't set value for field: %s", f.StructField.Name,
		))
	}

	fieldType := f.StructField.Type.String()
	valueType := reflect.TypeOf(value).String()
	if fieldType != valueType {
		return errors.New(fmt.Sprintf(
			"Incorrect type for field: %s. Your: %s. Actual: %s",
			f.StructField.Name,
			valueType,
			f.StructField.Type,
		))
	}

	f.Value.Set(reflect.ValueOf(value))

	return nil
}
