package controlwire

import (
	"encoding/json"
	"os"
	"sort"
	"testing"
)

// SessionCommand is one command a view calls on the core to maintain the session index.
//
// The index is written by the component that has both halves — the owner issued the id and the view
// knows the coordinate — and read by the core. Neither side can check the other's spelling: the
// coupling law forbids a repository from reading another's source, so a name that differs by one
// character produces an empty index and no error.
//
// Measured 2026-09-04: a terminal kit called session_attach with viewId. The command is
// session.attach with view. Both mismatches were silent, because the runner answers a refusal
// rather than raising one.
type SessionCommand struct {
	Command string   `json:"command"`
	Params  []string `json:"params"`
	Purpose string   `json:"purpose"`
}

// SessionCommands is the set both sides read. A caller spells a name from here; a server registers
// one from here.
func SessionCommands(t *testing.T) []SessionCommand {
	t.Helper()
	body, err := os.ReadFile("session-command-vectors.json")
	if err != nil {
		t.Fatalf("reading the session command vectors: %v", err)
	}
	var commands []SessionCommand
	if err := json.Unmarshal(body, &commands); err != nil {
		t.Fatalf("parsing the session command vectors: %v", err)
	}
	return commands
}

func TestTheSessionCommandVectorsAreWellFormed(t *testing.T) {
	commands := SessionCommands(t)
	if len(commands) == 0 {
		t.Fatal("the vectors name no command")
	}

	seen := map[string]bool{}
	for _, one := range commands {
		if one.Command == "" || one.Purpose == "" {
			t.Fatalf("a vector names no command or no purpose: %+v", one)
		}
		if seen[one.Command] {
			t.Fatalf("%s appears twice", one.Command)
		}
		seen[one.Command] = true

		// A dotted name. The core serves session_attach on its Go registry and session.attach on
		// its command catalog; a view reaches the catalog, so the dotted name is the one both sides
		// have to agree on.
		if !isDottedSessionName(one.Command) {
			t.Fatalf("%s is not a dotted session command name", one.Command)
		}

		parameters := append([]string(nil), one.Params...)
		sort.Strings(parameters)
		for index, name := range parameters {
			if name == "" {
				t.Fatalf("%s declares an empty parameter", one.Command)
			}
			if index > 0 && parameters[index-1] == name {
				t.Fatalf("%s declares %s twice", one.Command, name)
			}
		}
	}
}

func isDottedSessionName(name string) bool {
	if len(name) < len("session.x") || name[:len("session.")] != "session." {
		return false
	}
	for _, letter := range name[len("session."):] {
		if letter < 'a' || letter > 'z' {
			return false
		}
	}
	return true
}
