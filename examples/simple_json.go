package examples

import (
	"errors"
	"fmt"
	"github.com/shindakioku/tagger"
	"log"
	"reflect"
	"regexp"
	"strconv"
)

type Profile struct {
	Email *string `my_json:"email"`
}

type User struct {
	ID       uint     `my_json:"user_id"`
	Username *string  `my_json:"username"`
	Profile  Profile  `my_json:"profile"`
	Profile2 *Profile `my_json:"profile2"`
}

func MyJsonIn(data any, field *tagger.Field, in *reflect.Value) error {
	tags := field.Tag.Values()
	if len(tags) != 1 {
		return errors.New("too many values for tag")
	}

	var pattern string
	if field.Type() == reflect.Uint {
		pattern = fmt.Sprintf(`"%s":%s`, tags[0], `([\d+])`)
	} else {
		pattern = fmt.Sprintf(`"%s":%s`, tags[0], `"(.*?)"`)
	}

	// Skip User
	if field.IsStruct {
		return nil
	}

	if field.ParentStruct != nil {
		pattern = fmt.Sprintf(
			`"%s":{"%s":%s}`,
			field.ParentStruct.ParentField.Tag.Values()[0],
			tags[0],
			`"(.*?)"`,
		)
	}

	match := regexp.MustCompile(pattern).FindStringSubmatch(data.(string))
	if len(match) == 0 {
		return errors.New(fmt.Sprintf("Can't find key in json: %s", tags[0]))
	}

	if field.Type() == reflect.Uint {
		v, _ := strconv.Atoi(match[1])

		return field.Set(uint(v))
	}

	if field.IsPtr() {
		return field.Set(&match[1])
	}

	return field.Set(match[1])
}

type JsonOut struct {
	Fields map[string]any
	Output string
}

func MyJsonOut(data any, field *tagger.Field, in *reflect.Value) (any, error) {
	if field.IsStruct {
		return data, nil
	}

	value := field.Get()
	if field.IsPtr() {
		value = field.GetFromPointer()
	}

	key := field.Tag.Values()[0]
	if field.ParentStruct != nil {
		parentKey := field.ParentStruct.ParentField.Tag.Values()[0]
		data.(*JsonOut).Fields[parentKey] = make(map[string]any)
		data.(*JsonOut).Fields[parentKey].(map[string]any)[key] = value
	} else {
		data.(*JsonOut).Fields[key] = value
	}

	return data, nil
}

func simpleJsonInOut() {
	tag := tagger.NewReflectionTagger().
		Add(tagger.
			New("my_json").
			InFunction(MyJsonIn).
			OutFunction(MyJsonOut).
			Symbols(";", ","),
		)

	username := "foo"
	email := "foo@gmail.com"
	user := &User{
		ID:       1,
		Username: &username,
		Profile: Profile{
			Email: &email,
		},
		Profile2: &Profile{
			Email: &email,
		},
	}
	log.Println(tag.Out(&JsonOut{Fields: make(map[string]any)}, &user, ""))

	var userIn User
	log.Println(tag.In(
		`"user_id":1,"username":"foo","profile":{"email":"foo@gmail.com"},"profile2":{"email":"11"}`,
		&userIn,
		"",
	))
	log.Println(userIn)
}
