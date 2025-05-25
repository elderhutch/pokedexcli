package main

import (
	"strings"
	"testing"
)

func cleanInput(text string) []string {
	// Split the input string by spaces and remove empty strings
	words := []string{}
	for _, word := range strings.Fields(text) {
		if word != "" {
			words = append(words, word)
		}
	}
	return words

}

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		// add more cases here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Expected %s, but got %s", expectedWord, word)
			}
		}
	}
}
