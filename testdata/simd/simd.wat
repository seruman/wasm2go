(module
  (memory 1)
  (global $mask v128 (v128.const i8x16 -1 -1 -1 -1 -1 -1 -1 -1 0 0 0 0 0 0 0 0))

  (func (export "global") (result i32)
    global.get $mask
    i8x16.bitmask)

  (func (export "block") (result i32)
    block (result v128)
      v128.const i8x16 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1 -1
    end
    i8x16.bitmask)

  (func (export "mask") (result i32)
    v128.const i8x16 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15
    v128.const i8x16 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15
    i8x16.eq
    i8x16.bitmask)

  (func (export "load_store") (result i32)
    i32.const 32
    v128.const i8x16 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15
    v128.store offset=3
    i32.const 33
    v128.load 0 offset=2
    i8x16.extract_lane_u 15)

  (func (export "load_lane") (result i32)
    i32.const 64
    i32.const 171
    i32.store8
    i32.const 64
    v128.const i8x16 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
    v128.load8_lane 7
    i8x16.extract_lane_u 7)

  (func (export "store_lane") (result i32)
    i32.const 96
    v128.const i8x16 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15
    v128.store8_lane 15
    i32.const 96
    i32.load8_u)

  (func (export "shr_u") (result i64)
    v128.const i64x2 -1 -9223372036854775808
    i32.const 4
    i64x2.shr_u
    i64x2.extract_lane 0))
