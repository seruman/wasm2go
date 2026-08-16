//go:build goexperiment.simd && (amd64 || arm64 || wasm)

package helpers

import (
	"encoding/binary"
	"runtime"
	"simd/archsimd"
)

// SIMD helpers.

type v128 = archsimd.Uint8x16

func v128_const(lo, hi uint64) v128 {
	var b [16]byte
	binary.LittleEndian.PutUint64(b[:8], lo)
	binary.LittleEndian.PutUint64(b[8:], hi)
	return archsimd.LoadUint8x16Array(&b)
}

func v128_check_cpu() {
	if runtime.GOARCH == "amd64" && !archsimd.X86.AVX2() {
		panic("wasm2go: SIMD requires AVX2 on amd64")
	}
}

func v128_bytes(x v128) [16]byte {
	var b [16]byte
	x.StoreArray(&b)
	return b
}

func v128_from_bytes(b [16]byte) v128 {
	return archsimd.LoadUint8x16Array(&b)
}

func asI8(x v128) archsimd.Int8x16 {
	return x.BitsToInt8()
}

func asU16(x v128) archsimd.Uint16x8 {
	return x.ReshapeToUint16s()
}

func asI16(x v128) archsimd.Int16x8 {
	return x.ReshapeToUint16s().BitsToInt16()
}

func asU32(x v128) archsimd.Uint32x4 {
	return x.ReshapeToUint32s()
}

func asI32(x v128) archsimd.Int32x4 {
	return x.ReshapeToUint32s().BitsToInt32()
}

func asU64(x v128) archsimd.Uint64x2 {
	return x.ReshapeToUint64s()
}

func asI64(x v128) archsimd.Int64x2 {
	return x.ReshapeToUint64s().BitsToInt64()
}

func asF32(x v128) archsimd.Float32x4 {
	return x.ReshapeToUint32s().BitsToFloat32()
}

func asF64(x v128) archsimd.Float64x2 {
	return x.ReshapeToUint64s().BitsToFloat64()
}

func fromI8(x archsimd.Int8x16) v128 {
	return x.ToBits()
}

func fromU16(x archsimd.Uint16x8) v128 {
	return x.ReshapeToUint8s()
}

func fromI16(x archsimd.Int16x8) v128 {
	return x.ToBits().ReshapeToUint8s()
}

func fromU32(x archsimd.Uint32x4) v128 {
	return x.ReshapeToUint8s()
}

func fromI32(x archsimd.Int32x4) v128 {
	return x.ToBits().ReshapeToUint8s()
}

func fromU64(x archsimd.Uint64x2) v128 {
	return x.ReshapeToUint8s()
}

func fromI64(x archsimd.Int64x2) v128 {
	return x.ToBits().ReshapeToUint8s()
}

func fromF32(x archsimd.Float32x4) v128 {
	return x.ToBits().ReshapeToUint8s()
}

func fromF64(x archsimd.Float64x2) v128 {
	return x.ToBits().ReshapeToUint8s()
}

func fromM8(x archsimd.Mask8x16) v128 {
	return x.ToInt8x16().ToBits()
}

func fromM16(x archsimd.Mask16x8) v128 {
	return x.ToInt16x8().ToBits().ReshapeToUint8s()
}

func fromM32(x archsimd.Mask32x4) v128 {
	return x.ToInt32x4().ToBits().ReshapeToUint8s()
}

func fromM64(x archsimd.Mask64x2) v128 {
	return x.ToInt64x2().ToBits().ReshapeToUint8s()
}

func v128_load[T uint32 | int64](mem []byte, addr T) v128 {
	return archsimd.LoadUint8x16(mem[addr:])
}

func v128_store[T uint32 | int64](mem []byte, addr T, x v128) {
	x.Store(mem[addr:])
}

func v128_load32x2_s[T uint32 | int64](mem []byte, addr T) v128 {
	b := mem[addr:]
	v := [2]int64{
		int64(int32(binary.LittleEndian.Uint32(b))),
		int64(int32(binary.LittleEndian.Uint32(b[4:]))),
	}
	return fromI64(archsimd.LoadInt64x2Array(&v))
}

func v128_load32x2_u[T uint32 | int64](mem []byte, addr T) v128 {
	b := mem[addr:]
	v := [2]uint64{
		uint64(binary.LittleEndian.Uint32(b)),
		uint64(binary.LittleEndian.Uint32(b[4:])),
	}
	return fromU64(archsimd.LoadUint64x2Array(&v))
}

func v128_load8_splat[T uint32 | int64](mem []byte, addr T) v128 {
	return archsimd.BroadcastUint8x16(mem[addr])
}

func v128_load16_splat[T uint32 | int64](mem []byte, addr T) v128 {
	return fromU16(archsimd.BroadcastUint16x8(binary.LittleEndian.Uint16(mem[addr:])))
}

func v128_load32_splat[T uint32 | int64](mem []byte, addr T) v128 {
	return fromU32(archsimd.BroadcastUint32x4(binary.LittleEndian.Uint32(mem[addr:])))
}

func v128_load64_splat[T uint32 | int64](mem []byte, addr T) v128 {
	return fromU64(archsimd.BroadcastUint64x2(binary.LittleEndian.Uint64(mem[addr:])))
}

func v128_load32_zero[T uint32 | int64](mem []byte, addr T) v128 {
	return fromU32(archsimd.BroadcastUint32x4(0).SetElem(0, binary.LittleEndian.Uint32(mem[addr:])))
}

func v128_load64_zero[T uint32 | int64](mem []byte, addr T) v128 {
	return fromU64(archsimd.BroadcastUint64x2(0).SetElem(0, binary.LittleEndian.Uint64(mem[addr:])))
}

func v128_load8_lane[T uint32 | int64](mem []byte, addr T, x v128, lane uint8) v128 {
	return x.SetElem(lane, mem[addr])
}

func v128_load16_lane[T uint32 | int64](mem []byte, addr T, x v128, lane uint8) v128 {
	return fromU16(asU16(x).SetElem(lane, binary.LittleEndian.Uint16(mem[addr:])))
}

func v128_load32_lane[T uint32 | int64](mem []byte, addr T, x v128, lane uint8) v128 {
	return fromU32(asU32(x).SetElem(lane, binary.LittleEndian.Uint32(mem[addr:])))
}

func v128_load64_lane[T uint32 | int64](mem []byte, addr T, x v128, lane uint8) v128 {
	return fromU64(asU64(x).SetElem(lane, binary.LittleEndian.Uint64(mem[addr:])))
}

func v128_store8_lane[T uint32 | int64](mem []byte, addr T, x v128, lane uint8) {
	mem[addr] = x.GetElem(lane)
}

func v128_store16_lane[T uint32 | int64](mem []byte, addr T, x v128, lane uint8) {
	binary.LittleEndian.PutUint16(mem[addr:], asU16(x).GetElem(lane))
}

func v128_store32_lane[T uint32 | int64](mem []byte, addr T, x v128, lane uint8) {
	binary.LittleEndian.PutUint32(mem[addr:], asU32(x).GetElem(lane))
}

func v128_store64_lane[T uint32 | int64](mem []byte, addr T, x v128, lane uint8) {
	binary.LittleEndian.PutUint64(mem[addr:], asU64(x).GetElem(lane))
}

func v128_not(x v128) v128 {
	return x.Not()
}

func v128_and(x, y v128) v128 {
	return x.And(y)
}

func v128_andnot(x, y v128) v128 {
	return x.AndNot(y)
}

func v128_or(x, y v128) v128 {
	return x.Or(y)
}

func v128_xor(x, y v128) v128 {
	return x.Xor(y)
}

func v128_bitselect(x, y, c v128) v128 {
	return c.And(x).Or(c.Not().And(y))
}

func v128_any_true(x v128) int32 {
	b := v128_bytes(x)
	for _, v := range b {
		if v != 0 {
			return 1
		}
	}
	return 0
}

func i8x16_shuffle(x, y v128, l0, l1, l2, l3, l4, l5, l6, l7, l8, l9, l10, l11, l12, l13, l14, l15 uint8) v128 {
	xb := v128_bytes(x)
	yb := v128_bytes(y)
	lanes := [...]uint8{l0, l1, l2, l3, l4, l5, l6, l7, l8, l9, l10, l11, l12, l13, l14, l15}
	var b [16]byte
	for i, lane := range lanes {
		if lane < 16 {
			b[i] = xb[lane]
		} else {
			b[i] = yb[lane-16]
		}
	}
	return v128_from_bytes(b)
}

func i8x16_splat(x int32) v128 {
	return archsimd.BroadcastUint8x16(byte(x))
}

func i16x8_splat(x int32) v128 {
	return fromU16(archsimd.BroadcastUint16x8(uint16(x)))
}

func i32x4_splat(x int32) v128 {
	return fromI32(archsimd.BroadcastInt32x4(x))
}

func i64x2_splat(x int64) v128 {
	return fromI64(archsimd.BroadcastInt64x2(x))
}

func f32x4_splat(x float32) v128 {
	return fromF32(archsimd.BroadcastFloat32x4(x))
}

func f64x2_splat(x float64) v128 {
	return fromF64(archsimd.BroadcastFloat64x2(x))
}

func i8x16_extract_lane_s(x v128, lane uint8) int32 {
	return int32(int8(x.GetElem(lane)))
}

func i8x16_extract_lane_u(x v128, lane uint8) int32 {
	return int32(x.GetElem(lane))
}

func i8x16_replace_lane(x v128, lane uint8, val int32) v128 {
	return x.SetElem(lane, byte(val))
}

func i16x8_extract_lane_s(x v128, lane uint8) int32 {
	return int32(asI16(x).GetElem(lane))
}

func i16x8_extract_lane_u(x v128, lane uint8) int32 {
	return int32(asU16(x).GetElem(lane))
}

func i16x8_replace_lane(x v128, lane uint8, val int32) v128 {
	return fromU16(asU16(x).SetElem(lane, uint16(val)))
}

func i32x4_extract_lane(x v128, lane uint8) int32 {
	return asI32(x).GetElem(lane)
}

func i32x4_replace_lane(x v128, lane uint8, val int32) v128 {
	return fromI32(asI32(x).SetElem(lane, val))
}

func i64x2_extract_lane(x v128, lane uint8) int64 {
	return asI64(x).GetElem(lane)
}

func i64x2_replace_lane(x v128, lane uint8, val int64) v128 {
	return fromI64(asI64(x).SetElem(lane, val))
}

func f32x4_extract_lane(x v128, lane uint8) float32 {
	return asF32(x).GetElem(lane)
}

func f32x4_replace_lane(x v128, lane uint8, val float32) v128 {
	return fromF32(asF32(x).SetElem(lane, val))
}

func f64x2_extract_lane(x v128, lane uint8) float64 {
	return asF64(x).GetElem(lane)
}

func f64x2_replace_lane(x v128, lane uint8, val float64) v128 {
	return fromF64(asF64(x).SetElem(lane, val))
}

func i8x16_eq(x, y v128) v128 {
	return fromM8(x.Equal(y))
}

func i8x16_ne(x, y v128) v128 {
	return fromM8(x.NotEqual(y))
}

func i8x16_lt_s(x, y v128) v128 {
	return fromM8(asI8(x).Less(asI8(y)))
}

func i8x16_lt_u(x, y v128) v128 {
	return fromM8(x.Less(y))
}

func i8x16_gt_s(x, y v128) v128 {
	return fromM8(asI8(x).Greater(asI8(y)))
}

func i8x16_gt_u(x, y v128) v128 {
	return fromM8(x.Greater(y))
}

func i8x16_le_s(x, y v128) v128 {
	return fromM8(asI8(x).LessEqual(asI8(y)))
}

func i8x16_le_u(x, y v128) v128 {
	return fromM8(x.LessEqual(y))
}

func i8x16_ge_s(x, y v128) v128 {
	return fromM8(asI8(x).GreaterEqual(asI8(y)))
}

func i8x16_ge_u(x, y v128) v128 {
	return fromM8(x.GreaterEqual(y))
}

func i16x8_eq(x, y v128) v128 {
	return fromM16(asU16(x).Equal(asU16(y)))
}

func i16x8_ne(x, y v128) v128 {
	return fromM16(asU16(x).NotEqual(asU16(y)))
}

func i16x8_lt_s(x, y v128) v128 {
	return fromM16(asI16(x).Less(asI16(y)))
}

func i16x8_lt_u(x, y v128) v128 {
	return fromM16(asU16(x).Less(asU16(y)))
}

func i16x8_gt_s(x, y v128) v128 {
	return fromM16(asI16(x).Greater(asI16(y)))
}

func i16x8_gt_u(x, y v128) v128 {
	return fromM16(asU16(x).Greater(asU16(y)))
}

func i16x8_le_s(x, y v128) v128 {
	return fromM16(asI16(x).LessEqual(asI16(y)))
}

func i16x8_le_u(x, y v128) v128 {
	return fromM16(asU16(x).LessEqual(asU16(y)))
}

func i16x8_ge_s(x, y v128) v128 {
	return fromM16(asI16(x).GreaterEqual(asI16(y)))
}

func i16x8_ge_u(x, y v128) v128 {
	return fromM16(asU16(x).GreaterEqual(asU16(y)))
}

func i32x4_eq(x, y v128) v128 {
	return fromM32(asU32(x).Equal(asU32(y)))
}

func i32x4_ne(x, y v128) v128 {
	return fromM32(asU32(x).NotEqual(asU32(y)))
}

func i32x4_lt_s(x, y v128) v128 {
	return fromM32(asI32(x).Less(asI32(y)))
}

func i32x4_lt_u(x, y v128) v128 {
	return fromM32(asU32(x).Less(asU32(y)))
}

func i32x4_gt_s(x, y v128) v128 {
	return fromM32(asI32(x).Greater(asI32(y)))
}

func i32x4_gt_u(x, y v128) v128 {
	return fromM32(asU32(x).Greater(asU32(y)))
}

func i32x4_le_s(x, y v128) v128 {
	return fromM32(asI32(x).LessEqual(asI32(y)))
}

func i32x4_le_u(x, y v128) v128 {
	return fromM32(asU32(x).LessEqual(asU32(y)))
}

func i32x4_ge_s(x, y v128) v128 {
	return fromM32(asI32(x).GreaterEqual(asI32(y)))
}

func i32x4_ge_u(x, y v128) v128 {
	return fromM32(asU32(x).GreaterEqual(asU32(y)))
}

func i64x2_eq(x, y v128) v128 {
	return fromM64(asU64(x).Equal(asU64(y)))
}

func i64x2_ne(x, y v128) v128 {
	return fromM64(asU64(x).NotEqual(asU64(y)))
}

func i8x16_all_true(x v128) int32 {
	b := v128_bytes(x)
	for _, v := range b {
		if v == 0 {
			return 0
		}
	}
	return 1
}

func i8x16_bitmask(x v128) int32 {
	b := v128_bytes(x)
	var r int32
	for i, v := range b {
		r |= int32(v>>7) << i
	}
	return r
}

func i8x16_min_u(x, y v128) v128 {
	return x.Min(y)
}

func i16x8_all_true(x v128) int32 {
	b := v128_bytes(x)
	for i := range 8 {
		if binary.LittleEndian.Uint16(b[i*2:]) == 0 {
			return 0
		}
	}
	return 1
}

func i16x8_bitmask(x v128) int32 {
	b := v128_bytes(x)
	var r int32
	for i := range 8 {
		r |= int32(binary.LittleEndian.Uint16(b[i*2:])>>15) << i
	}
	return r
}

func i16x8_min_u(x, y v128) v128 {
	return fromU16(asU16(x).Min(asU16(y)))
}

func i32x4_all_true(x v128) int32 {
	b := v128_bytes(x)
	for i := range 4 {
		if binary.LittleEndian.Uint32(b[i*4:]) == 0 {
			return 0
		}
	}
	return 1
}

func i32x4_bitmask(x v128) int32 {
	b := v128_bytes(x)
	var r int32
	for i := range 4 {
		r |= int32(binary.LittleEndian.Uint32(b[i*4:])>>31) << i
	}
	return r
}

func i32x4_add(x, y v128) v128 {
	return fromI32(asI32(x).Add(asI32(y)))
}

func i8x16_narrow_i16x8_s(x, y v128) v128 {
	xb, yb := v128_bytes(x), v128_bytes(y)
	var b [16]byte
	for i := range 8 {
		b[i] = byte(min(max(int32(int16(binary.LittleEndian.Uint16(xb[i*2:]))), -128), 127))
		b[i+8] = byte(min(max(int32(int16(binary.LittleEndian.Uint16(yb[i*2:]))), -128), 127))
	}
	return v128_from_bytes(b)
}

func i8x16_narrow_i16x8_u(x, y v128) v128 {
	xb, yb := v128_bytes(x), v128_bytes(y)
	var b [16]byte
	for i := range 8 {
		b[i] = byte(min(max(int32(int16(binary.LittleEndian.Uint16(xb[i*2:]))), 0), 255))
		b[i+8] = byte(min(max(int32(int16(binary.LittleEndian.Uint16(yb[i*2:]))), 0), 255))
	}
	return v128_from_bytes(b)
}

func i16x8_narrow_i32x4_s(x, y v128) v128 {
	xb, yb := v128_bytes(x), v128_bytes(y)
	var b [16]byte
	for i := range 4 {
		xv := int64(int32(binary.LittleEndian.Uint32(xb[i*4:])))
		yv := int64(int32(binary.LittleEndian.Uint32(yb[i*4:])))
		binary.LittleEndian.PutUint16(b[i*2:], uint16(min(max(xv, -32768), 32767)))
		binary.LittleEndian.PutUint16(b[(i+4)*2:], uint16(min(max(yv, -32768), 32767)))
	}
	return v128_from_bytes(b)
}

func i16x8_narrow_i32x4_u(x, y v128) v128 {
	xb, yb := v128_bytes(x), v128_bytes(y)
	var b [16]byte
	for i := range 4 {
		xv := int64(int32(binary.LittleEndian.Uint32(xb[i*4:])))
		yv := int64(int32(binary.LittleEndian.Uint32(yb[i*4:])))
		binary.LittleEndian.PutUint16(b[i*2:], uint16(min(max(xv, 0), 65535)))
		binary.LittleEndian.PutUint16(b[(i+4)*2:], uint16(min(max(yv, 0), 65535)))
	}
	return v128_from_bytes(b)
}

func i16x8_extend_low_i8x16_s(x v128) v128 {
	return fromI16(asI8(x).ExtendLo8ToInt16())
}

func i16x8_extend_high_i8x16_s(x v128) v128 {
	b := v128_bytes(x)
	var v [8]int16
	for i := range v {
		v[i] = int16(int8(b[i+8]))
	}
	return fromI16(archsimd.LoadInt16x8Array(&v))
}

func i16x8_extend_low_i8x16_u(x v128) v128 {
	return fromU16(x.ExtendLo8ToUint16())
}

func i16x8_extend_high_i8x16_u(x v128) v128 {
	b := v128_bytes(x)
	var v [8]uint16
	for i := range v {
		v[i] = uint16(b[i+8])
	}
	return fromU16(archsimd.LoadUint16x8Array(&v))
}

func i32x4_extend_low_i16x8_s(x v128) v128 {
	return fromI32(asI16(x).ExtendLo4ToInt32())
}

func i32x4_extend_high_i16x8_s(x v128) v128 {
	b := v128_bytes(x)
	var v [4]int32
	for i := range v {
		v[i] = int32(int16(binary.LittleEndian.Uint16(b[(i+4)*2:])))
	}
	return fromI32(archsimd.LoadInt32x4Array(&v))
}

func i32x4_extend_low_i16x8_u(x v128) v128 {
	return fromU32(asU16(x).ExtendLo4ToUint32())
}

func i32x4_extend_high_i16x8_u(x v128) v128 {
	b := v128_bytes(x)
	var v [4]uint32
	for i := range v {
		v[i] = uint32(binary.LittleEndian.Uint16(b[(i+4)*2:]))
	}
	return fromU32(archsimd.LoadUint32x4Array(&v))
}

func i64x2_extend_low_i32x4_s(x v128) v128 {
	return fromI64(asI32(x).ExtendLo2ToInt64())
}

func i64x2_extend_high_i32x4_s(x v128) v128 {
	b := v128_bytes(x)
	v := [2]int64{
		int64(int32(binary.LittleEndian.Uint32(b[8:]))),
		int64(int32(binary.LittleEndian.Uint32(b[12:]))),
	}
	return fromI64(archsimd.LoadInt64x2Array(&v))
}

func i64x2_extend_low_i32x4_u(x v128) v128 {
	return fromU64(asU32(x).ExtendLo2ToUint64())
}

func i64x2_extend_high_i32x4_u(x v128) v128 {
	b := v128_bytes(x)
	v := [2]uint64{
		uint64(binary.LittleEndian.Uint32(b[8:])),
		uint64(binary.LittleEndian.Uint32(b[12:])),
	}
	return fromU64(archsimd.LoadUint64x2Array(&v))
}

func i64x2_shl(x v128, y int32) v128 {
	return fromU64(asU64(x).ShiftAllLeft(uint64(y & 63)))
}

func f64x2_add(x, y v128) v128 {
	return fromF64(asF64(x).Add(asF64(y)))
}

func f64x2_mul(x, y v128) v128 {
	return fromF64(asF64(x).Mul(asF64(y)))
}
