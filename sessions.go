package controlwire

// SessionsCommand is what every component that owns sessions answers.
//
// The core holds an index of what exists and each owner holds records. A process exit between them
// changes one and not the other, and reconciling the two is the same question whatever the sessions
// are: a shell, a page, anything a person or a program started that outlives the view showing it.
//
// One name for it. A name per owner would make the core know which owner it is addressing before it
// could ask, and knowing that means knowing the owner's contract — the coupling this envelope
// exists to remove. An owner that serves no sessions refuses the name with a reason, which is a
// different thing from not having it.
//
// It starts no process and ends none. A caller reconciling its index must be able to ask without
// changing what it is asking about.
const SessionsCommand = "system.sessions"

// The outcomes a start can end in for one session. A caller reads exactly one of these per session
// it named: a fifth value is one a caller has no handling for, and an absent one is a caller
// guessing.
const (
	// SessionFull is a record whose owner marked it as stopped on purpose. Its state is what that
	// stop wrote.
	SessionFull = "full"
	// SessionDegraded is a record no stop marked, or one holding the creation facts alone. It ends
	// wherever the last write reached, and nothing states where it was meant to end.
	SessionDegraded = "degraded"
	// SessionFailed is a record that exists and could not be used. The owner keeps it rather than
	// deleting it: it is the only evidence of what was lost, and a later start may stand it up.
	SessionFailed = "failed"
	// SessionLost is a session the owner holds no record for. Nothing recovers it.
	SessionLost = "lost"
)

// SessionsRequest asks what became of the sessions a caller holds.
//
// The caller names them because only the caller knows what it held. An owner answering from its own
// store alone could never report a session whose record is gone, and that session is exactly the
// one a caller needs to hear about. An empty Sessions asks for every session the owner knows of,
// which is what a caller with no index of its own reads.
type SessionsRequest struct {
	Sessions []string `json:"sessions,omitempty"`
}

// SessionOutcome is what became of one session.
//
// The id travels as text and this envelope reads nothing out of it. Its form is the owner's: a
// number here would make the envelope decide what an id may be, and a large one does not survive a
// generic JSON parse intact — a parser holding it as a float64 is exact only to 2^53.
type SessionOutcome struct {
	Session string `json:"session"`
	Outcome string `json:"outcome"`
	// Reason states what stopped a failed restore. Empty for every other outcome.
	Reason string `json:"reason,omitempty"`
}

// SessionReport is the answer to SessionsRequest.
type SessionReport struct {
	// Complete states that the owner finished reading its store. A caller that took an unfinished
	// report as final would count a session lost that the owner had not looked for yet.
	Complete bool             `json:"complete"`
	Sessions []SessionOutcome `json:"sessions"`
}
