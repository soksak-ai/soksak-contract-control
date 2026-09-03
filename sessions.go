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

// SessionCloseCommand ends one session, and every component that owns sessions answers it.
//
// Closing removes the owner's record, and the core does not write an owner's store. So the close is
// the owner's act and the core orders it — the same split as the report, and one name for the same
// reason: a caller must not have to know which owner it is addressing before it can ask.
//
// It is explicit. Closing a window, a workspace or a pane detaches the sessions inside it and closes
// none of them, so nothing but this ends a session.
const SessionCloseCommand = "system.sessionClose"

// SessionCloseRequest names the one session to end.
//
// One rather than a list. A close over a list leaves a caller reading a partial failure with no way
// to say which halves happened, and the caller wanted a definite answer per session.
type SessionCloseRequest struct {
	Session string `json:"session"`
}

// SessionCloseResult is what became of the close.
type SessionCloseResult struct {
	Session string `json:"session"`
	// Closed states that no record for this session remains. A session the owner never held is
	// closed as far as a caller is concerned: the outcome it wanted is the outcome it has.
	Closed bool `json:"closed"`
	// Held states that the owner was holding the session when the close arrived. It separates
	// ending a running session from finding nothing to end, which a caller reconciling an index
	// reads differently even though both leave no record.
	Held bool `json:"held"`
}

// SessionNoticeEvent is what an owner emits for a session it stood back up, and what the core
// delivers to whatever is attached to that session.
//
// An owner that restarted holds sessions whose state changed while it was gone. A consumer attached
// to one and told nothing resumes against state the session does not have and reports it as the
// session's own — a degraded restore read as a full one, silently.
//
// The core delivers it and does not act on it. What an attachment does with the news is that
// component's: the core owns which sessions exist, never what a session's content means.
const SessionNoticeEvent = "session.restored"

// SessionNotice is what that event holds.
//
// The outcome is one of the four a start ends in. A notice stating only that something changed
// leaves a consumer knowing the state moved and not how, which is the same as not knowing.
type SessionNotice struct {
	Session string `json:"session"`
	Owner   string `json:"owner"`
	Outcome string `json:"outcome"`
	// Reason states what stopped a failed restore. Empty for every other outcome.
	Reason string `json:"reason,omitempty"`
}
