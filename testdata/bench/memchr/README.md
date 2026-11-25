# Rust memchr: scalar vs SIMD

This benchmarks the same byte search from memchr 2.8.3, built with and without
the Rust `simd128` target feature and translated with the current checkout.
The needle is at the end of a 64-byte, 4-KiB, or 64-KiB buffer. The timed loop
does not instantiate modules, copy input, or call Rust through an FFI.

With Cargo, a Rust toolchain with the `wasm32-unknown-unknown` target,
and a SIMD-capable Go host:

```sh
rustup target add wasm32-unknown-unknown
GOTOOLCHAIN=go1.27.0 ./testdata/bench/memchr/bench.sh -benchtime=500ms -count=5
```

To collect a CPU profile of just the SIMD search, pass
`-bench='BenchmarkMemchr/simd/65536$' -cpuprofile="$PWD/memchr.cpu"`.
Cargo and Go build in a temporary directory; no Rust toolchain or crate is
needed for normal wasm2go builds and tests.

As a reference, the 64-KiB late-match search on an Apple M5 Pro (Go 1.27.0,
rustc 1.95.0, median of five 350-ms runs) took **8.34 µs** with scalar Wasm,
**5.77 µs** with SIMD translated by `e3caf1c`, and **1.75 µs** with the
optimized shared SIMD output. The last two used the same Wasm binary. All
three measured zero allocations. This tests a hot, cache-resident byte search,
not cold-memory throughput or an end-to-end codec.
