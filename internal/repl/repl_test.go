package main

import (
	"fmt"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  I like to workout   ",
			expected: []string{"i", "like", "to", "workout"},
		},
		{
			input:    "  Making a new project   ",
			expected: []string{"making", "a", "new", "project"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		expected := c.expected

		fmt.Println("about to compare lengths")
		if len(actual) != len(expected) {
			fmt.Print("lengths do not match")
			t.Errorf("expected lenght %d, actual length %d", len(actual), len(expected))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("actual word: %s, expected word: %s", word, expectedWord)
			}
		}
	}
}
