//go:build go1.27 && goexperiment.simd && (amd64 || arm64)

package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGeneratedShuffle(t *testing.T) {
	wasm, err := os.ReadFile("testdata/shuffle/shuffle.wasm")
	if err != nil {
		t.Fatal(err)
	}
	tests, err := os.ReadFile("testdata/shuffle/shuffle_test.go")
	if err != nil {
		t.Fatal(err)
	}
	oldArch, oldPkg := *simdarch, *pkg
	defer func() { *simdarch, *pkg = oldArch, oldPkg }()
	*pkg = "shuffle"
	dir := t.TempDir()
	for _, arch := range []string{"", runtime.GOARCH} {
		*simdarch = arch
		var out bytes.Buffer
		if err := translate(bytes.NewReader(wasm), &out); err != nil {
			t.Fatal(err)
		}
		sub := filepath.Join(dir, "shuffle"+arch)
		if err := os.Mkdir(sub, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(sub, "shuffle.go"), out.Bytes(), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(sub, "shuffle_test.go"), tests, 0644); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GO111MODULE=off", "GOTOOLCHAIN=local")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generated shuffle tests: %v\n%s", err, out)
	}
}
