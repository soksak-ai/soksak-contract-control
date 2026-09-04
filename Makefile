SHELL := /bin/sh
registry_flags = --@soksak:registry=$(REGISTRY) --@soksak-ai:registry=$(REGISTRY) --config.minimum-release-age=0
publish_flags = --registry "$(REGISTRY)" --@soksak:registry="$(REGISTRY)" --@soksak-ai:registry="$(REGISTRY)" --no-git-checks
# REGISTRY is accepted from the make command line only ($(origin) must be "command line").
# GNU make's own environment channels (MAKEFLAGS, GNUMAKEFLAGS, MAKEFILES, -e) are outside this
# Makefile's control and are not refused; setting them is a deliberate act of the caller.

.PHONY: preflight guard lock prepare build verify require-registry publish

preflight:
	@scripts/check-build-environment.sh

lock: guard preflight
	@cargo metadata --format-version 1 > /dev/null
	@CI=1 PNPM_DISABLE_SELF_UPDATE_CHECK=1 pnpm install --lockfile-only $(if $(findstring command line,$(origin REGISTRY)),$(registry_flags))

# A failed install exits with the pnpm status; the pnpm-workspace.yaml digest is compared only after
# a successful install.
prepare: guard preflight
	@go mod download
	@cargo fetch --locked
	@before=$$(shasum -a 256 pnpm-workspace.yaml); CI=1 PNPM_DISABLE_SELF_UPDATE_CHECK=1 pnpm install --frozen-lockfile $(if $(findstring command line,$(origin REGISTRY)),$(registry_flags)) || exit $$?; test "$$before" = "$$(shasum -a 256 pnpm-workspace.yaml)" || { echo 'pnpm install rewrote pnpm-workspace.yaml' >&2; exit 65; }

build: prepare
	@go build ./...
	@cargo build --locked --release
	@CI=1 PNPM_DISABLE_SELF_UPDATE_CHECK=1 pnpm $(if $(findstring command line,$(origin REGISTRY)),$(registry_flags)) typecheck

verify: prepare
	@go mod tidy -diff
	@go test -count=1 ./...
	@go vet ./...
	@cargo test --locked
	@CI=1 PNPM_DISABLE_SELF_UPDATE_CHECK=1 pnpm $(if $(findstring command line,$(origin REGISTRY)),$(registry_flags)) test
	@CI=1 PNPM_DISABLE_SELF_UPDATE_CHECK=1 pnpm $(if $(findstring command line,$(origin REGISTRY)),$(registry_flags)) typecheck

# The Node distribution of this contract. REGISTRY is accepted from the make command line
# only ($(origin) must be "command line").
guard:
	@case "$(origin REGISTRY)" in undefined|"command line") ;; *) echo 'REGISTRY from the $(origin REGISTRY) is refused: make publish REGISTRY=http://host:port/' >&2; exit 64 ;; esac
	@case "$(origin REGISTRY):$(REGISTRY)" in undefined:|"command line:http://"*|"command line:https://"*) ;; *) echo 'REGISTRY must be an absolute URL: make publish REGISTRY=http://host:port/' >&2; exit 64 ;; esac

# This repository is a build input in three languages, not a component installed into a soksak home.
# Go and Rust consumers resolve it from git; Node consumers resolve it from the registry by name and
# version. There is no portable release: the SDK's release builder takes exactly one of package.json
# or Cargo.toml, and this repository has both.
require-registry: guard
	@test "$(origin REGISTRY)" = "command line" || { echo 'REGISTRY must be given on the make command line: make publish REGISTRY=http://host:port/' >&2; exit 64; }

publish: require-registry verify
	@CI=1 PNPM_DISABLE_SELF_UPDATE_CHECK=1 pnpm publish $(publish_flags)
