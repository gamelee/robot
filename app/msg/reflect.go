package msg

import (
	"reflect"
	"strings"
)

var reqBody map[string]interface{}

func init() {
	reqBody = TypeDefault(reflect.TypeOf(Req{})).(map[string]interface{})
}

func GetReq(name string) interface{} {
	return reqBody[name]
}
func GetAllReqBody() map[string]interface{} {
	return reqBody
}

func TypeDefault(t reflect.Type) interface{} {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		tmpInfo := make(map[string]interface{})
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !strings.HasPrefix(f.Name, "XXX_") {
				tmpInfo[f.Name] = TypeDefault(f.Type)
			}
		}
		return tmpInfo
	} else if t.Kind() == reflect.Slice {
		tmpInfo := make([]interface{}, 1)
		tmpInfo[0] = TypeDefault(t.Elem())
		return tmpInfo
	} else if t.Kind() == reflect.Array {
		tmpInfo := make([]interface{}, t.Len())
		if t.Len() == 0 {
			return tmpInfo
		}
		for i := 0; i < t.Len(); i++ {
			tmpInfo[i] = TypeDefault(t.Elem())
		}
		tmpInfo[0] = TypeDefault(t.Elem())
		return tmpInfo
	} else {
		return reflect.New(t).Interface()
	}
}
