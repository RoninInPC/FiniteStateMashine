package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/fnv"
)

func hashValue(value any) uint32 {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value) // Сериализуем значение в байты
	if err != nil {
		panic(err)
	}

	h := fnv.New32a()
	h.Write(buf.Bytes()) // Вычисляем хэш от байтов
	return h.Sum32()
}

func main() {
	a := A{"21", 0}
	b := A{"21", 0}
	fmt.Println(hashValue(a), hashValue(b))
}
