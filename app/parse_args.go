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
	escapeQuoteRuneType    runeType = 4
	escapeRuneType         runeType = 5
)

func classifyRune(r rune) runeType {
	typeByRune := map[rune]runeType{
		' ':  spaceRuneType,
		'\'': nonEscapeQuoteRuneType,
		'"':  escapeQuoteRuneType,
		'\\': escapeRuneType,
	}

	return typeByRune[r]
}

type parserState int

const (
	startState            parserState = 1
	inWordState           parserState = 2
	inNonEscapeQuoteState parserState = 3
	inEscapeQuoteState    parserState = 4
	inEscapeState         parserState = 5
	inEscapeQuotedState   parserState = 6
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
					state = inNonEscapeQuoteState
				}
			case escapeQuoteRuneType:
				{
					state = inEscapeQuoteState
				}
			case escapeRuneType:
				{
					state = inEscapeState
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
					state = inNonEscapeQuoteState
				}
			case escapeQuoteRuneType:
				{
					state = inEscapeQuoteState
				}
			case escapeRuneType:
				{
					state = inEscapeState
				}
			default:
				{
					token = append(token, nextRune)
				}
			}
		case inNonEscapeQuoteState:
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
		case inEscapeQuoteState:
			switch nextRuneType {
			case eofRuneType:
				{
					args = append(args, string(token))
					token = []rune{}
					state = startState
					break loop
				}
			case escapeQuoteRuneType:
				{
					state = inWordState
				}
			case escapeRuneType:
				{
					state = inEscapeQuotedState
				}
			default:
				{
					token = append(token, nextRune)
				}
			}
		case inEscapeState:
			switch nextRuneType {
			case eofRuneType:
				{
					args = append(args, string(token))
					token = []rune{}
					state = startState
					break loop
				}
			default:
				token = append(token, nextRune)
				state = inWordState
			}
		case inEscapeQuotedState:
			switch nextRuneType {
			case eofRuneType:
				{
					args = append(args, string(token))
					token = []rune{}
					state = startState
					break loop
				}
			default:
				token = append(token, nextRune)
				state = inEscapeQuoteState
			}
		}
	}

	return args, nil
}
