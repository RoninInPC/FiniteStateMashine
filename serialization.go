package main

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	"reflect"
)

func isStruct(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Struct
}

func isFunc(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Func
}

func isChannel(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Chan
}

func channelToAny(value any) any {
	if value == nil || !isChannel(value) {
		return value
	}
	return reflect.ValueOf(value).Pointer()
}

func funcToAny(value any) any {
	if value == nil || !isFunc(value) {
		return value
	}
	return reflect.ValueOf(value).Pointer()
}

func structToArrayAny(value any) []any {
	answer := make([]any, 0)
	if value == nil || !isStruct(value) {
		return answer
	}
	valueOf := reflect.ValueOf(value)
	for i := 0; i < valueOf.Type().NumField(); i++ {
		field := valueOf.Field(i).Interface()
		if isStruct(field) {
			answer = append(answer, structToArrayAny(field)...)
			continue
		}
		if isFunc(field) || isChannel(field) {
			answer = append(answer, reflect.ValueOf(field).Pointer())
			continue
		}
		answer = append(answer, field)
	}
	return answer
}

func anyToCorrectAny(value any) any {
	if isStruct(value) {
		return structToArrayAny(value)
	}
	if isChannel(value) {
		return channelToAny(value)
	}
	if isFunc(value) {
		return funcToAny(value)
	}
	return value
}

type Serialization func(value any) (uint64, error)

func SerializationValueDefault(value any) (uint64, error) {
	correctAny := anyToCorrectAny(value)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(correctAny) // Сериализуем значение в байты
	if err != nil {
		return 0, err
	}

	h := fnv.New64a()
	h.Write(buf.Bytes()) // Вычисляем хэш от байтов
	return h.Sum64(), nil
}
