package core

import "testing"

func TestRunFilterPipeline(t *testing.T) {
	tests := []struct {
		msg      string
		expected bool
	}{
		{"BUY NOW!!!", true},
		{"Hi, how are you?", false},
	}

	for _, tt := range tests {
		result := RunFilterPipeline(tt.msg)
		if result != tt.expected {
			t.Errorf("RunFilterPipeline(%q) = %v; want %v", tt.msg, result, tt.expected)
		}
	}
}