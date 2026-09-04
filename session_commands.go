package controlwire

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

// The commands a view calls to write and read the core's session index, named once here so both
// sides grade themselves against one file rather than against each other.
//
// Measured 2026-09-04: a view called `session_attach` with `viewId` while the command it could run
// was `session.attach` with `view`. Either mismatch alone left the index empty, and neither was
// visible: the runner answers `{ok:false}` rather than raising, so a name nothing serves reads as a
// success to a caller that only catches.
//
//go:embed session-command-vectors.json
var sessionCommandVectors []byte

// SessionCommand is one command a view calls, with the exact parameters it takes.
type SessionCommand struct {
	Command string   `json:"command"`
	Params  []string `json:"params"`
	Purpose string   `json:"purpose"`
}

// SessionCommands returns the declared commands. A consumer grades its own registration or its own
// calls against these; it does not read another component's source.
func SessionCommands() ([]SessionCommand, error) {
	var commands []SessionCommand
	if err := json.Unmarshal(sessionCommandVectors, &commands); err != nil {
		return nil, fmt.Errorf("session command vectors: %w", err)
	}
	return commands, nil
}
