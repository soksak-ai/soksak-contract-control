# soksak-contract-control

Shared control-envelope wire for Soksak Core and process sidecars.

The contract defines the greeting, request, response, announcement, command-table, and local address
shapes. Go and Rust implementations consume the same address vectors. Unix addresses are filesystem
sockets; Windows addresses are opaque current-machine named pipes derived from runtime identity.
Domain payloads remain with their owning contracts.

## Verification

```sh
make verify
```

`go.mod` and `rust-toolchain.toml` are the exact Go and Rust owners. Make verifies the shared Go/Rust
wire vectors without invoking a consumer implementation. The workflow injects those owner versions
and calls the same command; it does not duplicate language commands or cache a nonexistent lockfile.
