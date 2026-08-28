# soksak-contract-control

Shared control-envelope wire for Soksak Core and process sidecars.

The contract defines the greeting, request, response, announcement, command-table, and local address
shapes. Go and Rust implementations consume the same address vectors. Unix addresses are filesystem
sockets; Windows addresses are opaque current-machine named pipes derived from runtime identity.
Domain payloads remain with their owning contracts.

## Process labels

Protocol 2 announcements and greetings carry `processLabel`. The launcher supplies the same label
through `SOKSAK_PROCESS_LABEL` to every process it owns. A label is an ASCII diagnostic token of 1
through 64 characters: it starts with an alphanumeric character and the rest may also contain `.`,
`_`, or `-`. Go and Rust consume `process-label-vectors.json` as the shared verdict.

The label never identifies an installation, resolves a dependency, grants permission, selects a
socket, or establishes ownership. It exists only so operating-system process tools and public status
can distinguish two otherwise identical runs. A process must announce the label it actually applied;
inheriting an unused environment value is not success.

## Verification

```sh
make verify
```

`go.mod` and `rust-toolchain.toml` are the exact Go and Rust owners. Make verifies the shared Go/Rust
wire vectors without invoking a consumer implementation. The workflow injects those owner versions
and calls the same command; it does not duplicate language commands or cache a nonexistent lockfile.
