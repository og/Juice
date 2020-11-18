package jhttp

import (
	"reflect"
)

type RequestUnmarshaler interface {
	UnmarshalRequest()  (string ,error)
}
var requestUnmarshalerType = reflect.TypeOf((*RequestUnmarshaler)(nil)).Elem()

type RequestMarshaler interface {
	MarshalRequest(value string) error
}
var requestMarshalerType = reflect.TypeOf((*RequestMarshaler)(nil)).Elem()
