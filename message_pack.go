package i18n

import (
	"github.com/gopi-frame/contract/translator"
	"golang.org/x/text/language"
)

// MessagePack returns a [translator.MessagePack] from the given messages and language tag.
func MessagePack(messages []translator.Message, tag language.Tag) translator.MessagePack {
	return &messagePack{
		Messages:    messages,
		LanguageTag: tag,
	}
}

type messagePack struct {
	Messages    []translator.Message
	LanguageTag language.Tag
}

func (m *messagePack) GetMessages() []translator.Message {
	return m.Messages
}

func (m *messagePack) GetLanguageTag() language.Tag {
	return m.LanguageTag
}
