package controlwire

import (
	"errors"
	"regexp"
)

var ErrInvalidProcessRole = errors.New("invalid process role")
var processRolePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,17}$`)

// FormatProcessName combines one run label and one component role into the complete operating-
// system process name. The ASCII budgets guarantee the result fits Darwin's 31-byte proc_name.
func FormatProcessName(label, role string) (string, error) {
	if _, err := ParseProcessLabel(label); err != nil {
		return "", err
	}
	if !processRolePattern.MatchString(role) {
		return "", ErrInvalidProcessRole
	}
	return label + "-" + role, nil
}
