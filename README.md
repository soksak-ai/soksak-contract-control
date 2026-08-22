# soksak-contract-control

Shared control-envelope wire for Soksak Core and process sidecars.

The contract defines the greeting, request, response, announcement, command-table, and local address
shapes. Go and Rust implementations consume the same address vectors. Unix addresses are filesystem
sockets; Windows addresses are opaque current-machine named pipes derived from runtime identity.
Domain payloads remain with their owning contracts.

## Verification

```sh
go test ./...
go vet ./...
cargo test --locked
```
