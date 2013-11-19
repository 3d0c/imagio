package utils

import (
	"log"
	"reflect"
)

func Construct(i interface{}, args ...interface{}) interface{} {
	v := reflect.ValueOf(i)

	if method := v.MethodByName("Construct"); method.IsValid() {
		return method.Call([]reflect.Value{reflect.ValueOf(args)})[0].Interface()
	}

	log.Println("Type", reflect.TypeOf(i), "doesn't have `Construct` method.")

	return nil
}
