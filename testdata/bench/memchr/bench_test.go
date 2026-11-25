package memchrbench

import (
	"bytes"
	"fmt"
	"testing"

	"memchrbench/scalar"
	"memchrbench/simd"
)

func BenchmarkMemchr(b *testing.B) {
	s, v := scalar.New(), simd.New()
	for _, m := range []struct {
		name string
		buf  []byte
		find func(int32, int32) int32
	}{
		{"scalar", (*s.Xmemory().Slice())[s.Xinput():][:65536], s.Xfind},
		{"simd", (*v.Xmemory().Slice())[v.Xinput():][:65536], v.Xfind},
	} {
		for _, n := range []int{64, 4096, 65536} {
			b.Run(fmt.Sprintf("%s/%d", m.name, n), func(b *testing.B) {
				buf := m.buf[:n]
				for i := range buf {
					buf[i] = byte(i % 127)
				}
				const needle = 0xff
				buf[n-1] = needle
				want := int32(bytes.IndexByte(buf, needle))
				if got := m.find(needle, int32(n)); got != want {
					b.Fatalf("got %d; want %d", got, want)
				}
				b.SetBytes(int64(n))
				b.ReportAllocs()
				for b.Loop() {
					m.find(needle, int32(n))
				}
			})
		}
	}
}
