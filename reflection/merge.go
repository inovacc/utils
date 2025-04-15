package reflection

import (
	"reflect"
)

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
