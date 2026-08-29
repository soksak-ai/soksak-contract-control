package controlwire

import "errors"

import "regexp"

const ProcessLabelEnvironment = "SOKSAK_PROCESS_LABEL"

// SidecarNameEnvironment carries the installer-materialized process name. It is the sole source
// for identity-scoped PTY endpoints; a daemon must not infer it from its executable path.
const SidecarNameEnvironment = "SOKSAK_SIDECAR_NAME"

// SidecarBindingsEnvironment is a JSON object whose keys are installed Sidecar component ids and
// whose values are their materialized process names for the current environment.
const SidecarBindingsEnvironment = "SOKSAK_SIDECAR_BINDINGS"
const DefaultProcessLabel = "soksak"

var ErrInvalidProcessLabel = errors.New("invalid process label")
var processLabelPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,30}$`)

// ParseProcessLabel validates one explicit launch label. An empty environment value is resolved to
// the canonical label by the launcher before it reaches this contract.
func ParseProcessLabel(value string) (string, error) {
	if !processLabelPattern.MatchString(value) {
		return "", ErrInvalidProcessLabel
	}
	return value, nil
}
