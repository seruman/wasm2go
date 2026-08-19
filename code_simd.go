package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

func (t *translator) readOpcodeSIMD(fn *funcCompiler) error {
	t.simd = true
	code, err := readLEB128(t.in)
	if err != nil {
		return err
	}
	switch code {
	case 0x00: // v128.load
		return fn.simdLoad("v128_load")
	case 0x05: // v128.load32x2_s
		return fn.simdLoad("v128_load32x2_s")
	case 0x06: // v128.load32x2_u
		return fn.simdLoad("v128_load32x2_u")
	case 0x07: // v128.load8_splat
		return fn.simdLoad("v128_load8_splat")
	case 0x08: // v128.load16_splat
		return fn.simdLoad("v128_load16_splat")
	case 0x09: // v128.load32_splat
		return fn.simdLoad("v128_load32_splat")
	case 0x0a: // v128.load64_splat
		return fn.simdLoad("v128_load64_splat")
	case 0x0b: // v128.store
		return fn.simdStore("v128_store")
	case 0x0c: // v128.const
		expr, err := t.constV128()
		if err != nil {
			return err
		}
		fn.pushConst(expr)
	case 0x0d: // i8x16.shuffle
		var lanes [16]ast.Expr
		for i := range lanes {
			lane, err := t.in.ReadByte()
			if err != nil {
				return err
			}
			lanes[i] = &ast.BasicLit{Kind: token.INT, Value: formatUint(uint64(lane))}
		}
		y := fn.pop()
		x := fn.pop()
		args := append([]ast.Expr{x, y}, lanes[:]...)
		fn.push(fn.simdCall("i8x16_shuffle", args...))
	case 0x0f: // i8x16.splat
		fn.uniHelper("i8x16_splat")
	case 0x10: // i16x8.splat
		fn.uniHelper("i16x8_splat")
	case 0x11: // i32x4.splat
		fn.uniHelper("i32x4_splat")
	case 0x12: // i64x2.splat
		fn.uniHelper("i64x2_splat")
	case 0x13: // f32x4.splat
		fn.uniHelper("f32x4_splat")
	case 0x14: // f64x2.splat
		fn.uniHelper("f64x2_splat")
	case 0x15: // i8x16.extract_lane_s
		return fn.simdExtractLane("i8x16_extract_lane_s")
	case 0x16: // i8x16.extract_lane_u
		return fn.simdExtractLane("i8x16_extract_lane_u")
	case 0x17: // i8x16.replace_lane
		return fn.simdReplaceLane("i8x16_replace_lane")
	case 0x18: // i16x8.extract_lane_s
		return fn.simdExtractLane("i16x8_extract_lane_s")
	case 0x19: // i16x8.extract_lane_u
		return fn.simdExtractLane("i16x8_extract_lane_u")
	case 0x1a: // i16x8.replace_lane
		return fn.simdReplaceLane("i16x8_replace_lane")
	case 0x1b: // i32x4.extract_lane
		return fn.simdExtractLane("i32x4_extract_lane")
	case 0x1c: // i32x4.replace_lane
		return fn.simdReplaceLane("i32x4_replace_lane")
	case 0x1d: // i64x2.extract_lane
		return fn.simdExtractLane("i64x2_extract_lane")
	case 0x1e: // i64x2.replace_lane
		return fn.simdReplaceLane("i64x2_replace_lane")
	case 0x1f: // f32x4.extract_lane
		return fn.simdExtractLane("f32x4_extract_lane")
	case 0x20: // f32x4.replace_lane
		return fn.simdReplaceLane("f32x4_replace_lane")
	case 0x21: // f64x2.extract_lane
		return fn.simdExtractLane("f64x2_extract_lane")
	case 0x22: // f64x2.replace_lane
		return fn.simdReplaceLane("f64x2_replace_lane")
	case 0x23: // i8x16.eq
		fn.binHelper("i8x16_eq")
	case 0x24: // i8x16.ne
		fn.binHelper("i8x16_ne")
	case 0x25: // i8x16.lt_s
		fn.binHelper("i8x16_lt_s")
	case 0x26: // i8x16.lt_u
		fn.binHelper("i8x16_lt_u")
	case 0x27: // i8x16.gt_s
		fn.binHelper("i8x16_gt_s")
	case 0x28: // i8x16.gt_u
		fn.binHelper("i8x16_gt_u")
	case 0x29: // i8x16.le_s
		fn.binHelper("i8x16_le_s")
	case 0x2a: // i8x16.le_u
		fn.binHelper("i8x16_le_u")
	case 0x2b: // i8x16.ge_s
		fn.binHelper("i8x16_ge_s")
	case 0x2c: // i8x16.ge_u
		fn.binHelper("i8x16_ge_u")
	case 0x2d: // i16x8.eq
		fn.binHelper("i16x8_eq")
	case 0x2e: // i16x8.ne
		fn.binHelper("i16x8_ne")
	case 0x2f: // i16x8.lt_s
		fn.binHelper("i16x8_lt_s")
	case 0x30: // i16x8.lt_u
		fn.binHelper("i16x8_lt_u")
	case 0x31: // i16x8.gt_s
		fn.binHelper("i16x8_gt_s")
	case 0x32: // i16x8.gt_u
		fn.binHelper("i16x8_gt_u")
	case 0x33: // i16x8.le_s
		fn.binHelper("i16x8_le_s")
	case 0x34: // i16x8.le_u
		fn.binHelper("i16x8_le_u")
	case 0x35: // i16x8.ge_s
		fn.binHelper("i16x8_ge_s")
	case 0x36: // i16x8.ge_u
		fn.binHelper("i16x8_ge_u")
	case 0x37: // i32x4.eq
		fn.binHelper("i32x4_eq")
	case 0x38: // i32x4.ne
		fn.binHelper("i32x4_ne")
	case 0x39: // i32x4.lt_s
		fn.binHelper("i32x4_lt_s")
	case 0x3a: // i32x4.lt_u
		fn.binHelper("i32x4_lt_u")
	case 0x3b: // i32x4.gt_s
		fn.binHelper("i32x4_gt_s")
	case 0x3c: // i32x4.gt_u
		fn.binHelper("i32x4_gt_u")
	case 0x3d: // i32x4.le_s
		fn.binHelper("i32x4_le_s")
	case 0x3e: // i32x4.le_u
		fn.binHelper("i32x4_le_u")
	case 0x3f: // i32x4.ge_s
		fn.binHelper("i32x4_ge_s")
	case 0x40: // i32x4.ge_u
		fn.binHelper("i32x4_ge_u")
	case 0x4d: // v128.not
		fn.uniHelper("v128_not")
	case 0x4e: // v128.and
		fn.binHelper("v128_and")
	case 0x4f: // v128.andnot
		fn.binHelper("v128_andnot")
	case 0x50: // v128.or
		fn.binHelper("v128_or")
	case 0x51: // v128.xor
		fn.binHelper("v128_xor")
	case 0x52: // v128.bitselect
		c := fn.pop()
		y := fn.pop()
		x := fn.pop()
		fn.push(fn.simdCall("v128_bitselect", x, y, c))
	case 0x53: // v128.any_true
		fn.uniHelper("v128_any_true")
	case 0x54: // v128.load8_lane
		return fn.simdLoadLane("v128_load8_lane")
	case 0x55: // v128.load16_lane
		return fn.simdLoadLane("v128_load16_lane")
	case 0x56: // v128.load32_lane
		return fn.simdLoadLane("v128_load32_lane")
	case 0x57: // v128.load64_lane
		return fn.simdLoadLane("v128_load64_lane")
	case 0x58: // v128.store8_lane
		return fn.simdStoreLane("v128_store8_lane")
	case 0x59: // v128.store16_lane
		return fn.simdStoreLane("v128_store16_lane")
	case 0x5a: // v128.store32_lane
		return fn.simdStoreLane("v128_store32_lane")
	case 0x5b: // v128.store64_lane
		return fn.simdStoreLane("v128_store64_lane")
	case 0x5c: // v128.load32_zero
		return fn.simdLoad("v128_load32_zero")
	case 0x5d: // v128.load64_zero
		return fn.simdLoad("v128_load64_zero")
	case 0x63: // i8x16.all_true
		fn.uniHelper("i8x16_all_true")
	case 0x64: // i8x16.bitmask
		fn.uniHelper("i8x16_bitmask")
	case 0x65: // i8x16.narrow_i16x8_s
		fn.binHelper("i8x16_narrow_i16x8_s")
	case 0x66: // i8x16.narrow_i16x8_u
		fn.binHelper("i8x16_narrow_i16x8_u")
	case 0x77: // i8x16.min_u
		fn.binHelper("i8x16_min_u")
	case 0x83: // i16x8.all_true
		fn.uniHelper("i16x8_all_true")
	case 0x84: // i16x8.bitmask
		fn.uniHelper("i16x8_bitmask")
	case 0x85: // i16x8.narrow_i32x4_s
		fn.binHelper("i16x8_narrow_i32x4_s")
	case 0x86: // i16x8.narrow_i32x4_u
		fn.binHelper("i16x8_narrow_i32x4_u")
	case 0x87: // i16x8.extend_low_i8x16_s
		fn.uniHelper("i16x8_extend_low_i8x16_s")
	case 0x88: // i16x8.extend_high_i8x16_s
		fn.uniHelper("i16x8_extend_high_i8x16_s")
	case 0x89: // i16x8.extend_low_i8x16_u
		fn.uniHelper("i16x8_extend_low_i8x16_u")
	case 0x8a: // i16x8.extend_high_i8x16_u
		fn.uniHelper("i16x8_extend_high_i8x16_u")
	case 0x97: // i16x8.min_u
		fn.binHelper("i16x8_min_u")
	case 0xa3: // i32x4.all_true
		fn.uniHelper("i32x4_all_true")
	case 0xa4: // i32x4.bitmask
		fn.uniHelper("i32x4_bitmask")
	case 0xa7: // i32x4.extend_low_i16x8_s
		fn.uniHelper("i32x4_extend_low_i16x8_s")
	case 0xa8: // i32x4.extend_high_i16x8_s
		fn.uniHelper("i32x4_extend_high_i16x8_s")
	case 0xa9: // i32x4.extend_low_i16x8_u
		fn.uniHelper("i32x4_extend_low_i16x8_u")
	case 0xaa: // i32x4.extend_high_i16x8_u
		fn.uniHelper("i32x4_extend_high_i16x8_u")
	case 0xae: // i32x4.add
		fn.binHelper("i32x4_add")
	case 0xc7: // i64x2.extend_low_i32x4_s
		fn.uniHelper("i64x2_extend_low_i32x4_s")
	case 0xc8: // i64x2.extend_high_i32x4_s
		fn.uniHelper("i64x2_extend_high_i32x4_s")
	case 0xc9: // i64x2.extend_low_i32x4_u
		fn.uniHelper("i64x2_extend_low_i32x4_u")
	case 0xca: // i64x2.extend_high_i32x4_u
		fn.uniHelper("i64x2_extend_high_i32x4_u")
	case 0xcb: // i64x2.shl
		fn.binHelper("i64x2_shl")
	case 0xcd: // i64x2.shr_u
		fn.binHelper("i64x2_shr_u")
	case 0xd6: // i64x2.eq
		fn.binHelper("i64x2_eq")
	case 0xd7: // i64x2.ne
		fn.binHelper("i64x2_ne")
	case 0xf0: // f64x2.add
		fn.binHelper("f64x2_add")
	case 0xf2: // f64x2.mul
		fn.binHelper("f64x2_mul")
	default:
		return fmt.Errorf("unsupported opcode (SIMD): 0xFD 0x%02X", code)
	}
	return nil
}

func (fn *funcCompiler) simdCall(name string, args ...ast.Expr) *ast.CallExpr {
	fn.helpers.add(name)
	return &ast.CallExpr{Fun: newID(name), Args: args}
}

func (fn *funcCompiler) simdMemArg() (uint64, error) {
	_, err := readLEB128(fn.translator.in) // align
	if err != nil {
		return 0, err
	}
	return readLEB128(fn.translator.in)
}

func (fn *funcCompiler) simdLoad(name string) error {
	offset, err := fn.simdMemArg()
	if err != nil {
		return err
	}
	addr := fn.popAddr(offset)
	fn.push(fn.simdCall(name, fn.memory.selector, addr))
	return nil
}

func (fn *funcCompiler) simdStore(name string) error {
	offset, err := fn.simdMemArg()
	if err != nil {
		return err
	}
	val := fn.pop()
	addr := fn.popAddr(offset)
	fn.emit(&ast.ExprStmt{X: fn.simdCall(name, fn.memory.selector, addr, val)})
	return nil
}

func (fn *funcCompiler) simdLoadLane(name string) error {
	offset, err := fn.simdMemArg()
	if err != nil {
		return err
	}
	lane, err := fn.translator.in.ReadByte()
	if err != nil {
		return err
	}
	vec := fn.pop()
	addr := fn.popAddr(offset)
	laneExpr := &ast.BasicLit{Kind: token.INT, Value: formatUint(uint64(lane))}
	fn.push(fn.simdCall(name, fn.memory.selector, addr, vec, laneExpr))
	return nil
}

func (fn *funcCompiler) simdStoreLane(name string) error {
	offset, err := fn.simdMemArg()
	if err != nil {
		return err
	}
	lane, err := fn.translator.in.ReadByte()
	if err != nil {
		return err
	}
	vec := fn.pop()
	addr := fn.popAddr(offset)
	laneExpr := &ast.BasicLit{Kind: token.INT, Value: formatUint(uint64(lane))}
	fn.emit(&ast.ExprStmt{X: fn.simdCall(name, fn.memory.selector, addr, vec, laneExpr)})
	return nil
}

func (fn *funcCompiler) simdExtractLane(name string) error {
	lane, err := fn.translator.in.ReadByte()
	if err != nil {
		return err
	}
	fn.push(fn.simdCall(name, fn.pop(), &ast.BasicLit{Kind: token.INT, Value: formatUint(uint64(lane))}))
	return nil
}

func (fn *funcCompiler) simdReplaceLane(name string) error {
	lane, err := fn.translator.in.ReadByte()
	if err != nil {
		return err
	}
	val := fn.pop()
	vec := fn.pop()
	fn.push(fn.simdCall(name, vec, &ast.BasicLit{Kind: token.INT, Value: formatUint(uint64(lane))}, val))
	return nil
}
