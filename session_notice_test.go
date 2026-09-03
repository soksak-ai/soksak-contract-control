package controlwire

import (
	"encoding/json"
	"testing"
)

// An owner that restarted holds sessions whose state changed while it was gone. A consumer attached
// to one and told nothing resumes against state the session does not have, and reports it as the
// session's own.
func TestTheNoticeHasOneNameEveryOwnerEmits(t *testing.T) {
	if SessionNoticeEvent == "" {
		t.Fatal("the notice has no name")
	}
	for _, taken := range []string{HelloCommand, SessionsCommand, SessionCloseCommand} {
		if SessionNoticeEvent == taken {
			t.Fatalf("the notice took the name %q", taken)
		}
	}
}

// The notice states the session and what its restore ended in. Without the outcome a consumer knows
// the state changed and not how, which is the same as not knowing.
func TestTheNoticeStatesTheSessionAndTheOutcome(t *testing.T) {
	body, err := json.Marshal(SessionNotice{
		Session: "7", Owner: "pty", Outcome: SessionDegraded, Reason: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	var back map[string]any
	if err := json.Unmarshal(body, &back); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"session", "owner", "outcome"} {
		if _, present := back[name]; !present {
			t.Errorf("the notice states no %q", name)
		}
	}
}

// The outcome a notice carries is one of the four a start ends in. A fifth would be one no consumer
// has handling for.
func TestTheNoticeCarriesOneOfTheFourOutcomes(t *testing.T) {
	for _, outcome := range []string{SessionFull, SessionDegraded, SessionFailed, SessionLost} {
		notice := SessionNotice{Session: "7", Owner: "pty", Outcome: outcome}
		body, err := json.Marshal(notice)
		if err != nil {
			t.Fatal(err)
		}
		var back SessionNotice
		if err := json.Unmarshal(body, &back); err != nil {
			t.Fatal(err)
		}
		if back.Outcome != outcome {
			t.Fatalf("the outcome came back as %q", back.Outcome)
		}
	}
}
