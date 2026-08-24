SHELL := /bin/sh

.PHONY: preflight prepare build verify

preflight:
	@scripts/check-build-environment.sh

prepare: preflight
	@go mod download
	@cargo fetch --locked

build: prepare
	@go build ./...
	@cargo build --locked --release

verify: prepare
	@go mod tidy -diff
	@go test -count=1 ./...
	@go vet ./...
	@cargo test --locked
