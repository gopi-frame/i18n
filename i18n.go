package i18n

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"github.com/gopi-frame/collection/kv"

	"github.com/gopi-frame/contract/translator"
	"github.com/gopi-frame/exception"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var defaultMessages = *kv.NewMap[string, string]()

func SetDefaultMessage(id, message string) {
	defaultMessages.Lock()
	defer defaultMessages.Unlock()
	defaultMessages.Set(id, message)
}

func GetDefaultMessage(id string) string {
	defaultMessages.RLock()
	defer defaultMessages.RUnlock()
	message, _ := defaultMessages.Get(id)
	return message
}

func SetDefaultMessages(messages map[string]string) {
	defaultMessages.Lock()
	defer defaultMessages.Unlock()
	for id, message := range messages {
		defaultMessages.Set(id, message)
	}
}

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
// If the length of data is even, it will be used as key-value pairs.
// If the length of data is greater than 1 and is odd, it will be used as a slice.
func (i *I18n) T(id string, data ...any) string {
	lc := &i18n.LocalizeConfig{
		MessageID: id,
	}
	if len(data) > 1 {
		if len(data)%2 == 0 {
			var d = make(map[any]any, len(data)/2)
			for i := 0; i < len(data); i += 2 {
				d[data[i]] = data[i+1]
			}
			lc.TemplateData = d
		} else {
			lc.TemplateData = data
		}
	} else if len(data) == 1 {
		lc.TemplateData = data[0]
	}
	r, err := i.localizer.Localize(lc)
	if err != nil {
		if defaultMessage := GetDefaultMessage(id); defaultMessage != "" {
			return defaultMessage
		}
		return id
	}
	return r
}

// P returns the translation for the given id and plural count.
// If the length of data is even, it will be used as key-value pairs.
// If the length of data is greater than 1 and is odd, it will be used as a slice.
func (i *I18n) P(id string, pluralCount any, data ...any) string {
	lc := &i18n.LocalizeConfig{
		MessageID:   id,
		PluralCount: pluralCount,
	}
	if len(data) > 1 {
		if len(data)%2 == 0 {
			var d = make(map[any]any, len(data)/2)
			for i := 0; i < len(data); i += 2 {
				d[data[i]] = data[i+1]
			}
			lc.TemplateData = d
		} else {
			lc.TemplateData = data
		}
	} else if len(data) == 1 {
		lc.TemplateData = data[0]
	}
	r, err := i.localizer.Localize(lc)
	if err != nil {
		if defaultMessage := GetDefaultMessage(id); defaultMessage != "" {
			return defaultMessage
		}
		return id
	}
	return r
}

// M returns the translation for the given [translator.Message].
// If the length of data is even, it will be used as key-value pairs.
// If the length of data is greater than 1 and is odd, it will be used as a slice.
func (i *I18n) M(message translator.Message, pluralCount any, data ...any) string {
	lc := &i18n.LocalizeConfig{
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
		PluralCount: pluralCount,
	}
	if len(data) > 1 {
		if len(data)%2 == 0 {
			var d = make(map[any]any, len(data)/2)
			for i := 0; i < len(data); i += 2 {
				d[data[i]] = data[i+1]
			}
			lc.TemplateData = d
		} else {
			lc.TemplateData = data
		}
	} else if len(data) == 1 {
		lc.TemplateData = data[0]
	}
	if defaultMessage := GetDefaultMessage(message.GetID()); defaultMessage != "" {
		lc.DefaultMessage = &i18n.Message{
			ID:    message.GetID(),
			Other: defaultMessage,
		}
	} else {
		lc.DefaultMessage = &i18n.Message{
			ID:    message.GetID(),
			Other: message.GetID(),
		}
	}
	return i.localizer.MustLocalize(lc)
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
