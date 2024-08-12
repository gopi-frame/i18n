package i18n

import (
	"fmt"
	"github.com/gopi-frame/contract/translator"
	"github.com/gopi-frame/exception"
	"io"
	"net/http"
)

type loaderFunc func() ([]byte, error)

func (l loaderFunc) Load() ([]byte, error) {
	return l()
}

// LoaderFunc returns a loader from a function.
func LoaderFunc(fn func() ([]byte, error)) translator.Loader {
	return loaderFunc(fn)
}

// RemoteLoader is the default loader to load messages from a remote url.
// If the response status code is not 200, an error is returned.
type RemoteLoader struct {
	Req    *http.Request
	Client *http.Client
}

// Load loads messages from a remote url.
func (r *RemoteLoader) Load() ([]byte, error) {
	resp, err := r.Client.Do(r.Req)
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
}
