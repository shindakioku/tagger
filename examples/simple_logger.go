package examples

import (
	"github.com/shindakioku/tagger"
	"log"
	"net/http"
	"reflect"
)

type RequestData struct {
	HTTPStatusCode int
	Headers        map[string]string
}

type MyDataForLogging struct {
	ParsedBody  []byte       `my_logger:"key:summary | to:string"`
	RequestData *RequestData `my_logger:"key:data"`
	IgnoreField string       `my_logger:"-"`
}

type LoggerMessage struct {
	Summary string
	Data    map[string]any
}

type MyLoggerOutData struct {
	LoggerMessage LoggerMessage
}

func MyLoggerOut(data any, field *tagger.Field, in *reflect.Value) (any, error) {
	if field.IsStruct {
		return data, nil
	}

	value := field.Get()
	if field.IsPtr() {
		value = field.GetFromPointer()
	}

	myLoggerOutData := data.(*MyLoggerOutData)
	if field.Tag.IsEmpty() && field.ParentStruct != nil {
		parentTags := field.ParentStruct.ParentField.Tag
		keyValue, exists := parentTags.FindByKey("key")
		if exists && keyValue == "data" {
			if myLoggerOutData.LoggerMessage.Data == nil {
				myLoggerOutData.LoggerMessage.Data = make(map[string]any)
			}

			myLoggerOutData.LoggerMessage.Data[field.Name()] = value
		}

		return data, nil
	}

	keyValue, exists := field.Tag.FindByKey("key")
	if exists && keyValue == "summary" {
		convertToType, exists := field.Tag.FindByKey("to")
		if exists {
			if convertToType == "string" {
				myLoggerOutData.LoggerMessage.Summary = string(value.([]byte))
			}
		}
	}

	return data, nil
}

func main() {
	tag := tagger.NewReflectionTagger().
		Add(tagger.
			New("my_logger").
			OutFunction(MyLoggerOut).
			Symbols(":", " | "),
		)

	req := MyDataForLogging{
		RequestData: &RequestData{
			HTTPStatusCode: http.StatusOK,
			Headers: map[string]string{
				"Accept": "application/json",
			},
		},
		ParsedBody:  []byte("hello world!"),
		IgnoreField: "...",
	}
	log.Println(tag.Out(&MyLoggerOutData{}, &req, "my_logger"))
}
