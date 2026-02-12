package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tokenize(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		raw    string
		expect []Token
	}{
		{
			name:   "empty",
			raw:    "",
			expect: []Token{},
		},
		{
			name:   "one",
			raw:    "aa",
			expect: []Token{NewTokenWord("aa")},
		},
		{
			name:   "one with spaces",
			raw:    "  aa  ",
			expect: []Token{NewTokenWord("aa")},
		},
		{
			name:   "two",
			raw:    "aa   bb",
			expect: []Token{NewTokenWord("aa"), NewTokenWord("bb")},
		},
		{
			name: "random",
			raw:  "   abc   cdv asds   cdcd  ",
			expect: []Token{
				NewTokenWord("abc"),
				NewTokenWord("cdv"),
				NewTokenWord("asds"),
				NewTokenWord("cdcd"),
			},
		},
		{
			name: "consecutive spaces ",
			raw:  "echo hello    world",
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("hello"),
				NewTokenWord("world"),
			},
		},
		{
			name: "single quotes",
			raw:  "echo 'hello    world'",
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("hello    world"),
			},
		},
		{
			name: "Adjacent quoted strings concatenated",
			raw:  "echo 'hello''world'",
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("helloworld"),
			},
		},
		{
			name: "Empty quotes ignored",
			raw:  "echo hello''world",
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("helloworld"),
			},
		},
		{
			name: "WTF1",
			raw:  "echo hello'world",
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("helloworld"),
			},
		},
		{
			name: "WTF2",
			raw:  "echo 'text''''",
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("text"),
			},
		},
		{
			name: "double quotes/Multiple spaces preserved",
			raw:  `echo "hello    world"`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("hello    world"),
			},
		},
		{
			name: "double quotes/concatenated",
			raw:  `echo "hello""world"`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("helloworld"),
			},
		},
		{
			name: "double quotes/separate",
			raw:  `echo "hello" "world"`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("hello"),
				NewTokenWord("world"),
			},
		},
		{
			name: "double quotes/Single quotes inside are literal",
			raw:  `echo "shell's test"`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("shell's test"),
			},
		},
		{
			name: "backslash/thee spaces",
			raw:  `echo three\ \ \ spaces`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("three   spaces"),
			},
		},
		{
			name: "backslash/first space",
			raw:  `echo before\     after`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("before "),
				NewTokenWord("after"),
			},
		},
		{
			name: "backslash/just n",
			raw:  `echo test\nexample`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("testnexample"),
			},
		},
		{
			name: "backslash/two",
			raw:  `echo hello\\world`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord(`hello\world`),
			},
		},
		{
			name: "backslash/escape quote",
			raw:  `echo \'hello\'`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("'hello'"),
			},
		},
		{
			name: "backslashes in single quotes",
			raw:  `echo 'shell\\\nscript'`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord(`shell\\\nscript`),
			},
		},
		{
			name: "backslashes in single quotes/second",
			raw:  `echo 'example\"test'`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord(`example\"test`),
			},
		},
		{
			name: "Backslashes in Double Quotes/1",
			raw:  `echo "A \\ escapes itself"`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord(`A \ escapes itself`),
			},
		},
		{
			name: "Backslashes in Double Quotes/2",
			raw:  `echo "A \" inside double quotes"`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord(`A " inside double quotes`),
			},
		},
		{
			name: "Redirect Stdout",
			raw:  `echo > file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenRedirect("file"),
			},
		},
		{
			name: "Redirect Stdout/stick right",
			raw:  `echo >file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenRedirect("file"),
			},
		},
		{
			name: "Redirect Stdout/stick",
			raw:  `echo>file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenRedirect("file"),
			},
		},
		{
			name: "Redirect Stdout/left exist",
			raw:  `echo hello>file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("hello"),
				NewTokenRedirect("file"),
			},
		},
		{
			name: "Redirect Stdout/left right exist",
			raw:  `echo  > file hello`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenRedirect("file"),
				NewTokenWord("hello"),
			},
		},
		{
			name: "Redirect Stdout Number",
			raw:  `echo 1> file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenRedirect("file"),
			},
		},
		{
			name: "Redirect Stdout Number/stick right",
			raw:  `echo 1>file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenRedirect("file"),
			},
		},
		{
			name: "Redirect Stdout Number/left right exist",
			raw:  `echo  1> file hello`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenRedirect("file"),
				NewTokenWord("hello"),
			},
		},
		{
			name: "Redirect Stdout Number/stuck right",
			raw:  `echo  hello1>file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("hello1"),
				NewTokenRedirect("file"),
			},
		},
		{
			name: "Redirect Stdout Number/words",
			raw:  `echo  hello 1`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("hello"),
				NewTokenRedirect("1"),
			},
		},
		{
			name: "Redirect Stderr",
			raw:  `echo 2> file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenRedirectErr("file"),
			},
		},
		{
			name: "Redirect Stderr/stick right",
			raw:  `echo 2>file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenRedirectErr("file"),
			},
		},
		{
			name: "Redirect Stderr/left right exist",
			raw:  `echo  2> file hello`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenRedirectErr("file"),
				NewTokenWord("hello"),
			},
		},
		{
			name: "Redirect Stdout/stuck all",
			raw:  `echo  hello2>file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("hello2"),
				NewTokenRedirect("file"),
			},
		},
		{
			name: "Append Redirect",
			raw:  `echo >> file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenAppendRedirect("file"),
			},
		},
		{
			name: "Append Redirect/stick right",
			raw:  `echo >>file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenAppendRedirect("file"),
			},
		},
		{
			name: "Append Redirect/left right exist",
			raw:  `echo  >> file hello`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenAppendRedirect("file"),
				NewTokenWord("hello"),
			},
		},
		{
			name: "Append Redirect/stuck all",
			raw:  `echo  hello>>file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenWord("hello"),
				NewTokenAppendRedirect("file"),
			},
		}, ///
		{
			name: "Append Redirect Err",
			raw:  `echo 2>> file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenAppendRedirectErr("file"),
			},
		},
		{
			name: "Append Redirect/stick right",
			raw:  `echo 2>>file`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenAppendRedirectErr("file"),
			},
		},
		{
			name: "Append Redirect/left right exist",
			raw:  `echo  2>> file hello`,
			expect: []Token{
				NewTokenWord("echo"),
				NewTokenAppendRedirectErr("file"),
				NewTokenWord("hello"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, _ := Tokenize(tc.raw)

			assert.Equal(t, tc.expect, result)
		})
	}
}
