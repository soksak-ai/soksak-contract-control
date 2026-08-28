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
An absent launch value resolves once to the canonical label `soksak`.

The label never identifies an installation, resolves a dependency, grants permission, selects a
socket, or establishes ownership. It exists so public process inventory and monitoring UI can
distinguish two otherwise identical runs. A process must announce the validated label it accepted;
silently inheriting an unread environment value is not success.

An operating-system process name is a platform projection of process metadata, not component
identity. On macOS, the standard C library's `setprogname` changes the name returned by `proc_name`;
changing `argv[0]`, an `NSProcessInfo` value, or a pthread name is not equivalent. The macOS
projection is limited to 31 observable bytes and `setprogname` retains the supplied pointer for the
process lifetime. Wails application and WebKit helper display names are a separate build-owned
surface derived from bundle product metadata. A platform implementation must therefore prove its
projection against the operating-system observation API without changing component IDs, executable
paths, routing, or dependency identity.

## Verification

```sh
make verify
```

`go.mod` and `rust-toolchain.toml` are the exact Go and Rust owners. Make verifies the shared Go/Rust
wire vectors without invoking a consumer implementation. The workflow injects those owner versions
and calls the same command; it does not duplicate language commands or cache a nonexistent lockfile.
