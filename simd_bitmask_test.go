package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
	"testing"

	"github.com/ncruces/wasm2go/internal/offset"
)

func TestBitmaskLowering(t *testing.T) {
	const (
		mask  = "\x20\x00\xfd\x64"
		zero  = "\x41\x00"
		one   = "\x41\x01"
		cond  = "i8x16_bitmask_nonzero(v0)"
		exact = "i8x16_bitmask(v0)"
	)
	for _, tt := range []struct{ name, code, want string }{
		{"br_if", one + mask + "\x0d\x00\x1a" + zero, "if " + cond},
		{"eqz", mask + "\x45", "if !" + cond},
		{"eq_zero", mask + zero + "\x46", "if !" + cond},
		{"zero_ne", zero + mask + "\x47", "if " + cond},
		{"return", mask, "return " + exact},
		{"eq_one", mask + one + "\x46", exact + " == i32(1)"},
		{"tee", mask + "\x22\x01\x04\x7f\x20\x01\x68\x05" + one + "\x0b", "v1 = " + exact},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tr := translator{
				in:      offset.NewReader(strings.NewReader("\x01\x01\x7f" + tt.code + "\x0b")),
				helpers: set[string]{}, packages: set[string]{},
			}
			typ := funcType{params: "\x7b", results: "\x7f"}
			fn := funcCompiler{typ: typ, decl: &ast.FuncDecl{Name: newID("Xtest"), Type: typ.toAST(true)}}
			if err := tr.readCodeForFunction(&fn); err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			if err := format.Node(&out, token.NewFileSet(), fn.decl.Body); err != nil {
				t.Fatal(err)
			}
			want := tt.want
			if *noopt {
				want = "i8x16_bitmask("
			}
			if !strings.Contains(out.String(), want) {
				t.Errorf("wanted %q in:\n%s", want, &out)
			}
		})
	}
}
