package reflection

import (
	"reflect"
)

// MergeZeroFields copies non-zero fields from src to dst, but only for fields in dst that are currently zero.
// dst and src must be pointers to structs of the same type.
//
// Example:
//
//	MergeZeroFields(&userConfig, &defaultConfig)
//
// Only fields in userConfig that are zero-valued will be replaced by corresponding values from defaultConfig.
func MergeZeroFields(dst, src any) {
	vDst := reflect.ValueOf(dst).Elem()
	vSrc := reflect.ValueOf(src).Elem()

	for i := 0; i < vDst.NumField(); i++ {
		fieldDst := vDst.Field(i)
		fieldSrc := vSrc.Field(i)

		// Only set if field is zero (e.g., 0, "", nil)
		if fieldDst.CanSet() && fieldDst.IsZero() {
			fieldDst.Set(fieldSrc)
		}
	}
}

// CopyExportedFields copies all exported fields from src to dst, regardless of their current values.
// Useful for copying full configurations or cloning objects.
//
// Note: Unexported fields (lowercase) will be ignored.
func CopyExportedFields(dst, src any) {
	vDst := reflect.ValueOf(dst).Elem()
	vSrc := reflect.ValueOf(src).Elem()
	tSrc := vSrc.Type()

	for i := 0; i < tSrc.NumField(); i++ {
		field := tSrc.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}
		fieldDst := vDst.FieldByName(field.Name)
		if fieldDst.IsValid() && fieldDst.CanSet() {
			fieldDst.Set(vSrc.Field(i))
		}
	}
}

// ZeroStruct resets all settable fields in a struct to their zero value.
// Takes a pointer to a struct.
//
// Example:
//
//	ZeroStruct(&config) // config.Name = "", config.Count = 0, config.Enabled = false, etc.
func ZeroStruct(v any) {
	val := reflect.ValueOf(v).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.CanSet() {
			field.Set(reflect.Zero(field.Type()))
		}
	}
}

// StructToMap converts an exported struct (or a pointer to one) into a map[string]any.
// Only exported fields will be included in the map.
//
// Useful for logging, serializing, or dynamic field access.
//
// Example:
//
//	m := StructToMap(user)
//	fmt.Println(m["Name"], m["Active"])
func StructToMap(input any) map[string]any {
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	result := make(map[string]any)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath == "" { // exported
			result[field.Name] = val.Field(i).Interface()
		}
	}
	return result
}
