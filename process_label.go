package controlwire

import "errors"

const ProcessLabelEnvironment = "SOKSAK_PROCESS_LABEL"

var ErrInvalidProcessLabel = errors.New("invalid process label")

// ParseProcessLabel validates one explicit launch label. An empty environment value is resolved to
// the canonical label by the launcher before it reaches this contract.
func ParseProcessLabel(string) (string, error) {
	return "", ErrInvalidProcessLabel
}
