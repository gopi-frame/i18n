package i18n

import (
	"fmt"
	"github.com/gopi-frame/contract/translator"
	"github.com/gopi-frame/exception"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io"
	"io/fs"
	"net/http"
)

// I18n is a wrapper around [i18n.Bundle] and an implementation of [translator.Translator].
type I18n struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
}

// New creates a new i18n instance with the given default language.
func New(defaultLanguage string) (*I18n, error) {
	languageTag, err := language.Parse(defaultLanguage)
	if err != nil {
		return nil, err
	}
	i := new(I18n)
	i.bundle = i18n.NewBundle(languageTag)
	i.localizer = i18n.NewLocalizer(i.bundle, defaultLanguage)
	return i, nil
}

// T returns the translation for the given id.
func (i *I18n) T(id string, data ...any) string {
	lc := &i18n.LocalizeConfig{
		MessageID: id,
	}
	if len(data) > 1 {
		lc.TemplateData = data
	} else if len(data) == 1 {
		lc.TemplateData = data[0]
	}
	return i.localizer.MustLocalize(lc)
}

// P returns the translation for the given id and plural count.
func (i *I18n) P(id string, pluralCount any, data ...any) string {
	lc := &i18n.LocalizeConfig{
		MessageID:   id,
		PluralCount: pluralCount,
	}
	if len(data) > 1 {
		lc.TemplateData = data
	} else if len(data) == 1 {
		lc.TemplateData = data[0]
	}
	return i.localizer.MustLocalize(lc)
}

// M returns the translation for the given [translator.Message].
func (i *I18n) M(message translator.Message) string {
	return i.localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          message.GetID(),
			Hash:        message.GetHash(),
			Description: message.GetDescription(),
			LeftDelim:   message.GetLeftDelim(),
			RightDelim:  message.GetRightDelim(),
			Zero:        message.GetZero(),
			One:         message.GetOne(),
			Two:         message.GetTwo(),
			Few:         message.GetFew(),
			Many:        message.GetMany(),
			Other:       message.GetOther(),
		},
	})
}

// Locale returns a translator for the given languages.
func (i *I18n) Locale(languages ...string) translator.Translator {
	l := new(I18n)
	l.bundle = i.bundle
	l.localizer = i18n.NewLocalizer(i.bundle, languages...)
	return l
}

// AddMessages adds messages to the bundle.
func (i *I18n) AddMessages(l string, messages ...translator.Message) error {
	languageTag, err := language.Parse(l)
	if err != nil {
		return err
	}
	return i.AddMessagesByLanguageTag(languageTag, messages...)
}

// AddMessagesByLanguageTag adds messages to the bundle by language tag.
func (i *I18n) AddMessagesByLanguageTag(languageTag language.Tag, messages ...translator.Message) error {
	var msgList []*i18n.Message
	for _, message := range messages {
		msgList = append(msgList, &i18n.Message{
			ID:          message.GetID(),
			Hash:        message.GetHash(),
			Description: message.GetDescription(),
			LeftDelim:   message.GetLeftDelim(),
			RightDelim:  message.GetRightDelim(),
			Zero:        message.GetZero(),
			One:         message.GetOne(),
			Two:         message.GetTwo(),
			Few:         message.GetFew(),
			Many:        message.GetMany(),
			Other:       message.GetOther(),
		})
	}
	return i.bundle.AddMessages(languageTag, msgList...)
}

// RegisterUnmarshalFunc registers a custom unmarshal function for the given format.
func (i *I18n) RegisterUnmarshalFunc(format string, unmarshaller func(data []byte, v any) error) {
	i.bundle.RegisterUnmarshalFunc(format, unmarshaller)
}

// LoadMessage loads messages from the given loader and parser.
func (i *I18n) LoadMessage(loader translator.Loader, parser translator.Parser) error {
	content, err := loader.Load()
	if err != nil {
		return err
	}
	messagePack, err := parser.Parse(content)
	if err != nil {
		return err
	}
	return i.AddMessagesByLanguageTag(messagePack.GetLanguageTag(), messagePack.GetMessages()...)
}

// LoadMessageFile loads messages from a file.
func (i *I18n) LoadMessageFile(path string) error {
	_, err := i.bundle.LoadMessageFile(path)
	return err
}

// LoadMessageFileFS loads messages from a file from the given file system.
func (i *I18n) LoadMessageFileFS(fsys fs.FS, path string) error {
	_, err := i.bundle.LoadMessageFileFS(fsys, path)
	return err
}

// LoadMessageRemote loads messages from a remote url.
func (i *I18n) LoadMessageRemote(remote string, parser translator.Parser) error {
	req, err := http.NewRequest(http.MethodGet, remote, nil)
	if err != nil {
		return err
	}
	return i.LoadMessageRemoteRequest(req, parser)
}

// LoadMessageRemoteRequest loads messages from a custom request.
func (i *I18n) LoadMessageRemoteRequest(req *http.Request, parser translator.Parser, clientOpts ...func(client *http.Client) error) error {
	client := new(http.Client)
	for _, clientOpt := range clientOpts {
		if err := clientOpt(client); err != nil {
			return err
		}
	}
	return i.LoadMessage(LoaderFunc(func() ([]byte, error) {
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, exception.New(fmt.Sprintf("invalid status code responsed: %d", resp.StatusCode))
		}
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				panic(err)
			}
		}()
		return content, nil
	}), parser)
}

// LanguageTags returns the list of language tags of the bundle.
func (i *I18n) LanguageTags() []language.Tag {
	return i.bundle.LanguageTags()
}
