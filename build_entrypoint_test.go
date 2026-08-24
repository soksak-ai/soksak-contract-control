package controlwire

import (
	"os"
	"strings"
	"testing"
)

func TestBuildEntrypointProjectsGoAndRustOwners(t *testing.T) {
	makefile, err := os.ReadFile("Makefile")
	if err != nil {
		t.Fatal(err)
	}
	source := string(makefile)
	for _, target := range []string{"preflight:", "prepare:", "build:", "verify:"} {
		if !strings.Contains(source, target) {
			t.Errorf("Makefile omits %s", target)
		}
	}
	workflow, err := os.ReadFile(".github/workflows/verify.yml")
	if err != nil {
		t.Fatal(err)
	}
	text := string(workflow)
	for _, required := range []string{"go-version-file: go.mod", "rust-toolchain.toml", "make verify"} {
		if !strings.Contains(text, required) {
			t.Errorf("workflow omits %s", required)
		}
	}
	for _, duplicate := range []string{"go-version: \"1.25.0\"", "toolchain: \"1.96.0\"", "go test ./...", "cargo test --locked", "cache-dependency-path: go.sum"} {
		if strings.Contains(text, duplicate) {
			t.Errorf("workflow duplicates owner command or metadata: %s", duplicate)
		}
	}
}
