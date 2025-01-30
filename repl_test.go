package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO   WORLD\tGO",
			expected: []string{"hello", "world", "go"},
		},
		{
			input:    "Trailing   spaces   ",
			expected: []string{"trailing", "spaces"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("CleanInput(%q) = %v, want %v", c.input, actual, c.expected)
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("CleanInput(%q) = %v, want %v", c.input, word, expectedWord)
			}
		}
	}
}
