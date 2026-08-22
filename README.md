# soksak-contract-control

Shared control-envelope wire for Soksak Core and process sidecars.

The contract defines the greeting, request, response, announcement and command-table shapes. Domain
payloads remain with their owning contracts.

## Verification

```sh
go test ./...
go vet ./...
```
