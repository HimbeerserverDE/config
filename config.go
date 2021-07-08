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
	rv := reflect.ValueOf(v)
	rt := reflect.TypeOf(v)

	ymap := make(map[interface{}]interface{})
	for i := 0; i < rv.NumField(); i++ {
		ft := rt.Field(i).Tag.Get("conf")
		fv := rv.Field(i).Interface()

		ymap[ft] = fv
	}

	out, err := yaml.Marshal(ymap)
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

	rv := reflect.Indirect(reflect.ValueOf(v))
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		ft := rt.Field(i).Tag.Get("conf")

		if !rv.Field(i).CanSet() {
			return ErrFieldUnsettable
		}

		if _, ok := ymap[ft]; ok {
			rv.Field(i).Set(reflect.ValueOf(ymap[ft]))
		}
	}

	return nil
}
