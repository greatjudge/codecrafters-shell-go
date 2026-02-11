package tokenizer

type TokenType int

const (
	WordTokenType           TokenType = 1
	RedirectTokenType       TokenType = 2
	RedirectErrTokenType    TokenType = 3
	AppendRedirectTokenType TokenType = 4
)

type Token struct {
	Type  TokenType
	Value string
}

func NewTokenWord(value string) Token {
	return Token{
		Type:  WordTokenType,
		Value: value,
	}
}

func NewTokenRedirect(value string) Token {
	return Token{
		Type:  RedirectTokenType,
		Value: value,
	}
}

func NewTokenRedirectErr(value string) Token {
	return Token{
		Type:  RedirectErrTokenType,
		Value: value,
	}
}

func NewTokenAppendRedirect(value string) Token {
	return Token{
		Type:  AppendRedirectTokenType,
		Value: value,
	}
}
