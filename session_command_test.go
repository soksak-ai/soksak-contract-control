package controlwire

import (
	"sort"
	"testing"
)

// The declared commands are read through the package accessor so a consumer reads exactly what this
// test grades. A second reader here would be a second answer to one question.
func sessionCommands(t *testing.T) []SessionCommand {
	t.Helper()
	commands, err := SessionCommands()
	if err != nil {
		t.Fatalf("reading the session command vectors: %v", err)
	}
	return commands
}

func TestTheSessionCommandVectorsAreWellFormed(t *testing.T) {
	commands := sessionCommands(t)
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
