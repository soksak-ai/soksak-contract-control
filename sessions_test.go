package controlwire

import (
	"encoding/json"
	"testing"
)

// Every component that owns sessions answers this one command. The core holds an index of what
// exists and each owner holds records, and a process exit between them changes one and not the
// other; reconciling the two is the same question whatever the sessions are.
//
// One name for it. A name per owner would make the core know which owner it is addressing before it
// could ask, which is the coupling this envelope exists to remove.
func TestTheSessionReportHasOneNameEveryOwnerServes(t *testing.T) {
	if SessionsCommand == "" {
		t.Fatal("the session report has no command name")
	}
	if SessionsCommand == HelloCommand {
		t.Fatal("the session report took the greeting's name")
	}
}

// The caller names what it holds. An owner answering from its own store alone could never report a
// session whose record is gone, and that is the one a caller most needs to hear about.
func TestTheRequestNamesWhatTheCallerHolds(t *testing.T) {
	body, err := json.Marshal(SessionsRequest{Sessions: []string{"7", "8"}})
	if err != nil {
		t.Fatal(err)
	}
	var back SessionsRequest
	if err := json.Unmarshal(body, &back); err != nil {
		t.Fatal(err)
	}
	if len(back.Sessions) != 2 {
		t.Fatalf("the request named %d sessions", len(back.Sessions))
	}
}

// A session id is the owner's, in the form the owner issued. The envelope carries it as text and
// reads nothing out of it: a number would make the envelope decide what an id may be, and a large
// one does not survive a generic JSON parse intact.
func TestASessionIdTravelsAsTextTheEnvelopeDoesNotRead(t *testing.T) {
	const large = "101552085244916000000"
	body, err := json.Marshal(SessionReport{
		Complete: true,
		Sessions: []SessionOutcome{{Session: large, Outcome: SessionFull}},
	})
	if err != nil {
		t.Fatal(err)
	}
	var back SessionReport
	if err := json.Unmarshal(body, &back); err != nil {
		t.Fatal(err)
	}
	if back.Sessions[0].Session != large {
		t.Fatalf("the id came back as %q", back.Sessions[0].Session)
	}
}

// Four outcomes and no fifth. A caller has no handling for a value it does not know, and an absent
// one is a caller guessing.
func TestTheOutcomesAreFour(t *testing.T) {
	seen := map[string]bool{}
	for _, outcome := range []string{SessionFull, SessionDegraded, SessionFailed, SessionLost} {
		if outcome == "" {
			t.Fatal("an outcome has no name")
		}
		if seen[outcome] {
			t.Fatalf("two outcomes share the name %q", outcome)
		}
		seen[outcome] = true
	}
	if len(seen) != 4 {
		t.Fatalf("the envelope names %d outcomes", len(seen))
	}
}

// An unfinished report is not a final one. A caller that took it as final would count a session lost
// that the owner had not looked for yet.
func TestAnUnfinishedReportSaysSo(t *testing.T) {
	body, err := json.Marshal(SessionReport{Complete: false})
	if err != nil {
		t.Fatal(err)
	}
	var back SessionReport
	if err := json.Unmarshal(body, &back); err != nil {
		t.Fatal(err)
	}
	if back.Complete {
		t.Fatal("an unfinished report round-tripped as finished")
	}
}
