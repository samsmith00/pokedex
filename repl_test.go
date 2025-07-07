package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/samsmith00/pokedex/internal/pokecache"
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

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}
