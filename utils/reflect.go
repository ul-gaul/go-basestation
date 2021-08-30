package utils

import (
	"reflect"
	"unsafe"
)

// GetUnexportedField >>> UNSAFE
//
// Avoid using this function.
func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(),
		unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

// SetUnexportedField >>> UNSAFE
//
// Avoid using this function.
func SetUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(),
		unsafe.Pointer(field.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
}
