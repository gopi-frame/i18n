package i18n

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gopi-frame/contract/translator"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestI18n_T(t *testing.T) {
	t.Run("translate without data", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.AddMessages("en", Message(&i18n.Message{
				ID:          "test",
				Description: "test",
				Other:       "test one",
			}))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			message := i.T("test")
			assert.Equal(t, "test one", message)
		}
	})

	t.Run("translate with single data", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.AddMessages("en", Message(&i18n.Message{
				ID:          "greeting",
				Description: "greeting",
				Other:       "hello, {{.name}}",
			}))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			message := i.T("greeting", map[string]interface{}{"name": "world"})
			assert.Equal(t, "hello, world", message)
		}
	})

	t.Run("translate with multiple data", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.AddMessages("en", Message(&i18n.Message{
				ID:          "greeting",
				Description: "greeting",
				Other:       "hello, {{index . 0}} and {{index . 1}}",
			}))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			message := i.T("greeting", "world", "friend")
			assert.Equal(t, "hello, world and friend", message)
		}
	})
}

func TestI18n_P(t *testing.T) {
	t.Run("translate without data", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.AddMessages("en", Message(&i18n.Message{
				ID:          "test",
				Description: "test",
				One:         "test one",
				Other:       "test other",
			}))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			message := i.P("test", 1)
			assert.Equal(t, "test one", message)
			message = i.P("test", 2)
			assert.Equal(t, "test other", message)
		}
	})

	t.Run("translate with single data", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.AddMessages("en", Message(&i18n.Message{
				ID:          "greeting",
				Description: "greeting",
				One:         "hello, {{.name}}",
				Other:       "hello, {{.name}}2",
			}))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			message := i.P("greeting", 1, map[string]interface{}{"name": "world"})
			assert.Equal(t, "hello, world", message)
			message = i.P("greeting", 2, map[string]interface{}{"name": "world"})
			assert.Equal(t, "hello, world2", message)
		}
	})

	t.Run("translate with multiple data", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.AddMessages("en", Message(&i18n.Message{
				ID:          "greeting",
				Description: "greeting",
				One:         "hello, {{index . 0}} and {{index . 1}}",
				Other:       "hello, {{index . 0}} and {{index . 1}}2",
			}))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			message := i.P("greeting", 1, "world", "friend")
			assert.Equal(t, "hello, world and friend", message)
			message = i.P("greeting", 2, "world", "friend")
			assert.Equal(t, "hello, world and friend2", message)
		}
	})
}

func TestI18n_M(t *testing.T) {
	t.Run("translate without data", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			message := i.M(Message(&i18n.Message{
				ID:          "test",
				Description: "test",
				Other:       "test one",
			}), nil)
			assert.Equal(t, "test one", message)
		}
	})
}

func TestI18n_Locale(t *testing.T) {
	t.Run("translate without data", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.AddMessages("en", Message(&i18n.Message{
				ID:          "test",
				Description: "test",
				Other:       "test one",
			}))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			err = i.AddMessages("zh", Message(&i18n.Message{
				ID:          "test",
				Description: "test",
				Other:       "测试一",
			}))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			l := i.Locale("zh")
			message := l.T("test")
			assert.Equal(t, "测试一", message)
		}
	})

	t.Run("translate with single data", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.AddMessages("en", Message(&i18n.Message{
				ID:          "greeting",
				Description: "greeting",
				Other:       "hello, {{.name}}",
			}))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			err = i.AddMessages("zh", Message(&i18n.Message{
				ID:          "greeting",
				Description: "greeting",
				Other:       "你好，{{.name}}",
			}))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			l := i.Locale("zh")
			message := l.T("greeting", map[string]interface{}{"name": "world"})
			assert.Equal(t, "你好，world", message)
		}
	})

	t.Run("translate with multiple data", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.AddMessages("en", Message(&i18n.Message{
				ID:          "greeting",
				Description: "greeting",
				Other:       "hello, {{index . 0}} and {{index . 1}}",
			}))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			err = i.AddMessages("zh", Message(&i18n.Message{
				ID:          "greeting",
				Description: "greeting",
				Other:       "你好，{{index . 0}}和{{index . 1}}",
			}))
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			l := i.Locale("zh")
			message := l.T("greeting", "world", "friend")
			assert.Equal(t, "你好，world和friend", message)
		}
	})
}

func TestI18n_AddMessages(t *testing.T) {
	t.Run("invalid language", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.AddMessages("q", Message(&i18n.Message{
				ID:          "test",
				Description: "test",
				Other:       "test one",
			}))
			assert.Error(t, err)
		}
	})

	t.Run("valid language", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.AddMessages("en", Message(&i18n.Message{
				ID:          "test",
				Description: "test",
				Other:       "test one",
			}))
			assert.NoError(t, err)
			message := i.T("test")
			assert.Equal(t, "test one", message)
		}
	})
}

func TestI18n_LoadMessage(t *testing.T) {
	t.Run("load failed", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.LoadMessage(LoaderFunc(func() ([]byte, error) {
				return nil, errors.New("load failed")
			}), ParserFunc(func(data []byte) (translator.MessagePack, error) {
				return nil, nil
			}))
			assert.EqualError(t, err, "load failed")
		}
	})

	t.Run("parse failed", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.LoadMessage(LoaderFunc(func() ([]byte, error) {
				return []byte("test"), nil
			}), ParserFunc(func(data []byte) (translator.MessagePack, error) {
				return nil, errors.New("parse failed")
			}))
			assert.EqualError(t, err, "parse failed")
		}
	})

	t.Run("success", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.LoadMessage(LoaderFunc(func() ([]byte, error) {
				return []byte("test"), nil
			}), ParserFunc(func(data []byte) (translator.MessagePack, error) {
				return MessagePack([]translator.Message{
					Message(&i18n.Message{
						ID:          "test",
						Description: "test",
						Other:       "test one",
					}),
				}, language.English), nil
			}))
			assert.NoError(t, err)
			message := i.T("test")
			assert.Equal(t, "test one", message)
		}
	})
}

func TestI18n_LoadMessageFile(t *testing.T) {
	t.Run("load failed", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.LoadMessageFile("test.en.json")
			assert.ErrorIs(t, err, os.ErrNotExist)
		}
	})

	t.Run("load success", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.LoadMessageFile("testdata/test.en.json")
			assert.NoError(t, err)
			message := i.T("test")
			assert.Equal(t, "test one", message)
		}
	})
}

func TestI18n_LoadMessageFileFS(t *testing.T) {
	t.Run("load failed", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.LoadMessageFileFS(os.DirFS("test"), "test.en.json")
			assert.ErrorIs(t, err, os.ErrNotExist)
		}
	})

	t.Run("load success", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.LoadMessageFileFS(os.DirFS("testdata"), "test.en.json")
			assert.NoError(t, err)
			message := i.T("test")
			assert.Equal(t, "test one", message)
		}
	})
}

func TestI18n_LoadMessageRemote(t *testing.T) {
	t.Run("load success", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			srv := &http.Server{
				Addr: ":8080",
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					_, err := w.Write([]byte(`[{"id": "test", "description": "test", "other": "test one"}]`))
					if err != nil {
						assert.FailNow(t, err.Error())
					}
				}),
			}
			go func() {
				err := srv.ListenAndServe()
				if err != nil {
					if err != http.ErrServerClosed {
						assert.FailNow(t, err.Error())
					}
				}
			}()
			defer srv.Shutdown(context.Background())
			time.Sleep(time.Millisecond * 100)
			err = i.LoadMessageRemote("http://127.0.0.1:8080/test.en.json", ParserFunc(func(data []byte) (translator.MessagePack, error) {
				var items []map[string]any
				err := json.Unmarshal(data, &items)
				if err != nil {
					return nil, err
				}
				var messages []translator.Message
				for _, item := range items {
					messages = append(messages, Message(&i18n.Message{
						ID:          item["id"].(string),
						Description: item["description"].(string),
						Other:       item["other"].(string),
					}))
				}
				return MessagePack(messages, language.English), nil
			}))
			assert.Nil(t, err)
			assert.Equal(t, "test one", i.T("test"))
		}
	})

	t.Run("load failed", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			err = i.LoadMessageRemote("http://127.0.0.1:8080/test.en.json", ParserFunc(func(data []byte) (translator.MessagePack, error) {
				return nil, errors.New("parse failed")
			}))
			assert.Error(t, err)
		}
	})
}

func TestI18n_LoadMessageRemoteRequest(t *testing.T) {
	t.Run("load success", func(t *testing.T) {
		i, err := New("en")
		if err != nil {
			assert.FailNow(t, err.Error())
		} else {
			srv := &http.Server{
				Addr: ":8080",
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					_, err := w.Write([]byte(`[{"id": "test", "description": "test", "other": "test one"}]`))
					if err != nil {
						assert.FailNow(t, err.Error())
					}
				}),
			}
			go func() {
				err := srv.ListenAndServe()
				if err != nil {
					err := srv.ListenAndServe()
					if err != nil {
						if err != http.ErrServerClosed {
							assert.FailNow(t, err.Error())
						}
					}
				}
			}()
			defer srv.Shutdown(context.Background())
			time.Sleep(time.Millisecond * 100)
			req, err := http.NewRequest("GET", "http://127.0.0.1:8080/test.en.json", nil)
			if err != nil {
				assert.FailNow(t, err.Error())
			}
			err = i.LoadMessageRemoteRequest(req, ParserFunc(func(data []byte) (translator.MessagePack, error) {
				var items []map[string]any
				err := json.Unmarshal(data, &items)
				if err != nil {
					return nil, err
				}
				var messages []translator.Message
				for _, item := range items {
					messages = append(messages, Message(&i18n.Message{
						ID:          item["id"].(string),
						Description: item["description"].(string),
						Other:       item["other"].(string),
					}))
				}
				return MessagePack(messages, language.English), nil
			}))
			assert.Nil(t, err)
			assert.Equal(t, "test one", i.T("test"))
		}
	})
}
