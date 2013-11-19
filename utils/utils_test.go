package utils

import (
	"testing"
)

type Foo struct {
	Bar string
}

func (*Foo) Construct(args ...interface{}) *Foo {
	return &Foo{Bar: args[0].([]interface{})[0].(string)}
}

func TestConstruct(t *testing.T) {
	x := Construct(new(Foo), "xxx").(*Foo)
	if x.Bar != "xxx" {
		t.Errorf("x.Bar = '%v' want 'xxx'", x.Bar)
	}
}
