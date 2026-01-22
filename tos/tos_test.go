package tos

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestTos(t *testing.T) {
	ctx := context.Background()
	key, err := Put(ctx, "hello")
	if err != nil {
		t.Fatal(err)
	}
	b, err := Get(ctx, key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}

func TestTos2(t *testing.T) {
	ctx := context.Background()
	b, err := os.ReadFile("/tmp/test.zip")
	if err != nil {
		t.Fatal(err)
	}
	key, err := PutBytes(ctx, b)
	if err != nil {
		t.Fatal(err)
	}
	b, err = Get(ctx, key)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
}
