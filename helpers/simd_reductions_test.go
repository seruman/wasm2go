//go:build go1.27 && goexperiment.simd && (amd64 || arm64 || wasm)

package helpers

import (
	"encoding/binary"
	"math/rand/v2"
	"testing"
)

func TestSIMDReductions(t *testing.T) {
	for _, tc := range []struct {
		name      string
		width     int
		mask, all func(v128) int32
	}{
		{"i8x16", 1, i8x16_bitmask, i8x16_all_true},
		{"i16x8", 2, i16x8_bitmask, i16x8_all_true},
		{"i32x4", 4, i32x4_bitmask, i32x4_all_true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			check := func(b [16]byte) {
				t.Helper()
				var mask, any int32
				all := int32(1)
				for i := 0; i < len(b); i += tc.width {
					var lane byte
					for _, v := range b[i : i+tc.width] {
						lane |= v
					}
					if lane == 0 {
						all = 0
					} else {
						any = 1
					}
					mask |= int32(b[i+tc.width-1]>>7) << (i / tc.width)
				}
				x := v128_from_bytes(b)
				if tc.width == 1 && i8x16_bitmask_nonzero(x) != (mask != 0) {
					t.Fatalf("%x: incorrect bitmask condition for mask %#x", b, mask)
				}
				if m, a, v := tc.mask(x), tc.all(x), v128_any_true(x); m != mask || a != all || v != any {
					t.Fatalf("%x: mask/all/any = %#x/%d/%d, want %#x/%d/%d", b, m, a, v, mask, all, any)
				}
			}
			check([16]byte{})
			for i := range 128 {
				var b [16]byte
				b[i/8] = 1 << (i % 8)
				check(b)
				for j := range b {
					b[j] = ^b[j]
				}
				check(b)
			}
			r := rand.New(rand.NewPCG(17, 31))
			for range 256 {
				var b [16]byte
				binary.LittleEndian.PutUint64(b[:], r.Uint64())
				binary.LittleEndian.PutUint64(b[8:], r.Uint64())
				check(b)
				lane := r.IntN(len(b) / tc.width)
				clear(b[lane*tc.width : (lane+1)*tc.width])
				check(b)
			}
		})
	}
}

func TestSIMDNarrow(t *testing.T) {
	for _, tc := range []struct {
		name   string
		width  int
		lo, hi int32
		narrow func(v128, v128) v128
	}{
		{"i16_to_i8", 1, -128, 127, i8x16_narrow_i16x8_s},
		{"i16_to_u8", 1, 0, 255, i8x16_narrow_i16x8_u},
		{"i32_to_i16", 2, -32768, 32767, i16x8_narrow_i32x4_s},
		{"i32_to_u16", 2, 0, 65535, i16x8_narrow_i32x4_u},
	} {
		t.Run(tc.name, func(t *testing.T) {
			values := [...]int32{tc.lo - 1, tc.lo, tc.lo + 1, tc.hi - 1, tc.hi, tc.hi + 1, -1, 0, 1, 37, -67, 0x1234, -0x2356, -32768, 32767, -1 << 31, 1<<31 - 1}
			for start := range values {
				var src [32]byte
				var want [16]byte
				for i := 0; i < len(want)/tc.width; i++ {
					v := values[(start+i)%len(values)]
					if tc.width == 1 {
						v = int32(int16(v))
						binary.LittleEndian.PutUint16(src[2*i:], uint16(v))
						want[i] = byte(min(max(v, tc.lo), tc.hi))
					} else {
						binary.LittleEndian.PutUint32(src[4*i:], uint32(v))
						binary.LittleEndian.PutUint16(want[2*i:], uint16(min(max(v, tc.lo), tc.hi)))
					}
				}
				x, y := v128_from_bytes([16]byte(src[:16])), v128_from_bytes([16]byte(src[16:]))
				if got := v128_bytes(tc.narrow(x, y)); got != want {
					t.Fatalf("%x: got %x, want %x", src, got, want)
				}
			}
		})
	}
}

func BenchmarkSIMDReduction(b *testing.B) {
	inputs := [4]v128{
		v128_const(0, 0), v128_const(0, 1<<63),
		v128_const(0x12ff56ff80feff2d, 0xfc89141f4192a3ef), v128_const(0x8080808080808080, 0x0101010101010101),
	}
	for _, tc := range []struct {
		name string
		op   func(v128) int32
	}{
		{"bitmask8", i8x16_bitmask}, {"bitmask16", i16x8_bitmask}, {"bitmask32", i32x4_bitmask},
		{"any", v128_any_true}, {"all8", i8x16_all_true}, {"all16", i16x8_all_true}, {"all32", i32x4_all_true},
	} {
		b.Run(tc.name, func(b *testing.B) {
			i := 0
			for b.Loop() {
				tc.op(inputs[i&3])
				i++
			}
		})
	}
}

func BenchmarkSIMDNarrow(b *testing.B) {
	x, y := v128_const(0x00ff567880feff2d, 0xfc89141f4192a3ef), v128_const(0x0118807f00000001, 0xff09003e001a0599)
	for _, tc := range []struct {
		name string
		op   func(v128, v128) v128
	}{
		{"i16_to_i8", i8x16_narrow_i16x8_s}, {"i16_to_u8", i8x16_narrow_i16x8_u},
		{"i32_to_i16", i16x8_narrow_i32x4_s}, {"i32_to_u16", i16x8_narrow_i32x4_u},
	} {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				tc.op(x, y)
			}
		})
	}
}
