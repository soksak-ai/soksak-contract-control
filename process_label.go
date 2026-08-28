package controlwire

import "errors"

import "regexp"

const ProcessLabelEnvironment = "SOKSAK_PROCESS_LABEL"
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
