(module
  (func (export "block") (result i32)
    block (result v128)
      unreachable
    end
    drop
    i32.const 0))
