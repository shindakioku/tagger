package tagger

import (
	"reflect"
)

func containsInSlice[T comparable](need T, set []T) bool {
	for _, v := range set {
		if need == v {
			return true
		}
	}

	return false
}

func isStruct(valueOf reflect.Value) bool {
	if valueOf.Kind() != reflect.Ptr {
		return valueOf.Kind() == reflect.Struct
	}

	return reflect.TypeOf(valueOf.Interface()).Elem().Kind() == reflect.Struct
}

func ContainsInSlice[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
