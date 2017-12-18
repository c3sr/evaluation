package evaluation

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

func uptoIndex(arry []interface{}, idx int) int {
	if len(arry) <= idx {
		return len(arry) - 1
	}
	return idx
}

func toFloat64Slice(i interface{}) []float64 {
	res, _ := toFloat64SliceE(i)
	return res
}

func toFloat64SliceE(i interface{}) ([]float64, error) {
	if i == nil {
		return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
	}

	switch v := i.(type) {
	case []float64:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]float64, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToFloat64E(s.Index(j).Interface())
			if err != nil {
				return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
			}
			a[j] = val
		}
		return a, nil
	default:
		return []float64{}, fmt.Errorf("unable to cast %#v of type %T to []float64", i, i)
	}
}

func float64SliceToStringSlice(us []float64) []string {
	res := make([]string, len(us))
	for ii, u := range us {
		res[ii] = cast.ToString(u)
	}
	return res
}

func uint64SliceToStringSlice(us []uint64) []string {
	res := make([]string, len(us))
	for ii, u := range us {
		res[ii] = cast.ToString(u)
	}
	return res
}
