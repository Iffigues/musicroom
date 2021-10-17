package config

import (
    "testing"
)

func TestHelloEmpty(t *testing.T) {
	c := NewConf()
	go serve()
	err := c.NewConfType("test", false)
	if err != nil {
		t.Fatalf("err one")
	}
	err = c.NewConfType("test", false)
	if err == nil {
		t.Fatalf("error 2")
	}
	err = c.NewConfType("test", true)
	if err != nil {
		t.Fatal("error 3")
	}

	err = c.NewConfType("rd", true)
	if err != nil {
		t.Fatal("error 4")
	}
}
