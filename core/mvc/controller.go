package mvc

import (
	"context"
	gingrpc "github.com/dan-and-dna/gin-grpc"
	"google.golang.org/grpc"
)

type Controller interface {
	Initialize() error
	OnRun()
	OnStop()
	OnClean()
	HandleProto(pkg string, service string, method string, desc *grpc.ServiceDesc, handler gingrpc.Handler)
	StopHandleProto(pkg, service, method string)
	ListenProto(pkg, service, method string, listener func(context.Context, interface{}))
	StopListenProto(pkg, service, method string)
	JustAsController()
}
