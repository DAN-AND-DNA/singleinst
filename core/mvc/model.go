package mvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Model interface {
	Initialize() error
	OnRun()
	OnStop()
	OnClean()
	Notify(string, context.Context) (interface{}, error)
	ListenMsg(string, func(context.Context))
	StopListenMsg(string)
	HandleMsg(string, endpoint.Endpoint)
	StopHandleMsg(string)
	JustAsModel()
}
