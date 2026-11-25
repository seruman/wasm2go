//go:build go1.27 && goexperiment.simd && (amd64 || arm64 || wasm)

package shuffle

import (
	"simd/archsimd"
	"testing"
)

func TestShuffle(t *testing.T) {
	m := New()
	a := [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	b := [16]byte{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	x, y := archsimd.LoadUint8x16Array(&a), archsimd.LoadUint8x16Array(&b)
	for _, tt := range []struct {
		name string
		f    func(v128, v128) v128
		want [16]byte
	}{
		{"x", m.Xx, a},
		{"y", m.Xy, b},
		{"repeat", m.Xrepeat, [16]byte{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7}},
		{"reverse_y", m.Xreverse_y, [16]byte{31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16}},
		{"mixed", m.Xmixed, [16]byte{31, 0, 30, 1, 29, 2, 28, 3, 27, 4, 26, 5, 25, 6, 24, 7}},
		{"word", m.Xword, [16]byte{14, 15, 30, 31, 2, 3, 10, 11, 4, 5, 16, 17, 4, 5, 8, 9}},
		{"dword", m.Xdword, [16]byte{16, 17, 18, 19, 8, 9, 10, 11, 28, 29, 30, 31, 0, 1, 2, 3}},
		{"half", m.Xhalf, [16]byte{24, 25, 26, 27, 28, 29, 30, 31, 0, 1, 2, 3, 4, 5, 6, 7}},
	} {
		var got [16]byte
		tt.f(x, y).StoreArray(&got)
		if got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

func BenchmarkShuffleMixed(b *testing.B) {
	m := New()
	x, y := archsimd.BroadcastUint8x16(42), archsimd.BroadcastUint8x16(15)
	for b.Loop() {
		x = m.Xmixed(x, y)
	}
}
