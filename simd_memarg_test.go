package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestSIMDExplicitMemoryIndex(t *testing.T) {
	wasm, err := os.ReadFile("testdata/simd/simd.wasm")
	if err != nil {
		t.Fatal(err)
	}
	// Xload_store exercises v128.load 0 offset=2 at runtime. Keep the
	// explicit-index encoding in the fixture even if a WAT assembler omits it.
	load := []byte{0xfd, 0x00, 0x44, 0x00, 0x02}
	i := bytes.Index(wasm, load)
	if i == -1 {
		t.Fatal("missing v128.load with explicit memory index 0 in SIMD fixture")
	}
	wasm[i+3] = 1
	err = translate(bytes.NewReader(wasm), io.Discard)
	if err == nil || !strings.Contains(err.Error(), "unsupported memory index 1") {
		t.Fatalf("got %v, want an error for memory index 1", err)
	}
}
