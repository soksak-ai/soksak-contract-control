package controlwire

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
)

func Address(runtimeRoot, identifier string, windows bool) string {
	if !windows {
		return filepath.Join(runtimeRoot, identifier+".sock")
	}
	digest := sha256.Sum256([]byte(strings.ToLower(filepath.Clean(runtimeRoot)) + "\x00" + identifier))
	return `\\.\pipe\soksak-control-` + hex.EncodeToString(digest[:16])
}
