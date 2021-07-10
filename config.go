/*
Package config is a struct-based wrapper
for YAML (de)serialization
*/
package config

import (
	"errors"
	"io"
	"reflect"

	"gopkg.in/yaml.v2"
)

var ErrFieldUnsettable = errors.New("struct field not settable")
var ErrInvalidDataType = errors.New("use of invalid data type")

// Marshal writes a configuration to the specified Writer
func Marshal(w io.Writer, v interface{}) error {
	yml := marshal(v).(map[interface{}]interface{})
	if r := recover(); r != nil {
		if _, ok := r.(error); ok {
			return r.(error)
		}
	}

	out, err := yaml.Marshal(yml)
	if err != nil {
		return err
	}

	_, err = w.Write(out)
	return err
}

// Unmarshal reads a configuration from the specified Reader
func Unmarshal(r io.Reader, v interface{}) error {
	yml, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	ymap := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(yml, ymap); err != nil {
		return err
	}

	unmarshal(ymap, reflect.Indirect(reflect.ValueOf(v)))

	if r := recover(); r != nil {
		if _, ok := r.(error); ok {
			return r.(error)
		}
	}
	return nil
}

func marshal(v interface{}) interface{} {
	rt := reflect.TypeOf(v)
	switch rt.Kind() {
	case reflect.Struct:
		m := make(map[interface{}]interface{})
		rv := reflect.ValueOf(v)

		for i := 0; i < rv.NumField(); i++ {
			ft := rt.Field(i).Tag.Get("conf")
			fv := rv.Field(i).Interface()
			m[ft] = marshal(fv)
		}

		return m
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.UnsafePointer:
		panic(ErrInvalidDataType)
	default:
		return v
	}
}

func unmarshal(yml interface{}, rv reflect.Value) {
	rt := rv.Type()

	switch rv.Kind() {
	case reflect.Array:
		for i := 0; i < rv.Cap(); i++ {
			fv := reflect.Indirect(reflect.New(rv.Index(i).Type()))
			unmarshal(yml.([]interface{})[i], fv)

			rv.Index(i).Set(fv)
		}
	case reflect.Slice:
		yarr := yml.([]interface{})
		for i := 0; i < len(yarr); i++ {
			sv := reflect.MakeSlice(rv.Type(), 1, 1)
			fv := reflect.Indirect(reflect.New(sv.Index(0).Type()))
			unmarshal(yarr[i], fv)

			rv.Set(reflect.Append(rv, fv))
		}
	case reflect.Map:
		rv.Set(reflect.MakeMap(rv.Type()))

		ymap := yml.(map[interface{}]interface{})
		for k, v := range ymap {
			fv := reflect.Indirect(reflect.New(rv.Type().Elem()))
			unmarshal(v, fv)

			rv.SetMapIndex(reflect.ValueOf(k), fv)
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			ft := rt.Field(i).Tag.Get("conf")
			fv := reflect.Indirect(reflect.New(rv.Field(i).Type()))
			unmarshal(yml.(map[interface{}]interface{})[ft], fv)

			rv.Field(i).Set(fv)
		}
	case reflect.Complex64:
		fallthrough
	case reflect.Complex128:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Ptr:
		fallthrough
	case reflect.UnsafePointer:
		panic(ErrInvalidDataType)
	case reflect.Int8:
		rv.Set(reflect.ValueOf(int8(yml.(int))))
	case reflect.Int16:
		rv.Set(reflect.ValueOf(int16(yml.(int))))
	case reflect.Int32:
		rv.Set(reflect.ValueOf(int32(yml.(int))))
	case reflect.Int64:
		rv.Set(reflect.ValueOf(int64(yml.(int))))
	case reflect.Uint:
		rv.Set(reflect.ValueOf(uint(yml.(int))))
	case reflect.Uint8:
		rv.Set(reflect.ValueOf(uint8(yml.(int))))
	case reflect.Uint16:
		rv.Set(reflect.ValueOf(uint16(yml.(int))))
	case reflect.Uint32:
		rv.Set(reflect.ValueOf(uint32(yml.(int))))
	case reflect.Uint64:
		rv.Set(reflect.ValueOf(uint64(yml.(int))))
	case reflect.Uintptr:
		rv.Set(reflect.ValueOf(uintptr(yml.(int))))
	case reflect.Float32:
		rv.Set(reflect.ValueOf(float32(yml.(float64))))
	default:
		rv.Set(reflect.ValueOf(yml))
	}
}
