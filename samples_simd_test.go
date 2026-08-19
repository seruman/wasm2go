//go:build !generator && goexperiment.simd && (amd64 || arm64 || wasm)

package main

import (
	"testing"

	simd_test "github.com/ncruces/wasm2go/testdata/simd"
	simd_type_test "github.com/ncruces/wasm2go/testdata/simd_type"
)

func Test_simd(t *testing.T) {
	m := simd_test.New()

	tests := []struct {
		name string
		got  int32
		want int32
	}{
		{"block", m.Xblock(), 0xffff},
		{"global", m.Xglobal(), 0x00ff},
		{"mask", m.Xmask(), 0xffff},
		{"load_store", m.Xload_store(), 15},
		{"load_lane", m.Xload_lane(), 171},
		{"store_lane", m.Xstore_lane(), 15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %d, want %d", tt.got, tt.want)
			}
		})
	}
}

func Test_simdShift(t *testing.T) {
	m := simd_test.New()
	if got, want := m.Xshr_u(), int64(0x0fffffffffffffff); got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func Test_simdType(t *testing.T) {
	if got := simd_type_test.New(); got == nil {
		t.Error("New returned nil")
	}
}
