package main

import (
	"testing"
)

func TestStringUnpacking(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true}, // начинается с цифры
		{"", "", false},  // пустая строка
		{"qwe\\4\\5", "qwe45", false},
		{"qwe\\45", "qwe44444", false},
		{"a12", "aaaaaaaaaaaa", false},
		{"\\", "", true},     // заканчивается на '\'
		{"a\\", "", true},    // заканчивается на '\'
		{"a0b", "ab", false}, // цифра 0 (не добавляет символ)
	}

	for _, tt := range tests {
		result, err := stringUnpacking(tt.input)
		if (err != nil) != tt.hasError {
			t.Errorf("input: %q, unexpected error status: got %v, want error? %v", tt.input, err, tt.hasError)
		}
		if err == nil && result != tt.expected {
			t.Errorf("input: %q, expected: %q, got: %q", tt.input, tt.expected, result)
		}
	}
}
