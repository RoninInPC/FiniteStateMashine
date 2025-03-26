package main

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	"reflect"
)

func isArray(value any) bool {
	return reflect.TypeOf(value).Kind() == reflect.Array || reflect.TypeOf(value).Kind() == reflect.Slice
}

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

func arrayToArrayAny(value any) []any {
	answer := make([]any, 0)
	if value == nil || !isArray(value) {
		return answer
	}
	val := reflect.ValueOf(value)
	for i := 0; i < val.Len(); i++ {
		part := val.Index(i).Interface()
		if isStruct(part) {
			answer = append(answer, structToArrayAny(part)...)
			continue
		}
		if isFunc(part) || isChannel(part) {
			answer = append(answer, reflect.ValueOf(part).Pointer())
			continue
		}
		if isArray(part) {
			answer = append(answer, arrayToArrayAny(part)...)
			continue
		}
	}
	return answer
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
		if isArray(field) {
			answer = append(answer, arrayToArrayAny(field)...)
			continue
		}
		answer = append(answer, field)
	}
	return answer
}

func anyToCorrectAny(value any) any {
	if isArray(value) {
		return arrayToArrayAny(value)
	}
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

type SerializationCondition[conditionType ConditionType] func(value conditionType) (uint64, error)

type SerializationAlphabet[alphabetType AlphabetType] func(value alphabetType) (uint64, error)

func SerializationConditionDefault[conditionType ConditionType](value conditionType) (uint64, error) {
	return serializationAny(value)
}

func SerializationAlphabetDefault[alphabetType AlphabetType](value alphabetType) (uint64, error) {
	return serializationAny(value)
}

func serializationAny(value any) (uint64, error) {
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
