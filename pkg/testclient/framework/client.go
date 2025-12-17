package framework

import (
	"context"
	"net/http"
)

type Client interface {
	GetContext() context.Context
	Init() error
	Start() error
	Quit() error

	AddLogic(logic Logic) error
	GetLogics() []Logic

	AddCommand(cmd ClientCommand) error

	// http
	GetUrl() string
	Do(req *http.Request) (*http.Response, error)
}
