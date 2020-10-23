package utils

import (
    "reflect"
    "unsafe"
)

// GetUnexportedField >>> UNSAFE
//
// Should not be used.
func GetUnexportedField(field reflect.Value) interface{} {
    return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

// func SetUnexportedField(field reflect.Value, value interface{}) {
//     reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
//         Elem().
//         Set(reflect.ValueOf(value))
// }

