package main

import "testing"

func TestClassifySecret(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "empty", value: "  ", want: "empty"},
		{name: "v2", value: "lwsec:v2:abc", want: "v2"},
		{name: "v1", value: "legacy-ciphertext", want: "v1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := classifySecret(tt.value); got != tt.want {
				t.Fatalf("classifySecret() = %q, want %q", got, tt.want)
			}
		})
	}
}
