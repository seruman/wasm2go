//go:build go1.27 && goexperiment.simd

package helpers

import "simd/archsimd"

// The translator uses 0x80 to zero lanes from the other operand (VPSHUFB).
func i8x16_lookup(x v128, mask *[16]byte) v128 {
	return x.PermuteOrZero(archsimd.LoadUint8x16Array(mask).BitsToInt8())
}
