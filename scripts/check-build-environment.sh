#!/bin/sh
set -eu

[ "$#" -eq 0 ] || { echo 'BUILD_DECLARATION_INVALID: usage: check-build-environment.sh' >&2; exit 78; }
root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
go_expected=$(awk '$1 == "go" { value="go" $2; count++ } END { if (count == 1) print value; else exit 1 }' "$root/go.mod" 2>/dev/null || true)
go_actual=$(go env GOVERSION 2>/dev/null || true)
go_host_os=$(go env GOHOSTOS 2>/dev/null || true)
go_host_arch=$(go env GOHOSTARCH 2>/dev/null || true)
rust_expected=$(sed -n 's/^channel = "\([^"]*\)"$/\1/p' "$root/rust-toolchain.toml")
rust_actual=$(rustc --version 2>/dev/null | awk '{print $2}' || true)
rust_host=$(rustc -vV 2>/dev/null | sed -n 's/^host: //p' || true)

case "$(uname -s)-$(uname -m)" in
  Darwin-arm64) required_go_os=darwin; required_go_arch=arm64; required_rust=aarch64-apple-darwin ;;
  Darwin-x86_64) if [ "$(sysctl -n hw.optional.arm64 2>/dev/null || true)" = 1 ]; then required_go_os=darwin; required_go_arch=arm64; required_rust=aarch64-apple-darwin; else required_go_os=darwin; required_go_arch=amd64; required_rust=x86_64-apple-darwin; fi ;;
  Linux-aarch64|Linux-arm64) required_go_os=linux; required_go_arch=arm64; required_rust=aarch64-unknown-linux-gnu ;;
  Linux-x86_64) required_go_os=linux; required_go_arch=amd64; required_rust=x86_64-unknown-linux-gnu ;;
  MINGW*-x86_64|MSYS*-x86_64|CYGWIN*-x86_64) required_go_os=windows; required_go_arch=amd64; required_rust=x86_64-pc-windows-msvc ;;
  *) echo 'TOOLCHAIN_MISMATCH: unsupported host' >&2; exit 78 ;;
esac

if [ "$go_actual" != "$go_expected" ] || [ "$go_host_os" != "$required_go_os" ] || [ "$go_host_arch" != "$required_go_arch" ] || \
   [ "$rust_actual" != "$rust_expected" ] || [ "$rust_host" != "$required_rust" ]; then
  printf 'TOOLCHAIN_MISMATCH: expected go=%s/%s/%s rust=%s/%s; actual go=%s/%s/%s rust=%s/%s\n' \
    "${go_expected:-missing}" "$required_go_os" "$required_go_arch" "${rust_expected:-missing}" "$required_rust" \
    "${go_actual:-missing}" "${go_host_os:-unknown}" "${go_host_arch:-unknown}" "${rust_actual:-missing}" "${rust_host:-unknown}" >&2
  exit 78
fi
printf 'BUILD_ENVIRONMENT_READY go=%s/%s/%s rust=%s/%s\n' "$go_actual" "$go_host_os" "$go_host_arch" "$rust_actual" "$rust_host"
