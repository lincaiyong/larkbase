package larkbase

import (
	"context"
	"github.com/lincaiyong/larkbase/tos"
)

type TosPutFn func(ctx context.Context, bs []byte) (string, error)
type TosGetFn func(ctx context.Context, key string) ([]byte, error)

var tosPutFn TosPutFn
var tosGetFn TosGetFn

func init() {
	tosPutFn = tos.Put
	tosGetFn = tos.Get
}

func CustomizeTos(tosPut TosPutFn, tosGet TosGetFn) {
	tosPutFn = tosPut
	tosGetFn = tosGet
}
