package command

import (
	"github.com/codecrafters-io/shell-starter-go/app/tokenizer"
)

type RedirectionWriteType int

const (
	RewriteRedirectionWriteType RedirectionWriteType = 1
	AppendRedirectionWriteType  RedirectionWriteType = 2
)

type StreamType int

const (
	StdoutStreamType StreamType = 1
	StderrStreamType StreamType = 2
)

var (
	redirectionTokenTypes = map[tokenizer.TokenType]bool{
		tokenizer.RedirectTokenType:          true,
		tokenizer.AppendRedirectTokenType:    true,
		tokenizer.RedirectErrTokenType:       true,
		tokenizer.AppendRedirectErrTokenType: true,
	}
)

func checkTokenTypeIsRedirection(t tokenizer.TokenType) bool {
	return redirectionTokenTypes[t]
}

type Redirection struct {
	WriteType  RedirectionWriteType
	StreamType StreamType
	Filename   string
}

func (r Redirection) IsEmpty() bool {
	return r == Redirection{}
}

func NewRedirectionFromToken(t tokenizer.Token) Redirection {
	switch t.Type {
	case tokenizer.RedirectTokenType:
		return Redirection{
			WriteType:  RewriteRedirectionWriteType,
			StreamType: StdoutStreamType,
			Filename:   t.Value,
		}
	case tokenizer.AppendRedirectTokenType:
		return Redirection{
			WriteType:  AppendRedirectionWriteType,
			StreamType: StdoutStreamType,
			Filename:   t.Value,
		}
	case tokenizer.RedirectErrTokenType:
		return Redirection{
			WriteType:  RewriteRedirectionWriteType,
			StreamType: StderrStreamType,
			Filename:   t.Value,
		}
	case tokenizer.AppendRedirectErrTokenType:
		return Redirection{
			WriteType:  AppendRedirectionWriteType,
			StreamType: StderrStreamType,
			Filename:   t.Value,
		}
	default:
		return Redirection{}
	}
}
