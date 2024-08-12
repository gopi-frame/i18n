# Overview
[![Go Reference](https://pkg.go.dev/badge/github.com/gopi-frame/i18n.svg)](https://pkg.go.dev/github.com/gopi-frame/i18n)
[![Go](https://github.com/gopi-frame/i18n/actions/workflows/go.yml/badge.svg)](https://github.com/gopi-frame/i18n/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopi-frame/i18n)](https://goreportcard.com/report/github.com/gopi-frame/i18n)
[![codecov](https://codecov.io/gh/gopi-frame/i18n/graph/badge.svg?token=AQ4qBviH5M)](https://codecov.io/gh/gopi-frame/i18n)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

Package i18n provides an implementation of the i18n based on the [go-i18n](https://github.com/nicksnyder/go-i18n)
package.

# Installation

```shell
go get -u -v github.com/gopi-frame/i18n
```

# Import

```go
import "github.com/gopi-frame/i18n"
```

# Usage

```go
package main

import (
    "github.com/gopi-frame/i18n"
    "golang.org/x/text/language"
)
import i18nlib "github.com/nicksnyder/go-i18n/v2/i18n"

func main() {
    // create a new i18n instance with the default locale (en-US)
    i, err := i18n.New("en-US")
    if err != nil {
        panic(err)
    }
    // add a message to the specific locale
    err = i.AddMessages("en-US", i18n.Message(&i18nlib.Message{ID: "hello", Other: "Hello, world!"}))
    // or add a message to the specific locale by language tag
    // err = i.AddMessagesByLanguageTag(language.MustParse("en-US"), i18n.Message(&i18nlib.Message{ID: "hello", Other: "Hello, world!"}))
    if err!= nil {
        panic(err)
    }
    // translate message
    msg := i.T("hello")
    fmt.Println(msg) // Hello, world!
}
```

# Load messages

```go
// implement a loader
func Load() ([]byte, error) {
    // load content from somewhere
}

// implement a parser
func Parse(data []byte) (translator.MessagePack, error) {
    // parse data to message pack
}

// load messages from custom loader and parser
err := i.LoadMessage(i18n.LoaderFunc(Load), i18n.ParserFunc(Parse))
if err != nil {
    panic(err)
}
```

# Load messages from file

The file name should be ended with `.{locale}.{format}`, for example: `locale.en-US.json`.
The `{locale}` is the locale code, `{format}` is the format of the messages file.

```go
// load messages from json file
err := i.LoadMessageFile("locale.en-US.json")
if err != nil {
    panic(err)
}

// load message from other format file
// register unmarshal function for the format first
i.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
err := i.LoadMessageFile("locale.en-US.yaml")
```

# Load messages from file system

The file name should be ended with `.{locale}.{format}`, for example: `locale.en-US.json`.
The `{locale}` is the locale code, `{format}` is the format of the messages file.

```go
// load messages from json file system
err := i.LoadMessageFS(http.Dir("locales"), "locale.en-US.json")
if err != nil {
    panic(err)
}
```

# Load messages from remote url

```go
// load messages from remote url
err := i.LoadMessageRemote("https://example.com/locale.en-US.json", i18n.ParserFunc(func(data []byte) (translator.MessagePack, error) {
    // parse data to message pack
}))
if err != nil {
    panic(err)
}

// load messages from remote url with custom http client and request
err := i.LoadMessageRemoteRequest(req, i18n.ParserFunc(func(data []byte) (translator.MessagePack, error) {
    // parse data to message pack
}, func(client *http.Client) error {
    // customize http client
})
if err!= nil {
    panic(err)
}
```
