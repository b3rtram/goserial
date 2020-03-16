package goserial

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"math"
	"reflect"
)

//Serial will serialize all types of structs or pointer of structs in a binary format
func Serial(obj interface{}) []byte {

	var b bytes.Buffer

	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	serialize(t, v, &b)

	// fmt.Println(reflect.ValueOf(obj))

	return b.Bytes()
}

func serialize(t reflect.Type, v reflect.Value, b *bytes.Buffer) error {

	// fmt.Println(t)

	k := t.Kind()

	switch k {
	case reflect.Struct:
		b.Write([]byte{0x4c, 0x8c})
		b.Write([]byte(t.Name()))
		b.Write([]byte{0})
		for i := 0; t.NumField() > i; i++ {
			f := t.Field(i)
			nv := v.Field(i)
			typeSerialize(f, nv, b)
		}

		b.Write([]byte{0xc8, 0xc4})

	case reflect.Ptr:
		e := t.Elem()
		ve := v.Elem()
		serialize(e, ve, b)
	default:
		return errors.New("could not deserialize not struct types")
	}

	return nil
}

func typeSerialize(s reflect.StructField, v reflect.Value, b *bytes.Buffer) {

	log.Println(s.Name)
	switch s.Type.Kind() {
	case reflect.Int:
		b.Write([]byte{0xaa})
		b.Write([]byte(s.Name))
		b.Write([]byte{0})
		i := v.Int()
		binary.Write(b, binary.LittleEndian, encodeInt(i))
		b.Write([]byte{0})

	case reflect.Float64:
		b.Write([]byte{0xab})
		b.Write([]byte(s.Name))
		b.Write([]byte{0})
		f := v.Float()
		b.Write(Float64bytes(f))
		b.Write([]byte{0})

	case reflect.String:
		b.Write([]byte{0xac})
		b.Write([]byte(s.Name))
		b.Write([]byte{0})
		b.Write([]byte(v.String()))
		b.Write([]byte{0})

	// case reflect.Slice:
	// 	b.Write([]byte(t.Name()))
	// 	b.Write([]byte{0})

	// 	// v.Slice()
	// case reflect.Map:
	// 	b.Write([]byte(t.Name()))
	// 	b.Write([]byte{0})

	// 	// v.Map
	// case reflect.Array:
	// 	b.Write([]byte(t.Name()))
	// 	b.Write([]byte{0})

	case reflect.Bool:
		b.Write([]byte{0xad})
		b.Write([]byte(s.Name))
		b.Write([]byte{0})
		binary.Read(b, binary.LittleEndian, v.Bool())
		b.Write([]byte{0})

	}

}

func Deserial(b []byte, obj interface{}) error {

	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	log.Println(t.Kind())

	if t.Kind() != reflect.Ptr {
		return errors.New("could not deserialize into non pointer struct")
	}

	e := t.Elem()
	ve := v.Elem()

	log.Println(e.Kind())

	lb := b

	if lb[0] != 0x4c && lb[1] != 0x8c {
		return errors.New("not a goserial binary")
	}

	var split []bytes.Buffer
	a := 0
	split = append(split, *bytes.NewBuffer(nil))
	for i := 2; len(lb) > i; i++ {

		if lb[i] == 0 {
			a++
			split = append(split, *bytes.NewBuffer(nil))
		} else {
			split[a].Write([]byte{lb[i]})
		}
	}

	typeName := string(split[0].Bytes())
	if typeName != t.Name() {
		errors.New("Typename not equal")
	}

	for i := 1; i < len(split); i++ {
		fieldname := string(split[i].Bytes())
		log.Println(e.Kind())
		st, _ := e.FieldByName(fieldname)
		f := ve.FieldByName(string(split[i].Bytes()))

		i = i + 1

		switch st.Type.Kind() {
		case reflect.Int:
			bs := split[i].Bytes()
			bla := binary.LittleEndian.Uint64(bs)
			in := toInt(bla)
			f.SetInt(in)
			log.Printf("%d\n", in)

		case reflect.Float64:

		case reflect.Float32:

		case reflect.String:

		// case reflect.Slice:

		// 	// v.Slice()
		// case reflect.Map:

		// 	// v.Map
		// case reflect.Array:

		case reflect.Bool:

		}

	}

	return nil
}

func encodeInt(i int64) uint64 {
	var x uint64
	if i < 0 {
		x = uint64(^i<<1) | 1
	} else {
		x = uint64(i << 1)
	}

	return x
}

func toInt(x uint64) int64 {
	i := int64(x >> 1)
	if x&1 != 0 {
		i = ^i
	}
	return i
}

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
