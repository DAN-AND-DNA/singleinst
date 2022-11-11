package core

import (
	"context"
	"github.com/dan-and-dna/singleinstmodule"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type NetworkCore struct {
	singleinstmodule.SingleInstModuleCore `mapstructure:"-"`

	// 公共
	Enable     bool   `mapstructure:"enable"` // 是否启动模块
	ListenHttp bool   `mapstructure:"-"`      // 是否是http服务
	ListenGrpc bool   `mapstructure:"-"`      // 是否是grpc服务
	ListenIp   string `mapstructure:"-"`      // 监听ip
	ListenPort int    `mapstructure:"-"`      // 监听port

	// http
	HttpReadTimeOut     int                       `mapstructure:"-"` // http服务读超时
	HttpWriteTimeOut    int                       `mapstructure:"-"` // http服务写超时
	HttpMiddlewares     []gin.HandlerFunc         `mapstructure:"-"` // http中间件
	HttpMiddlewareRoute gin.HandlerFunc           `mapstructure:"-"` // http中间件路由
	PathToServiceName   func(*gin.Context) string // http路径转grpc的服务名

	// grpc
	GrpcMiddlewares []grpc.UnaryServerInterceptor `mapstructure:"-"` // grpc中间件

}

type Network interface {
	// 监听消息
	ListenProto(pkg, service, method string, listener func(context.Context, interface{}))
	// 处理消息
	HandleProto(pkg, service, method string, handler func(context.Context, interface{}) (interface{}, error))
	// 停止监听消息
	StopListenProto(pkg, service, method string)
	// 停止处理消息
	StopHandleProto(pkg, service, method string)
	// 获得当前配置
	GetConfig() NetworkCore
	// 更新配置
	UpdateCfg(cfg NetworkCore) error
	// 通知协议给监听者
	NotifyListeners(ctx context.Context, req interface{}, key string)
	// 通知协议给处理者
	NotifyHandler(ctx context.Context, req interface{}, key string) (interface{}, error)
	// 重新启动
	ReStart()
}

type ProtoListener interface {
	Listen(context.Context, interface{})
}

type ProtoHandler interface {
	Handle(context.Context, interface{}) (interface{}, error)
}
