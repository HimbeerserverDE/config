package config_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/HimbeerserverDE/config"
)

type TestConf struct {
	TFV bool              `conf:"tfv"`
	Num int               `conf:"num"`
	Str string            `conf:"str"`
	Arr []string          `conf:"arr"`
	Map map[string]string `conf:"map"`
	T2  TestConf2         `conf:"t2"`
}

type TestConf2 struct {
	Str string `conf:"str"`
}

var expectedYML = `arr:
- foo
- bar
map:
  foo: bar
num: 1337
str: foo
t2:
  str: bar
tfv: true
`

var expectedConf = TestConf{
	TFV: true,
	Num: 1337,
	Str: "foo",
	Arr: []string{"foo", "bar"},
	Map: map[string]string{"foo": "bar"},
	T2: TestConf2{
		Str: "bar",
	},
}

func TestMarshal(t *testing.T) {
	buf := &strings.Builder{}
	if err := config.Marshal(buf, expectedConf); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if yml := buf.String(); yml != expectedYML {
		t.Log("incorrect YAML result")
		t.Log(yml)
		t.FailNow()
	}
}

func TestUnmarshal(t *testing.T) {
	c := TestConf{}
	buf := strings.NewReader(expectedYML)
	if err := config.Unmarshal(buf, &c); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(c, expectedConf) {
		t.Log("incorrect configuration result")
		t.Log(c)
		t.FailNow()
	}
}
