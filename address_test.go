package controlwire

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddressVectors(t *testing.T) {
	body, err := os.ReadFile("address-vectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var vectors []struct {
		Runtime, Identifier, Address string
		Windows                      bool
	}
	if err := json.Unmarshal(body, &vectors); err != nil {
		t.Fatal(err)
	}
	for _, vector := range vectors {
		if got := Address(vector.Runtime, vector.Identifier, vector.Windows); got != vector.Address {
			t.Errorf("Address(%q, %q, %t) = %q, want %q", vector.Runtime, vector.Identifier, vector.Windows, got, vector.Address)
		}
	}
}

func TestAddressUsesRuntimePathOnUnix(t *testing.T) {
	got := Address(filepath.Join("tmp", "runtime"), "com.soksak.test", false)
	if got != filepath.Join("tmp", "runtime", "com.soksak.test.sock") {
		t.Fatalf("address = %s", got)
	}
}

func TestAddressUsesStableOpaqueWindowsPipe(t *testing.T) {
	first := Address(`C:\runtime\one`, "com.soksak.test", true)
	second := Address(`c:\RUNTIME\one`, "com.soksak.test", true)
	if first != second || !strings.HasPrefix(first, `\\.\pipe\soksak-control-`) || len(first) != len(`\\.\pipe\soksak-control-`)+32 {
		t.Fatalf("addresses = %q %q", first, second)
	}
	if first == Address(`C:\runtime\two`, "com.soksak.test", true) {
		t.Fatal("runtime isolation did not change the Windows pipe")
	}
}
