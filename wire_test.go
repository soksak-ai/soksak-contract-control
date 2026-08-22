package controlwire

import (
	"encoding/json"
	"testing"
)

func TestRequestAndResponseRoundTrip(t *testing.T) {
	request := Request{ID: "r1", Command: "system.hello", Args: map[string]json.RawMessage{"protocol": json.RawMessage("1")}}
	body, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Request
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.ID != request.ID || decoded.Command != request.Command {
		t.Fatalf("request changed: %+v", decoded)
	}
}
