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
	case reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.Ptr, reflect.UnsafePointer:
		panic(ErrInvalidDataType)
	default:
		return v
	}
}

func unmarshal(yml interface{}, rv reflect.Value) {
	rt := rv.Type()

	switch rv.Kind() {
	case reflect.Array:
		if yml == nil {
			return
		}

		for i := 0; i < rv.Cap(); i++ {
			fv := reflect.Indirect(reflect.New(rv.Index(i).Type()))
			unmarshal(yml.([]interface{})[i], fv)

			rv.Index(i).Set(fv)
		}
	case reflect.Slice:
		if yml == nil {
			return
		}

		yarr := yml.([]interface{})
		for i := 0; i < len(yarr); i++ {
			sv := reflect.MakeSlice(rv.Type(), 1, 1)
			fv := reflect.Indirect(reflect.New(sv.Index(0).Type()))
			unmarshal(yarr[i], fv)

			rv.Set(reflect.Append(rv, fv))
		}
	case reflect.Map:
		if yml == nil {
			return
		}

		rv.Set(reflect.MakeMap(rv.Type()))

		ymap := yml.(map[interface{}]interface{})
		for k, v := range ymap {
			fv := reflect.Indirect(reflect.New(rv.Type().Elem()))
			unmarshal(v, fv)

			rv.SetMapIndex(reflect.ValueOf(k), fv)
		}
	case reflect.Struct:
		if yml == nil {
			return
		}

		for i := 0; i < rv.NumField(); i++ {
			ft := rt.Field(i).Tag.Get("conf")
			fv := reflect.Indirect(reflect.New(rv.Field(i).Type()))
			unmarshal(yml.(map[interface{}]interface{})[ft], fv)

			rv.Field(i).Set(fv)
		}
	case reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.Ptr, reflect.UnsafePointer:
		panic(ErrInvalidDataType)
	case reflect.Int8:
		if yml != nil {
			rv.Set(reflect.ValueOf(int8(yml.(int))))
		} else {
			rv.Set(reflect.ValueOf(int8(0)))
		}
	case reflect.Int16:
		if yml != nil {
			rv.Set(reflect.ValueOf(int16(yml.(int))))
		} else {
			rv.Set(reflect.ValueOf(int16(0)))
		}
	case reflect.Int32:
		if yml != nil {
			rv.Set(reflect.ValueOf(int32(yml.(int))))
		} else {
			rv.Set(reflect.ValueOf(int32(0)))
		}
	case reflect.Int64:
		if yml != nil {
			rv.Set(reflect.ValueOf(int64(yml.(int))))
		} else {
			rv.Set(reflect.ValueOf(int64(0)))
		}
	case reflect.Uint:
		if yml != nil {
			rv.Set(reflect.ValueOf(uint(yml.(int))))
		} else {
			rv.Set(reflect.ValueOf(uint(0)))
		}
	case reflect.Uint8:
		if yml != nil {
			rv.Set(reflect.ValueOf(uint8(yml.(int))))
		} else {
			rv.Set(reflect.ValueOf(uint8(0)))
		}
	case reflect.Uint16:
		if yml != nil {
			rv.Set(reflect.ValueOf(uint16(yml.(int))))
		} else {
			rv.Set(reflect.ValueOf(uint16(0)))
		}
	case reflect.Uint32:
		if yml != nil {
			rv.Set(reflect.ValueOf(uint32(yml.(int))))
		} else {
			rv.Set(reflect.ValueOf(uint32(0)))
		}
	case reflect.Uint64:
		if yml != nil {
			rv.Set(reflect.ValueOf(uint64(yml.(int))))
		} else {
			rv.Set(reflect.ValueOf(uint64(0)))
		}
	case reflect.Uintptr:
		if yml != nil {
			rv.Set(reflect.ValueOf(uintptr(yml.(int))))
		} else {
			rv.Set(reflect.ValueOf(uintptr(0)))
		}
	case reflect.Float32:
		if yml != nil {
			rv.Set(reflect.ValueOf(float32(yml.(float64))))
		} else {
			rv.Set(reflect.ValueOf(float32(0)))
		}
	default:
		if yml != nil {
			rv.Set(reflect.ValueOf(yml))
		}
	}
}
