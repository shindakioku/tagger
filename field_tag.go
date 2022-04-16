package tagger

import (
	"reflect"
	"strings"
)

// FieldTag defined tag for the field.
type FieldTag struct {
	// reflection.Tag
	StructTag reflect.StructTag
	// Tag - data by the handler name (every handler must have unique name)
	//StructTagByKey string
	// Your defined symbols
	TagSymbols TagSymbols

	// Parsed tags by keys
	// tuple (0 - key name, 1 - value)
	ParsedTags [][2]string
}

// Values split field tags by key separator
// "foo | bar | baz" - key separator is '|'
func (t FieldTag) Values() []string {
	return strings.Split(t.ToString(), t.TagSymbols.KeysSeparator)
}

// TagsToParsed
// returns tuple (0 - key name, 1 - value)
func (t FieldTag) TagsToParsed() (output [][2]string) {
	for _, value := range t.Values() {
		exploded := strings.Split(value, t.TagSymbols.KeyValue)
		if len(exploded) == 0 {
			continue
		}

		if len(exploded) == 1 {
			if len(value) == 0 {
				continue
			}

			output = append(output, [2]string{strings.ToLower(value), exploded[0]})

			continue
		}

		output = append(output, [2]string{strings.ToLower(exploded[0]), exploded[1]})
	}

	return
}

// IsEmpty does field contains any tags
func (t FieldTag) IsEmpty() bool {
	return len(t.StructTag) == 0
}

func (t FieldTag) ToString() string {
	return string(t.StructTag)
}

// Exists returns true if key defined in the field tags
func (t FieldTag) Exists(key string) bool {
	key = strings.ToLower(key)
	for _, v := range t.ParsedTags {
		if v[0] == key {
			return true
		}
	}

	return false
}

// FindIndexByKey returns index and exists (bool) if key defined in the fields tags
func (t FieldTag) FindIndexByKey(key string) (int, bool) {
	key = strings.ToLower(key)
	for i, v := range t.ParsedTags {
		if v[0] == key {
			return i, true
		}
	}

	return -1, false
}

// FindByKey returns value and exists (bool) if key defined in the fields tags
func (t FieldTag) FindByKey(key string) (string, bool) {
	key = strings.ToLower(key)
	for _, v := range t.ParsedTags {
		if v[0] == key {
			return v[1], true
		}
	}

	return "", false
}

func NewFieldTag(tag reflect.StructTag) FieldTag {
	return FieldTag{
		StructTag: tag,
	}
}
