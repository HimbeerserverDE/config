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

// Marshal writes a configuration to the specified Writer
func Marshal(w io.Writer, v interface{}) error {
	yml := marshal(v).(map[interface{}]interface{})

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
	default:
		rv.Set(reflect.ValueOf(yml))
	}
}
