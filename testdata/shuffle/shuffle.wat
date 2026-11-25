(module
  (func (export "x") (param v128 v128) (result v128)
    local.get 0
    local.get 1
    i8x16.shuffle 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15)
  (func (export "y") (param v128 v128) (result v128)
    local.get 0
    local.get 1
    i8x16.shuffle 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31)
  (func (export "repeat") (param v128 v128) (result v128)
    local.get 0
    local.get 1
    i8x16.shuffle 7 7 7 7 7 7 7 7 7 7 7 7 7 7 7 7)
  (func (export "reverse_y") (param v128 v128) (result v128)
    local.get 0
    local.get 1
    i8x16.shuffle 31 30 29 28 27 26 25 24 23 22 21 20 19 18 17 16)
  (func (export "mixed") (param v128 v128) (result v128)
    local.get 0
    local.get 1
    i8x16.shuffle 31 0 30 1 29 2 28 3 27 4 26 5 25 6 24 7)
  (func (export "word") (param v128 v128) (result v128)
    local.get 0
    local.get 1
    i8x16.shuffle 14 15 30 31 2 3 10 11 4 5 16 17 4 5 8 9)
  (func (export "dword") (param v128 v128) (result v128)
    local.get 0
    local.get 1
    i8x16.shuffle 16 17 18 19 8 9 10 11 28 29 30 31 0 1 2 3)
  (func (export "half") (param v128 v128) (result v128)
    local.get 0
    local.get 1
    i8x16.shuffle 24 25 26 27 28 29 30 31 0 1 2 3 4 5 6 7)
)
