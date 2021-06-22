package appctx

import "github.com/flared/lokify/pkg/loki"

type Context struct {
	Loki loki.Client
}

func New(loki loki.Client) *Context {
	return &Context{
		Loki: loki,
	}
}
