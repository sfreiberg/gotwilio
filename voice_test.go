package gotwilio

import (
	"encoding/json"
	"testing"
)

func TestVoiceResponseQueueTime(t *testing.T) {
	tests := []struct {
		name string
		body string
		want *string
	}{
		{"string", `{"queue_time": "1000"}`, strPtr("1000")},
		{"null", `{"queue_time": null}`, nil},
		{"missing", `{}`, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var vr VoiceResponse
			if err := json.Unmarshal([]byte(tt.body), &vr); err != nil {
				t.Fatal(err)
			}
			if (vr.QueueTime == nil) != (tt.want == nil) || (vr.QueueTime != nil && *vr.QueueTime != *tt.want) {
				t.Errorf("QueueTime = %v, want %v", vr.QueueTime, tt.want)
			}
		})
	}
}

func strPtr(s string) *string { return &s }
