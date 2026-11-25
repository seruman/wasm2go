//go:build go1.27 && goexperiment.simd

package helpers

import "simd/archsimd"

func i8x16_lookup(x v128, mask *[16]byte) v128 {
	return x.BitsToInt8().LookupOrZero(archsimd.LoadUint8x16Array(mask).BitsToInt8()).ToBits()
}
