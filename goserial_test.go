package goserial

import (
	"log"
	"testing"
)

type T struct {
	A int
	B float32
}

func TestSerial(t *testing.T) {
	b := Serial(3.14)
	log.Println(b)
}

func TestSerialStruct(t *testing.T) {
	b := Serial(T{A: 3, B: 3.14})
	log.Println(b)
}

func BenchmarkSerialStructPerf(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Serial(T{A: 3, B: 3.14})
	}
	// log.Println(b)
}

func TestSerialPtr(t *testing.T) {
	b := Serial(&T{A: 3, B: 3.14})
	log.Println(b)
}

func TestDeserial(t *testing.T) {
	b := Serial(&T{A: 3, B: 3.14})
	log.Println(b)

	var s T
	Deserial(b, &s)
}
