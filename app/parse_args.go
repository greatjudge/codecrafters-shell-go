package main

import (
	"io"
	"strings"
)

type runeType int

const (
	unknownRuneType        runeType = 0
	spaceRuneType          runeType = 1
	eofRuneType            runeType = 2
	nonEscapeQuoteRuneType runeType = 3
)

type parserState int

const (
	startState       parserState = 1
	inWordState      parserState = 2
	inNonEscapeQuote parserState = 3
)

func parseArgs(raw string) ([]string, error) {
	var (
		reader = strings.NewReader(raw)
		state  = startState
		token  = []rune{}
		args   = []string{}
	)

loop:
	for {
		nextRune, _, err := reader.ReadRune()
		nextRuneType := classifyRune(nextRune)

		if err == io.EOF {
			nextRuneType = eofRuneType
		} else if err != nil {
			return nil, err
		}

		switch state {
		case startState:
			switch nextRuneType {
			case spaceRuneType:
				{
				}
			case eofRuneType:
				{
					state = startState
					break loop
				}
			case nonEscapeQuoteRuneType:
				{
					state = inNonEscapeQuote
				}
			default:
				{
					token = append(token, nextRune)
					state = inWordState
				}
			}
		case inWordState:
			switch nextRuneType {
			case eofRuneType:
				{
					args = append(args, string(token))
					token = []rune{}
					state = startState
					break loop
				}
			case spaceRuneType:
				{
					args = append(args, string(token))
					token = []rune{}
					state = startState
				}
			case nonEscapeQuoteRuneType:
				{
					state = inNonEscapeQuote
				}
			default:
				{
					token = append(token, nextRune)
				}
			}
		case inNonEscapeQuote:
			switch nextRuneType {
			case eofRuneType:
				{
					args = append(args, string(token))
					token = []rune{}
					state = startState
					break loop
				}
			case nonEscapeQuoteRuneType:
				{
					state = inWordState
				}
			default:
				{
					token = append(token, nextRune)
				}
			}
		}
	}

	return args, nil
}

func classifyRune(r rune) runeType {
	typeByRune := map[rune]runeType{
		' ':  spaceRuneType,
		'\'': nonEscapeQuoteRuneType,
	}

	return typeByRune[r]
}
