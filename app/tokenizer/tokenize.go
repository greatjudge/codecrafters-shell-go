package tokenizer

import (
	"fmt"
	"io"
	"strings"
)

type runeType int

const (
	spaceRuneType          runeType = 1
	eofRuneType            runeType = 2
	nonEscapeQuoteRuneType runeType = 3
	escapeQuoteRuneType    runeType = 4
	escapeRuneType         runeType = 5
	redirectRuneType       runeType = 6
	oneRuneType            runeType = 7
	twoRuneType            runeType = 8
)

var (
	typeByRune = map[rune]runeType{
		' ':  spaceRuneType,
		'\'': nonEscapeQuoteRuneType,
		'"':  escapeQuoteRuneType,
		'\\': escapeRuneType,
		'>':  redirectRuneType,
		'1':  oneRuneType,
		'2':  twoRuneType,
	}
)

func classifyRune(r rune) runeType {
	return typeByRune[r]
}

type parserState int

const (
	startState                  parserState = 1
	inWordState                 parserState = 2
	inNonEscapeQuoteState       parserState = 3
	inEscapeQuoteState          parserState = 4
	inEscapeState               parserState = 5
	inEscapeQuotedState         parserState = 6
	inRedirectSymb              parserState = 7
	inRedirectReadName          parserState = 8
	inOneState                  parserState = 9
	inRedirectErrSymb           parserState = 10
	inRedirectErrReadName       parserState = 11
	inTwo                       parserState = 12
	inAppendRedirect            parserState = 13
	inAppendRedirectReadName    parserState = 14
	inAppendRedirectErr         parserState = 15
	inAppendRedirectErrReadName parserState = 16
)

func Tokenize(raw string) ([]Token, error) {
	var (
		reader     = strings.NewReader(raw)
		state      = startState
		tokenValue = []rune{}
		tokens     = []Token{}
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
			case redirectRuneType:
				{
					state = inRedirectSymb
				}
			case oneRuneType:
				{
					state = inOneState
					tokenValue = append(tokenValue, nextRune)
				}
			case twoRuneType:
				{
					state = inTwo
					tokenValue = append(tokenValue, nextRune)
				}
			default:
				{
					tokenValue = append(tokenValue, nextRune)
					state = inWordState
				}
			}
		case inWordState:
			switch nextRuneType {
			case eofRuneType:
				{
					tokens = append(tokens, NewTokenWord(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
					break loop
				}
			case spaceRuneType:
				{
					tokens = append(tokens, NewTokenWord(string(tokenValue)))
					tokenValue = []rune{}
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
			case redirectRuneType:
				{
					tokens = append(tokens, NewTokenWord(string(tokenValue)))
					tokenValue = []rune{}
					state = inRedirectSymb
				}
			default:
				{
					tokenValue = append(tokenValue, nextRune)
				}
			}
		case inNonEscapeQuoteState:
			switch nextRuneType {
			case eofRuneType:
				{
					tokens = append(tokens, NewTokenWord(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
					break loop
				}
			case nonEscapeQuoteRuneType:
				{
					state = inWordState
				}
			default:
				{
					tokenValue = append(tokenValue, nextRune)
				}
			}
		case inEscapeQuoteState:
			switch nextRuneType {
			case eofRuneType:
				{
					tokens = append(tokens, NewTokenWord(string(tokenValue)))
					tokenValue = []rune{}
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
					tokenValue = append(tokenValue, nextRune)
				}
			}
		case inEscapeState:
			switch nextRuneType {
			case eofRuneType:
				{
					tokens = append(tokens, NewTokenWord(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
					break loop
				}
			default:
				tokenValue = append(tokenValue, nextRune)
				state = inWordState
			}
		case inEscapeQuotedState:
			switch nextRuneType {
			case eofRuneType:
				{
					tokens = append(tokens, NewTokenWord(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
					break loop
				}
			default:
				tokenValue = append(tokenValue, nextRune)
				state = inEscapeQuoteState
			}
		case inRedirectSymb:
			switch nextRuneType {
			case spaceRuneType:
				{
				}
			case eofRuneType:
				{
					return nil, fmt.Errorf("expect filename after redirect symbol")
				}
			case redirectRuneType:
				{
					state = inAppendRedirect
				}
			default:
				tokenValue = append(tokenValue, nextRune)
				state = inRedirectReadName
			}
		case inRedirectReadName:
			switch nextRuneType {
			case spaceRuneType:
				{
					tokens = append(tokens, NewTokenRedirect(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
				}
			case eofRuneType:
				{
					tokens = append(tokens, NewTokenRedirect(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
					break loop
				}
			default:
				{
					tokenValue = append(tokenValue, nextRune)
				}
			}
		case inOneState:
			switch nextRuneType {
			case spaceRuneType:
				{
					tokens = append(tokens, NewTokenRedirect(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
				}
			case eofRuneType:
				{
					tokens = append(tokens, NewTokenRedirect(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
					break loop
				}
			case redirectRuneType:
				{
					state = inRedirectSymb
					tokenValue = []rune{}
				}
			default:
				{
					tokenValue = append(tokenValue, nextRune)
					state = inWordState
				}
			}
		case inTwo:
			switch nextRuneType {
			case spaceRuneType:
				{
					tokens = append(tokens, NewTokenRedirect(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
				}
			case eofRuneType:
				{
					tokens = append(tokens, NewTokenRedirect(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
					break loop
				}
			case redirectRuneType:
				{
					state = inRedirectErrSymb
					tokenValue = []rune{}
				}
			default:
				{
					tokenValue = append(tokenValue, nextRune)
					state = inWordState
				}
			}
		case inRedirectErrSymb:
			switch nextRuneType {
			case spaceRuneType:
				{
				}
			case eofRuneType:
				{
					return nil, fmt.Errorf("expect filename after redirect symbol")
				}
			case redirectRuneType:
				{
					state = inAppendRedirectErr
				}
			default:
				tokenValue = append(tokenValue, nextRune)
				state = inRedirectErrReadName
			}
		case inRedirectErrReadName:
			switch nextRuneType {
			case spaceRuneType:
				{
					tokens = append(tokens, NewTokenRedirectErr(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
				}
			case eofRuneType:
				{
					tokens = append(tokens, NewTokenRedirectErr(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
					break loop
				}
			default:
				{
					tokenValue = append(tokenValue, nextRune)
				}
			}
		case inAppendRedirect:
			switch nextRuneType {
			case spaceRuneType:
				{
				}
			case eofRuneType:
				{
					return nil, fmt.Errorf("expect filename after redirect symbol")
				}
			default:
				tokenValue = append(tokenValue, nextRune)
				state = inAppendRedirectReadName
			}
		case inAppendRedirectReadName:
			switch nextRuneType {
			case spaceRuneType:
				{
					tokens = append(tokens, NewTokenAppendRedirect(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
				}
			case eofRuneType:
				{
					tokens = append(tokens, NewTokenAppendRedirect(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
					break loop
				}
			default:
				{
					tokenValue = append(tokenValue, nextRune)
				}
			}
		case inAppendRedirectErr:
			switch nextRuneType {
			case spaceRuneType:
				{
				}
			case eofRuneType:
				{
					return nil, fmt.Errorf("expect filename after redirect symbol")
				}
			default:
				tokenValue = append(tokenValue, nextRune)
				state = inAppendRedirectErrReadName
			}
		case inAppendRedirectErrReadName:
			switch nextRuneType {
			case spaceRuneType:
				{
					tokens = append(tokens, NewTokenAppendRedirectErr(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
				}
			case eofRuneType:
				{
					tokens = append(tokens, NewTokenAppendRedirectErr(string(tokenValue)))
					tokenValue = []rune{}
					state = startState
					break loop
				}
			default:
				{
					tokenValue = append(tokenValue, nextRune)
				}
			}
		}

	}

	return tokens, nil
}
