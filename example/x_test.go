package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestFoo(t *testing.T) {
	type Demo struct {
		Value  []string `json:"value,omitempty"`
		Value2 []string `json:"value2"`
	}
	demo := Demo{
		Value:  make([]string, 0, 10),
		Value2: make([]string, 0, 10),
	}
	b, _ := json.Marshal(demo)
	fmt.Println(string(b))
}
