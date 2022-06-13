package util

import (
	"github.com/streadway/amqp"
	"reflect"
	"unsafe"
)

type RabbitMQMsgHandleFunc func(amqp.Delivery)

func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func Must(err error) bool {
	if err != nil {
		panic(err)
	}
	return false
}

func FindStringInGeneric(slice []interface{}, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}



