package base

import (
	gingrpc "github.com/dan-and-dna/gin-grpc"
	"github.com/dan-and-dna/gin-grpc-network/modules/network"
	"google.golang.org/grpc"
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

func (baseCtrl *BaseCtrl) ListenProto(pkg, service, method string, desc *grpc.ServiceDesc, listener func(grpc.ServerStream) error) {
	network.ListenProto(pkg, service, method, desc, listener)
}

func (baseCtrl *BaseCtrl) StopListenProto(pkg, service, method string) {
	network.StopListenProto(pkg, service, method)
}
