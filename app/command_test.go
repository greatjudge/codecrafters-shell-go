package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		{
			name:   "double quotes/Multiple spaces preserved",
			raw:    `echo "hello    world"`,
			expect: []string{"echo", "hello    world"},
		},
		{
			name:   "double quotes/concatenated",
			raw:    `echo "hello""world"`,
			expect: []string{"echo", "helloworld"},
		},
		{
			name:   "double quotes/separate",
			raw:    `echo "hello" "world"`,
			expect: []string{"echo", "hello", "world"},
		},
		{
			name:   "double quotes/Single quotes inside are literal",
			raw:    `echo "shell's test"`,
			expect: []string{"echo", "shell's test"},
		},
		{
			name:   "backslash/thee spaces",
			raw:    `echo three\ \ \ spaces`,
			expect: []string{"echo", "three   spaces"},
		},
		{
			name:   "backslash/first space",
			raw:    `echo before\     after`,
			expect: []string{"echo", "before ", "after"},
		},
		{
			name:   "backslash/just n",
			raw:    `echo test\nexample`,
			expect: []string{"echo", "testnexample"},
		},
		{
			name:   "backslash/two",
			raw:    `echo hello\\world`,
			expect: []string{"echo", `hello\world`},
		},
		{
			name:   "backslash/escape quote",
			raw:    `echo \'hello\'`,
			expect: []string{"echo", "'hello'"},
		},
		{
			name:   "backslashes in single quotes",
			raw:    `echo 'shell\\\nscript'`,
			expect: []string{"echo", `shell\\\nscript`},
		},
		{
			name:   "backslashes in single quotes/second",
			raw:    `echo 'example\"test'`,
			expect: []string{"echo", `example\"test`},
		},
		{
			name:   "Backslashes in Double Quotes/1",
			raw:    `echo "A \\ escapes itself"`,
			expect: []string{"echo", `A \ escapes itself`},
		},
		{
			name:   "Backslashes in Double Quotes/2",
			raw:    `echo "A \" inside double quotes"`,
			expect: []string{"echo", `A " inside double quotes`},
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
