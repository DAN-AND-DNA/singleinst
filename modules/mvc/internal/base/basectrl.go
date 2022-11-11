package base

import (
	"context"
	gingrpc "github.com/dan-and-dna/gin-grpc"
	"google.golang.org/grpc"
	"singleinst/modules/network"
)

type BaseCtrl struct {
}

func (baseCtrl *BaseCtrl) Initialize() error {
	baseCtrl.reset()
	return nil
}

func (baseCtrl *BaseCtrl) JustAsController() {
}

func (baseCtrl *BaseCtrl) OnRun() {

}

func (baseCtrl *BaseCtrl) OnStop() {
}

func (baseCtrl *BaseCtrl) OnClean() {
	baseCtrl.reset()
}

func (baseCtrl *BaseCtrl) reset() {

}

func (baseCtrl *BaseCtrl) HandleProto(pkg, service, method string, desc *grpc.ServiceDesc, handler gingrpc.Handler) {
	network.HandleProto(pkg, service, method, desc, handler)

}

func (baseCtrl *BaseCtrl) StopHandleProto(pkg, service, method string) {
	network.StopHandleProto(pkg, service, method)

}

func (baseCtrl *BaseCtrl) ListenProto(pkg, service, method string, listener func(context.Context, interface{})) {
	network.ListenProto(pkg, service, method, listener)
}

func (baseCtrl *BaseCtrl) StopListenProto(pkg, service, method string) {
	network.StopListenProto(pkg, service, method)
}
