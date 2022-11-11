package base

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type BaseModel struct {
	listeners map[string]([]func(context.Context))
	handlers  map[string]endpoint.Endpoint
}

func (coreModel *BaseModel) Notify(msg string, ctx context.Context) (interface{}, error) {
	if listeners, ok := coreModel.listeners[msg]; ok {
		for _, listener := range listeners {
			listener(ctx)
		}
	}

	if handler, ok := coreModel.handlers[msg]; ok {
		return handler(ctx, nil)
	}

	return nil, nil
}

func (coreModel *BaseModel) ListenMsg(msg string, listener func(context.Context)) {
	if coreModel.listeners == nil {
		coreModel.listeners = make(map[string][]func(context.Context))
	}
	coreModel.listeners[msg] = append(coreModel.listeners[msg], listener)
}

func (coreModel *BaseModel) StopListenMsg(msg string) {
	delete(coreModel.listeners, msg)
}

func (coreModel *BaseModel) HandleMsg(msg string, handler endpoint.Endpoint) {
	if coreModel.handlers == nil {
		coreModel.handlers = make(map[string]endpoint.Endpoint)
	}

	coreModel.handlers[msg] = handler
}

func (coreModel *BaseModel) StopHandleMsg(msg string) {
	delete(coreModel.handlers, msg)
}

func (coreModel *BaseModel) Initialize() error {
	coreModel.reset()
	return nil
}

func (coreModel *BaseModel) OnStop() {
}

func (coreModel *BaseModel) OnClean() {
	coreModel.reset()
}

func (coreModel *BaseModel) OnRun() {
}

func (coreModel *BaseModel) JustAsModel() {
}

func (coreModel *BaseModel) reset() {
	coreModel.listeners = nil
	coreModel.handlers = nil
}
