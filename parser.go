package i18n

import "github.com/gopi-frame/contract/translator"

type parserFunc func([]byte) (translator.MessagePack, error)

func (p parserFunc) Parse(data []byte) (translator.MessagePack, error) {
	return p(data)
}

// ParserFunc returns a parser from a function.
func ParserFunc(fn func([]byte) (translator.MessagePack, error)) translator.Parser {
	return parserFunc(fn)
}
