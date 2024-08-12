package i18n

import (
	"github.com/gopi-frame/contract/translator"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Message returns a [translator.Message] from the given [i18n.Message].
func Message(m *i18n.Message) translator.Message {
	return &message{
		Message: *m,
	}
}

type message struct {
	i18n.Message
}

func (m *message) GetID() string {
	return m.ID
}

func (m *message) GetHash() string {
	return m.Hash
}

func (m *message) GetDescription() string {
	return m.Description
}

func (m *message) GetLeftDelim() string {
	return m.LeftDelim
}

func (m *message) GetRightDelim() string {
	return m.RightDelim
}

func (m *message) GetZero() string {
	return m.Zero
}

func (m *message) GetOne() string {
	return m.One
}

func (m *message) GetTwo() string {
	return m.Two
}

func (m *message) GetFew() string {
	return m.Few
}

func (m *message) GetMany() string {
	return m.Many
}

func (m *message) GetOther() string {
	return m.Other
}
