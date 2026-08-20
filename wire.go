// Package controlwire is the control-plane envelope, and nothing that answers on it.
//
// It exists apart from every implementation because it has more than one. An application answers on
// this envelope; a unit that runs in its own process answers on the same one, over its own socket.
// A wire defined inside one of them makes the other copy it, and two copies of a wire diverge
// without failing — they arrive as a different answer.
//
// One line of JSON per message, both directions. A length prefix would make the socket unreadable by
// hand, which is the difference between a control plane someone can operate and one they can only
// write a client for.
package controlwire

import "encoding/json"

// Protocol is the version of this envelope.
//
// A client asking for a version the other side does not have is refused during the greeting rather
// than at the first command that behaves differently: a mismatch found halfway through a session has
// already produced answers the caller trusted.
const Protocol = 1

// HelloCommand is the greeting's name. It is reserved — anything that registered it would replace
// the negotiation with something that answers differently.
const HelloCommand = "system.hello"

// Request is one command, as it arrives.
type Request struct {
	// ID is echoed on the answer. The caller chooses it and nothing on the other side interprets
	// it, so a client may pipeline and match on its own terms.
	ID string `json:"id"`
	// Command names an entry in the answering side's registry. Nothing outside it exists.
	Command string `json:"command"`
	// Args are the command's arguments, still encoded. Decoding happens per command, so this
	// boundary never has to know their shapes.
	Args map[string]json.RawMessage `json:"args,omitempty"`
	// Language is the language this caller reads. Empty means the language agreed in the greeting,
	// and no greeting means English.
	//
	// It rides on the request rather than only on the connection because one process serves several
	// readers, and the language is the asker's rather than the socket's.
	Language string `json:"language,omitempty"`
}

// Response is one answer.
//
// Ok is explicit rather than inferred from Error being empty: a command whose result is null and one
// that failed with an empty message would otherwise be the same three bytes on the wire.
type Response struct {
	ID     string `json:"id"`
	Ok     bool   `json:"ok"`
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

// Answer is the shape a generic caller parses every result as.
//
// Two shapes reaching one socket, with nothing in the answer stating which, forces a client to know
// who owns each command before it can parse the reply. Measured 2026-08-17: two readings in one
// session were taken against the wrong shape and reported the opposite of what was on screen.
type Answer struct {
	// Code is machine-readable. "OK" for a result, and a name for a refusal — a caller acts on this
	// and shows Error to a person.
	Code string `json:"code"`
	Data any    `json:"data"`
}

// Owner names who answers a command, and it travels in the greeting.
//
// It is on the wire because a client reads it: a command owned by the framework needs this host's
// window and cannot be answered headless, and a caller that could not tell which is which would ask
// a window-less process for a window's answer and read the refusal as a missing feature.
type Owner string

const (
	OwnerCore      Owner = "core"
	OwnerFramework Owner = "framework"
	OwnerPlugin    Owner = "plugin"
)

// Served describes a command the answering side answers.
type Served struct {
	Name  string `json:"name"`
	Owner Owner  `json:"owner"`
}

// Unserved describes a command the answering side refuses, and why.
type Unserved struct {
	Name string `json:"name"`
	// BlockedBy separates "not written yet" from "impossible here". A caller that receives only
	// "unknown command" re-investigates settled ground, or imitates the command.
	BlockedBy string `json:"blockedBy"`
}

// Table is what the answering side serves and what it refuses, together.
type Table struct {
	Commands []Served   `json:"commands"`
	Unserved []Unserved `json:"unserved"`
}

// Greeting is what HelloCommand answers.
type Greeting struct {
	// Protocol is what the answering side speaks, which the client compares to what it asked for.
	Protocol int `json:"protocol"`
	// Identity names the installation, so a client that found the wrong socket receives that at the
	// greeting rather than through surprising answers.
	Identity string `json:"identity"`
	// Commands is what is served and refused, with reasons. Sent in the greeting because a client
	// that must ask separately will act on a name it has not checked.
	Commands Table `json:"commands"`
	// Language is what will be answered in for this session, and Languages is everything available.
	// A client that asked for one that is absent is told here, rather than receiving sentences that
	// quietly stayed English.
	Language  string   `json:"language"`
	Languages []string `json:"languages"`
}

// Announcement is the first line a process that answers on this envelope prints to stdout, and it is
// the only readiness signal this wire defines.
//
// A socket file is not one. The path exists from the moment of bind and also for as long as the
// filesystem holds one a dead process left behind, so a stat answers "a file is there" and never
// "someone is listening" — and turning the first answer into the second means looking again, which
// is a poll.
//
// Both fields are pointers so "absent" and "sent as a zero value" stay apart. A process that printed
// an unrelated JSON log line announces nothing; one that sent protocol 0 tried to announce and got
// it wrong, and collapsing those two turns every JSON-logging development server into a broken unit.
type Announcement struct {
	Protocol *int    `json:"protocol"`
	Socket   *string `json:"socket"`
	// Token is what the greeting on that socket has to carry, and it travels here because the only
	// process that reads this line is the one that started the process that printed it.
	//
	// The alternative is a file both sides derive a path to, and that is what a peer which did not
	// start the process uses. A starter reading the file instead would be deriving something it was
	// already told, and the two can part: a process that generated a token and had not yet written
	// the file is a process whose file says something else.
	//
	// Absent means the socket takes an unauthenticated greeting. That is a statement, not a
	// default: a process that wanted a token and forgot to announce one is refused at its own
	// greeting rather than reached by anyone.
	Token *string `json:"token,omitempty"`
}

// NewAnnouncement builds the line a process prints once its listener is bound and before it serves.
func NewAnnouncement(socket string) Announcement {
	protocol := Protocol
	return Announcement{Protocol: &protocol, Socket: &socket}
}

// WithToken is the same announcement, naming the token a greeting must carry.
func (announcement Announcement) WithToken(token string) Announcement {
	announcement.Token = &token
	return announcement
}
