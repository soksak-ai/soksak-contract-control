package controlwire

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type processLabelVector struct {
	Input string `json:"input"`
	Valid bool   `json:"valid"`
}

func processLabelVectors(t *testing.T) []processLabelVector {
	t.Helper()
	body, err := os.ReadFile("process-label-vectors.json")
	if err != nil {
		t.Fatal(err)
	}
	var vectors []processLabelVector
	if err := json.Unmarshal(body, &vectors); err != nil {
		t.Fatal(err)
	}
	return vectors
}

func TestProcessLabelVectorsAreTheGoContract(t *testing.T) {
	for _, vector := range processLabelVectors(t) {
		label, err := ParseProcessLabel(vector.Input)
		if vector.Valid && (err != nil || label != vector.Input) {
			t.Errorf("valid label %q = %q, %v", vector.Input, label, err)
		}
		if !vector.Valid && err == nil {
			t.Errorf("invalid label %q was accepted as %q", vector.Input, label)
		}
	}
}

func TestProtocolTwoAnnouncementRequiresAProcessLabel(t *testing.T) {
	if Protocol != 2 {
		t.Fatalf("control protocol = %d, want 2 for process-label announcements", Protocol)
	}
	field, found := reflect.TypeOf(Announcement{}).FieldByName("ProcessLabel")
	if !found || field.Tag.Get("json") != "processLabel" {
		t.Fatal("announcement does not carry the required processLabel")
	}
}
