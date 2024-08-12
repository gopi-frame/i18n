package i18n

import (
	"github.com/gopi-frame/contract/translator"
)

type loaderFunc func() ([]byte, error)

func (l loaderFunc) Load() ([]byte, error) {
	return l()
}

// LoaderFunc returns a loader from a function.
func LoaderFunc(fn func() ([]byte, error)) translator.Loader {
	return loaderFunc(fn)
}
