package config_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/HimbeerserverDE/config"
)

type TestConf struct {
	Num int    `conf:"num"`
	Str string `conf:"str"`
}

var expectedYML = `Num: 99116
Str: ct
`

var expectedConf = TestConf{
	Num: 99116,
	Str: "ct",
}

func TestMarshal(t *testing.T) {
	buf := &bytes.Buffer{}
	if err := config.Marshal(buf, expectedConf); err != nil {
		t.Log(err)
		t.FailNow()
	}

	if yml := string(buf.Bytes()); yml != expectedYML {
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

	if c.Num != expectedConf.Num || c.Str != expectedConf.Str {
		t.Log("incorrect configuration result")
		t.Log(c)
		t.FailNow()
	}
}
