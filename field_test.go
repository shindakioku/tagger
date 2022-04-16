package tagger

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestField_Set(t *testing.T) {
	cases := []struct {
		name        string
		value       any
		valueForSet any
		expect      func(value any, err error) bool
	}{
		{
			name: "Test non pointer struct with error 'cannot set value'",
			value: struct {
				Name string
			}{},
			valueForSet: "",
			expect: func(value any, err error) bool {
				return assert.NotNil(t, err) &&
					assert.Equal(t, errors.New("Can't set value for field: Name"), err)
			},
		},
		{
			name: "Test difference type of passed value with field type",
			value: &struct {
				ID int
			}{},
			valueForSet: "foo",
			expect: func(value any, err error) bool {
				return assert.NotNil(t, err) &&
					assert.Equal(t, errors.New("Incorrect type for field: ID. Your: string. Actual: int"), err)
			},
		},
		{
			name: "Test correctly work for int",
			value: &struct {
				ID int
			}{},
			valueForSet: 1,
			expect: func(value any, err error) bool {
				return assert.Nil(t, err) &&
					assert.Equal(t, 1, value.(*struct{ ID int }).ID)
			},
		},
		{
			name: "Test correctly work for string",
			value: &struct {
				Username string
			}{},
			valueForSet: "foo",
			expect: func(value any, err error) bool {
				return assert.Nil(t, err) &&
					assert.Equal(t, "foo", value.(*struct{ Username string }).Username)
			},
		},
		{
			name:        "Test correctly work for nested struct",
			value:       &struct{ Foo struct{ ID int } }{},
			valueForSet: struct{ ID int }{ID: 1},
			expect: func(value any, err error) bool {
				return assert.Nil(t, err) &&
					assert.Equal(t, 1, value.(*struct{ Foo struct{ ID int } }).Foo.ID)
			},
		},
	}

	for _, c := range cases {
		valueOf := reflect.ValueOf(c.value)
		if valueOf.Kind() == reflect.Ptr {
			valueOf = valueOf.Elem()
		}

		typeOf := reflect.TypeOf(c.value)
		if typeOf.Kind() == reflect.Ptr {
			typeOf = typeOf.Elem()
		}

		field := Field{Value: valueOf.Field(0), StructField: typeOf.Field(0)}
		err := field.Set(c.valueForSet)
		if !c.expect(c.value, err) {
			t.Error(fmt.Sprintf("[TestField_Set] %s is not true", c.name))
		}
	}
}
