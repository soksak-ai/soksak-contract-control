SHELL := /bin/sh

.PHONY: preflight lock prepare build verify

preflight:
	@scripts/check-build-environment.sh

lock: preflight
	@cargo metadata --format-version 1 > /dev/null

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
