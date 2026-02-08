package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommand(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		raw    string
		expect Command
	}{
		{
			name: "several spaces",
			raw:  "echo arg   other",
			expect: Command{
				Name: "echo",
				Args: []string{
					"echo",
					"arg",
					"other",
				},
			},
		},
		{
			name: "with single quotes",
			raw:  "echo 'hello    world' other",
			expect: Command{
				Name: "echo",
				Args: []string{
					"echo",
					"hello    world",
					"other",
				},
			},
		},
		{
			name: "adjacent quoted strings",
			raw:  "echo 'hello''world'",
			expect: Command{
				Name: "echo",
				Args: []string{
					"echo",
					"helloworld",
				},
			},
		},
		{
			name: "empty quotes '' are ignored.",
			raw:  "echo hello''world",
			expect: Command{
				Name: "echo",
				Args: []string{
					"echo",
					"helloworld",
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := ParseCommand(tc.raw)

			assert.NoError(t, err)
			assert.Equal(t, tc.expect, result)
		})
	}
}

func Test_parseArgs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		raw    string
		expect []string
	}{
		{
			name:   "empty",
			raw:    "",
			expect: []string{},
		},
		{
			name:   "one",
			raw:    "aa",
			expect: []string{"aa"},
		},
		{
			name:   "one with spaces",
			raw:    "  aa  ",
			expect: []string{"aa"},
		},
		{
			name:   "two",
			raw:    "aa   bb",
			expect: []string{"aa", "bb"},
		},
		{
			name:   "random",
			raw:    "   abc   cdv asds   cdcd  ",
			expect: []string{"abc", "cdv", "asds", "cdcd"},
		},
		{
			name:   "consecutive spaces ",
			raw:    "echo hello    world",
			expect: []string{"echo", "hello", "world"},
		},
		{
			name:   "single quotes",
			raw:    "echo 'hello    world'",
			expect: []string{"echo", "hello    world"},
		},
		{
			name:   "Adjacent quoted strings concatenated",
			raw:    "echo 'hello''world'",
			expect: []string{"echo", "helloworld"},
		},
		{
			name:   "Empty quotes ignored",
			raw:    "echo hello''world",
			expect: []string{"echo", "helloworld"},
		},
		{
			name:   "WTF1",
			raw:    "echo hello'world",
			expect: []string{"echo", "helloworld"},
		},
		{
			name:   "WTF2",
			raw:    "echo 'text''''",
			expect: []string{"echo", "text"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, _ := parseArgs(tc.raw)

			assert.Equal(t, tc.expect, result)
		})
	}
}
