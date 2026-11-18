package structs

import (
	"encoding/json"
	"reflect"
)

func ToUnstructured(v any) (any, error) {
	t := reflect.Indirect(reflect.ValueOf(v)).Type()
	switch {
	case t.Kind() == reflect.Struct:
		return toUnstructured(v)
	case t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Struct:
		return toUnstructured(v)
	default:
		return v, nil
	}
}

func toUnstructured(v any) (any, error) {
	var vv any
	buf, err := json.Marshal(v)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(buf, &vv)
	if err != nil {
		return v, err
	}
	return vv, nil
}
