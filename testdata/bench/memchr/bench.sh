#!/bin/sh
# Requires Cargo, the wasm32-unknown-unknown Rust target, and Go 1.27+.
set -eu
src=$(cd "$(dirname "$0")" && pwd)
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
(cd "$src/../../.." && go build -o "$tmp/wasm2go" .)
cp -R "$src/src" "$src/Cargo.toml" "$src/Cargo.lock" "$src/bench_test.go" "$tmp/"
cd "$tmp"
go mod init memchrbench
for variant in scalar simd; do
    case "$variant" in
        scalar) feature=-simd128 ;;
        simd)   feature=+simd128 ;;
    esac
    RUSTFLAGS="-Ctarget-feature=$feature" cargo build --locked --release --target wasm32-unknown-unknown
    mkdir "$variant"
    ./wasm2go -pkg "$variant" -o "$variant/search.go" target/wasm32-unknown-unknown/release/search.wasm
done
GOEXPERIMENT=simd go test -run '^$' -bench . "$@"
