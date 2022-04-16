package tagger

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestFieldTag_TagsToParsed(t *testing.T) {
	cases := []struct {
		name    string
		tags    string
		symbols TagSymbols
		expect  func(values [][2]string) bool
	}{
		{
			name:    "Test empty tags",
			tags:    "",
			symbols: TagSymbols{":", ";"},
			expect: func(values [][2]string) bool {
				return assert.Len(t, values, 0)
			},
		},
		{
			name:    "Test tags without separators",
			tags:    "a|b|c",
			symbols: TagSymbols{":", ";"},
			expect: func(values [][2]string) bool {
				return assert.Len(t, values, 1) &&
					assert.Equal(t, [][2]string{{"a|b|c", "a|b|c"}}, values)
			},
		},
		{
			name:    "Test solo tag",
			tags:    "-",
			symbols: TagSymbols{":", ";"},
			expect: func(values [][2]string) bool {
				return assert.Len(t, values, 1) &&
					assert.Equal(t, [][2]string{{"-", "-"}}, values)
			},
		},
		{
			name:    "Test solo tag with one by key separator",
			tags:    "-;a",
			symbols: TagSymbols{":", ";"},
			expect: func(values [][2]string) bool {
				return assert.Len(t, values, 2) &&
					assert.Equal(t, [][2]string{{"-", "-"}, {"a", "a"}}, values)
			},
		},
		{
			name:    "Test solo tag with one by key separator with key value",
			tags:    "-;a:b",
			symbols: TagSymbols{":", ";"},
			expect: func(values [][2]string) bool {
				return assert.Len(t, values, 2) &&
					assert.Equal(t, [][2]string{{"-", "-"}, {"a", "b"}}, values)
			},
		},
		{
			name:    "Test two tags",
			tags:    "key:value;key2:value2",
			symbols: TagSymbols{":", ";"},
			expect: func(values [][2]string) bool {
				return assert.Len(t, values, 2) &&
					assert.Equal(t, [][2]string{{"key", "value"}, {"key2", "value2"}}, values)
			},
		},
	}

	for _, c := range cases {
		fieldTag := FieldTag{
			StructTag:  reflect.StructTag(c.tags),
			TagSymbols: c.symbols,
		}

		if !c.expect(fieldTag.TagsToParsed()) {
			t.Error(fmt.Sprintf("[TestFieldTag_TagsToParsed] %s is not true", c.name))
		}
	}
}

func TestFieldTag_Exists(t *testing.T) {
	cases := []struct {
		name    string
		tags    string
		symbols TagSymbols
		key     string
		expect  func(exists bool) bool
	}{
		{
			name:    "Must not exists",
			key:     "convert_to",
			tags:    "",
			symbols: TagSymbols{":", ";"},
			expect: func(exists bool) bool {
				return assert.False(t, exists)
			},
		},
		{
			name:    "Must exists",
			key:     "convert_to",
			tags:    "convert_to:value",
			symbols: TagSymbols{":", ";"},
			expect: func(exists bool) bool {
				return assert.True(t, exists)
			},
		},
	}

	for _, c := range cases {
		fieldTag := FieldTag{
			StructTag:  reflect.StructTag(c.tags),
			TagSymbols: c.symbols,
		}
		fieldTag.ParsedTags = fieldTag.TagsToParsed()

		if !c.expect(fieldTag.Exists(c.key)) {
			t.Error(fmt.Sprintf("[TestFieldTag_Exists] %s is not true", c.name))
		}
	}
}

func TestFieldTag_FindIndexByKey(t *testing.T) {
	cases := []struct {
		name    string
		tags    string
		symbols TagSymbols
		key     string
		expect  func(index int, exists bool) bool
	}{
		{
			name:    "Must not exists",
			key:     "convert_to",
			tags:    "",
			symbols: TagSymbols{":", ";"},
			expect: func(index int, exists bool) bool {
				return assert.Equal(t, -1, index) &&
					assert.False(t, exists)
			},
		},
		{
			name:    "Must be exists with 0 index",
			key:     "convert_to",
			tags:    "convert_to:value",
			symbols: TagSymbols{":", ";"},
			expect: func(index int, exists bool) bool {
				return assert.Equal(t, 0, index) &&
					assert.True(t, exists)
			},
		},
		{
			name:    "Must be exists with 1 index",
			key:     "convert_to",
			tags:    "key:value;convert_to:value",
			symbols: TagSymbols{":", ";"},
			expect: func(index int, exists bool) bool {
				return assert.Equal(t, 1, index) &&
					assert.True(t, exists)
			},
		},
	}

	for _, c := range cases {
		fieldTag := FieldTag{
			StructTag:  reflect.StructTag(c.tags),
			TagSymbols: c.symbols,
		}
		fieldTag.ParsedTags = fieldTag.TagsToParsed()

		if !c.expect(fieldTag.FindIndexByKey(c.key)) {
			t.Error(fmt.Sprintf("[FindIndexByKey] %s is not true", c.name))
		}
	}
}

func TestFieldTag_FindByKey(t *testing.T) {
	cases := []struct {
		name    string
		tags    string
		symbols TagSymbols
		key     string
		expect  func(value string, exists bool) bool
	}{
		{
			name:    "Must not exists",
			key:     "convert_to",
			tags:    "",
			symbols: TagSymbols{":", ";"},
			expect: func(value string, exists bool) bool {
				return assert.Len(t, value, 0) &&
					assert.False(t, exists)
			},
		},
		{
			name:    "Must be exists",
			key:     "convert_to",
			tags:    "convert_to:string",
			symbols: TagSymbols{":", ";"},
			expect: func(value string, exists bool) bool {
				return assert.Equal(t, "string", value) &&
					assert.True(t, exists)
			},
		},
	}

	for _, c := range cases {
		fieldTag := FieldTag{
			StructTag:  reflect.StructTag(c.tags),
			TagSymbols: c.symbols,
		}
		fieldTag.ParsedTags = fieldTag.TagsToParsed()

		if !c.expect(fieldTag.FindByKey(c.key)) {
			t.Error(fmt.Sprintf("[FindByKey] %s is not true", c.name))
		}
	}
}
