//go:build goexperiment.simd && (amd64 || arm64 || wasm)

package helpers

import (
	"encoding/binary"
	"testing"
)

func Test_v128_bitselect(t *testing.T) {
	x := v128_from_bytes([16]byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	})
	y := v128_from_bytes([16]byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	})
	mask := v128_from_bytes([16]byte{
		0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00,
		0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00,
	})
	want := [16]byte{
		0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00,
		0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff,
	}

	if got := v128_bytes(v128_bitselect(x, y, mask)); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_i8x16_shuffle(t *testing.T) {
	x := v128_from_bytes([16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	y := v128_from_bytes([16]byte{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31})
	want := [16]byte{31, 0, 30, 1, 29, 2, 28, 3, 27, 4, 26, 5, 25, 6, 24, 7}

	got := i8x16_shuffle(x, y, 31, 0, 30, 1, 29, 2, 28, 3, 27, 4, 26, 5, 25, 6, 24, 7)
	if got := v128_bytes(got); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_i8x16_narrow_i16x8_u(t *testing.T) {
	var xb, yb [16]byte
	x := [...]int16{-1, 0, 1, 254, 255, 256, 32767, -32768}
	y := [...]int16{42, 512, -42, 127, 128, 129, 300, -300}
	for i := range x {
		binary.LittleEndian.PutUint16(xb[i*2:], uint16(x[i]))
		binary.LittleEndian.PutUint16(yb[i*2:], uint16(y[i]))
	}
	want := [16]byte{0, 0, 1, 254, 255, 255, 255, 0, 42, 255, 0, 127, 128, 129, 255, 0}

	got := i8x16_narrow_i16x8_u(v128_from_bytes(xb), v128_from_bytes(yb))
	if got := v128_bytes(got); got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_i16x8_extend_high_i8x16_s(t *testing.T) {
	x := v128_from_bytes([16]byte{0, 1, 2, 3, 4, 5, 6, 7, 0x80, 0x81, 0xfe, 0xff, 0, 1, 126, 127})
	want := [8]int16{-128, -127, -2, -1, 0, 1, 126, 127}

	b := v128_bytes(i16x8_extend_high_i8x16_s(x))
	var got [8]int16
	for i := range got {
		got[i] = int16(binary.LittleEndian.Uint16(b[i*2:]))
	}
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func Test_v128_memory(t *testing.T) {
	mem := make([]byte, 32)
	binary.LittleEndian.PutUint32(mem[3:], 0x89abcdef)
	binary.LittleEndian.PutUint32(mem[7:], 0x01234567)
	wantLoad := [16]byte{
		0xef, 0xcd, 0xab, 0x89, 0, 0, 0, 0,
		0x67, 0x45, 0x23, 0x01, 0, 0, 0, 0,
	}
	if got := v128_bytes(v128_load32x2_u(mem, uint32(3))); got != wantLoad {
		t.Errorf("got %v, want %v", got, wantLoad)
	}

	clear(mem)
	x := v128_from_bytes([16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	v128_store64_lane(mem, int64(16), x, 1)
	wantStore := [32]byte{
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		8, 9, 10, 11, 12, 13, 14, 15,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	if got := [32]byte(mem); got != wantStore {
		t.Errorf("got %v, want %v", got, wantStore)
	}
}

func Test_i32x4(t *testing.T) {
	var xb, yb [16]byte
	x := [...]int32{-1, 0, 1, 0x40000000}
	y := [...]int32{1, 0, -1, 0x40000000}
	for i := range x {
		binary.LittleEndian.PutUint32(xb[i*4:], uint32(x[i]))
		binary.LittleEndian.PutUint32(yb[i*4:], uint32(y[i]))
	}
	xv := v128_from_bytes(xb)
	yv := v128_from_bytes(yb)

	wantCompare := [16]byte{
		0xff, 0xff, 0xff, 0xff, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	if got := v128_bytes(i32x4_gt_u(xv, yv)); got != wantCompare {
		t.Errorf("got %v, want %v", got, wantCompare)
	}

	wantAdd := [4]int32{0, 0, 0, -0x80000000}
	b := v128_bytes(i32x4_add(xv, yv))
	var got [4]int32
	for i := range got {
		got[i] = int32(binary.LittleEndian.Uint32(b[i*4:]))
	}
	if got != wantAdd {
		t.Errorf("got %v, want %v", got, wantAdd)
	}
}

func Test_f64x2(t *testing.T) {
	x := f64x2_replace_lane(f64x2_splat(1.5), 1, -2)
	y := f64x2_replace_lane(f64x2_splat(2), 1, 4)

	wantAdd := [2]float64{3.5, 2}
	add := f64x2_add(x, y)
	gotAdd := [2]float64{f64x2_extract_lane(add, 0), f64x2_extract_lane(add, 1)}
	if gotAdd != wantAdd {
		t.Errorf("got %v, want %v", gotAdd, wantAdd)
	}

	wantMul := [2]float64{3, -8}
	mul := f64x2_mul(x, y)
	gotMul := [2]float64{f64x2_extract_lane(mul, 0), f64x2_extract_lane(mul, 1)}
	if gotMul != wantMul {
		t.Errorf("got %v, want %v", gotMul, wantMul)
	}
}
