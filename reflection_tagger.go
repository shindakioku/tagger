package tagger

import (
	"errors"
	"fmt"
	"reflect"
)

// Tags [name] -> Tag
type Tags map[string]*Tag

type ReflectionTagger struct {
	tags Tags
}

func (r *ReflectionTagger) Add(tag *Tag) Tagger {
	r.tags[tag.Name] = tag

	return r
}

// Collect defined tags when passed some tags for process or use all defined
func (r ReflectionTagger) collectCorrectTags(tagsForProcess []string) (Tags, error) {
	tagsForWork := make(Tags)
	if len(tagsForProcess) == 0 {
		tagsForWork = r.tags
	} else {
		for _, tag := range tagsForProcess {
			t, exists := r.tags[tag]
			if !exists {
				return tagsForWork, errors.New(fmt.Sprintf("%s tag doesn't exists. Check provided arguments please.", tag))
			}

			tagsForWork[tag] = t
		}
	}

	return tagsForWork, nil
}

// Create reflection value from user input with some validation (must be only struct or ptr)
func (r ReflectionTagger) createValueOf(value any) (reflect.Value, error) {
	var valueOf reflect.Value
	if len(r.tags) == 0 {
		return valueOf, errors.New("you must register at least one tag")
	}

	if value == nil {
		return valueOf, errors.New("%in% cannot be empty")
	}

	valueOf = reflect.ValueOf(value)
	if !containsInSlice(valueOf.Kind(), []reflect.Kind{reflect.Struct, reflect.Ptr}) {
		return valueOf, errors.New("%in% must be a structure")
	}

	for valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}

	if valueOf.NumField() == 0 {
		return valueOf, errors.New("%in% required at least one field")
	}

	return valueOf, nil
}

// Collect all of a struct fields
func (r ReflectionTagger) collectFields(valueOf reflect.Value, typeOf reflect.Type, parentStruct *ParentStruct, needToSet bool) (fields []*Field) {
	for i := 0; i < valueOf.NumField(); i++ {
		v := valueOf.Field(i)

		// Nested struct with reference
		// We must create struct object and set it to the field before we can use it
		// needToSet - we don't need to rewrite user fields if it's [Out] action
		if v.Kind() == reflect.Ptr && needToSet {
			v.Set(reflect.New(v.Type().Elem()))
		}

		// Nested struct without pointer
		// We must create value which referencing to the user struct but with pointer
		// otherwise we can't do anything (for example: we can't set the value)
		if v.Kind() == reflect.Struct {
			v = reflect.ValueOf(v.Addr().Interface())
		}

		structField := typeOf.Field(i)

		fields = append(fields, &Field{
			Value:        v,
			StructField:  structField,
			ParentStruct: parentStruct,
			Index:        i,
			IsStruct:     isStruct(v),
			Tag:          NewFieldTag(structField.Tag),
			tags:         make(Tags),
		})
	}

	return
}

// Check is handler for fields without tag are defined
func (r ReflectionTagger) getHandlerForEmptyTag(name string, availableTags Tags) (handler *Tag, err error) {
	handler, handlerForEmptyFieldExists := availableTags[name]
	if !handlerForEmptyFieldExists {
		err = errors.New(fmt.Sprintf("%s doesn't exists (empty field parser)", name))
	}

	return
}

// Collect all available handler for a field
func (r ReflectionTagger) makeHandlersForField(field *Field, handlerForEmptyField *Tag, tagsForWork Tags) Tags {
	handlers := make(Tags)

	if field.Tag.IsEmpty() && handlerForEmptyField != nil {
		handlers[handlerForEmptyField.Name] = handlerForEmptyField

		return handlers
	}

	for name, handler := range tagsForWork {
		_, exists := field.StructField.Tag.Lookup(name)
		if !exists {
			continue
		}

		handlers[name] = handler
	}

	return handlers
}

func (r ReflectionTagger) In(data any, in any, tagForEmpty string, tags ...string) error {
	return r.in(nil, data, in, tagForEmpty, tags...)
}

func (r ReflectionTagger) in(parentStruct *ParentStruct, data any, in any, tagForEmpty string, tags ...string) error {
	tagsForWork, err := r.collectCorrectTags(tags)
	if err != nil {
		return err
	}

	valueOf, err := r.createValueOf(in)
	if err != nil {
		return err
	}

	typeOf := reflect.TypeOf(in)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	fields := r.collectFields(valueOf, typeOf, parentStruct, true)

	var handlerForEmptyField *Tag
	if len(tagForEmpty) > 0 {
		handlerForEmptyField, err = r.getHandlerForEmptyTag(tagForEmpty, tagsForWork)
		if err != nil {
			return err
		}
	}

	// Evaluate
	for i, field := range fields {
		// TODO: Reimplement more effective
		fields[i].tags = r.makeHandlersForField(field, handlerForEmptyField, tagsForWork)
	}

	for _, field := range fields {
		if err = r.callInHandlers(data, field.tags, field, &valueOf); err != nil {
			return err
		}

		if field.IsStruct {
			parentS := &ParentStruct{Value: valueOf, ParentField: field}
			if err = r.in(parentS, data, field.Value.Interface(), tagForEmpty, tags...); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r ReflectionTagger) nestedStructWithPointerToInterfaceType(value reflect.Value) any {
	pointer := reflect.PointerTo(value.Type().Elem())
	isPtr := pointer.Kind() == reflect.Ptr
	for isPtr {
		pointer = pointer.Elem()
		isPtr = pointer.Kind() == reflect.Ptr
	}

	return reflect.ValueOf(pointer).Interface()
}

func (r ReflectionTagger) Out(data any, out interface{}, tagForEmpty string, tags ...string) (output interface{}, err error) {
	return r.out(nil, data, out, tagForEmpty, tags...)
}

func (r ReflectionTagger) out(parentStruct *ParentStruct, data any, out interface{}, tagForEmpty string, tags ...string) (output interface{}, err error) {
	output = data
	tagsForWork, err := r.collectCorrectTags(tags)
	if err != nil {
		return
	}

	valueOf, err := r.createValueOf(out)
	if err != nil {
		return
	}

	typeOf := reflect.TypeOf(out)
	for typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}

	fields := r.collectFields(valueOf, typeOf, parentStruct, false)

	var handlerForEmptyField *Tag
	if len(tagForEmpty) > 0 {
		handlerForEmptyField, err = r.getHandlerForEmptyTag(tagForEmpty, tagsForWork)
		if err != nil {
			return
		}
	}

	// Evaluate
	for i, field := range fields {
		// TODO: Reimplement more effective
		fields[i].tags = r.makeHandlersForField(field, handlerForEmptyField, tagsForWork)
	}

	for _, field := range fields {
		if output, err = r.callOutHandlers(output, field.tags, field, &valueOf); err != nil {
			return
		}

		if field.IsStruct {
			parentS := &ParentStruct{Value: valueOf, ParentField: field}
			if output, err = r.out(parentS, output, field.Value.Interface(), tagForEmpty, tags...); err != nil {
				return
			}
		}
	}

	return
}

func (r ReflectionTagger) callOutHandlers(data any, tags Tags, field *Field, in *reflect.Value) (output interface{}, err error) {
	output = data
	for _, handler := range tags {
		if handler.OutHandlerF == nil && handler.OutHandlerC == nil {
			return output, errors.New("empty")
		}

		field.Tag.TagSymbols = handler.TagSymbols
		field.Tag.ParsedTags = field.Tag.TagsToParsed()
		if handler.OutHandlerF != nil {
			if output, err = handler.OutHandlerF(data, field, in); err != nil {
				return
			}

			continue
		}

		if output, err = handler.OutHandlerC.Handle(data, field, in); err != nil {
			return
		}
	}

	return
}

func (r ReflectionTagger) callInHandlers(data any, tags Tags, field *Field, in *reflect.Value) error {
	for _, handler := range tags {
		if handler.InHandlerF == nil && handler.InHandlerC == nil {
			return errors.New("empty")
		}

		field.Tag.TagSymbols = handler.TagSymbols
		field.Tag.ParsedTags = field.Tag.TagsToParsed()
		if handler.InHandlerF != nil {
			if err := handler.InHandlerF(data, field, in); err != nil {
				return err
			}

			continue
		}

		if err := handler.InHandlerC.Handle(data, field, in); err != nil {
			return err
		}
	}

	return nil
}

func NewReflectionTagger() Tagger {
	return &ReflectionTagger{
		tags: make(map[string]*Tag),
	}
}
