package controlwire

import (
	"encoding/json"
	"os"
	"testing"
)

func TestProcessNameVectorsAreTheGoContract(t *testing.T) {
	body, err := os.ReadFile("process-name-vectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var vectors []struct {
		Label string `json:"label"`
		Role  string `json:"role"`
		Name  string `json:"name"`
		Valid bool   `json:"valid"`
	}
	if err := json.Unmarshal(body, &vectors); err != nil {
		t.Fatal(err)
	}
	for _, vector := range vectors {
		name, err := FormatProcessName(vector.Label, vector.Role)
		if vector.Valid && (err != nil || name != vector.Name) {
			t.Errorf("FormatProcessName(%q, %q) = %q, %v; want %q", vector.Label, vector.Role, name, err, vector.Name)
		}
		if !vector.Valid && err == nil {
			t.Errorf("FormatProcessName(%q, %q) accepted %q", vector.Label, vector.Role, name)
		}
	}
}
