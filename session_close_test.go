package controlwire

import (
	"encoding/json"
	"testing"
)

// Closing a session is an explicit act on that session, and the component that owns it is the only
// one that can perform it: the core does not write an owner's store. One name for it, for the same
// reason the report has one.
func TestTheCloseHasOneNameEveryOwnerServes(t *testing.T) {
	if SessionCloseCommand == "" {
		t.Fatal("the close has no command name")
	}
	for _, taken := range []string{HelloCommand, SessionsCommand} {
		if SessionCloseCommand == taken {
			t.Fatalf("the close took the name %q", taken)
		}
	}
}

// The close names one session. A close that took a list would leave a caller reading a partial
// failure with no way to say which halves happened.
func TestTheCloseNamesOneSession(t *testing.T) {
	body, err := json.Marshal(SessionCloseRequest{Session: "7"})
	if err != nil {
		t.Fatal(err)
	}
	var back SessionCloseRequest
	if err := json.Unmarshal(body, &back); err != nil {
		t.Fatal(err)
	}
	if back.Session != "7" {
		t.Fatalf("the request named %q", back.Session)
	}
}

// The answer states whether the record is gone. A close of a session the owner does not hold is not
// a failure — the outcome a caller wanted is the outcome it has — and it says so.
func TestTheAnswerStatesWhetherTheRecordIsGone(t *testing.T) {
	body, err := json.Marshal(SessionCloseResult{Session: "7", Closed: true, Held: false})
	if err != nil {
		t.Fatal(err)
	}
	var back map[string]any
	if err := json.Unmarshal(body, &back); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"session", "closed", "held"} {
		if _, present := back[name]; !present {
			t.Errorf("the answer states no %q", name)
		}
	}
}
