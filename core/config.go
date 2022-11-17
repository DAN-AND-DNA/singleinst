package core

import "github.com/dan-and-dna/singleinstmodule"

type ConfigCore struct {
	singleinstmodule.SingleInstModuleCore `mapstructure:"-"`

	EnableGrpc bool   `mapstructure:"enable_grpc"` // 是否启动模块（走配置）
	EnableHttp bool   `mapstructure:"enable_http"` // 是否启动模块（走配置）
	ListenIp   string `mapstructure:"listen_ip"`   // http 监听ip
	ListenPort int    `mapstructure:"listen_port"` // http 监听端口

	GRPCUseMiddlewareTags    bool `mapstructure:"grpc_use_middleware_tags"`     // 是否使用tags
	GRPCUseMiddlewareZap     bool `mapstructure:"grpc_use_middleware_zap"`      // 是否使用zap日志
	GRPCUseMiddlewareTraceId bool `mapstructure:"grpc_use_middleware_trace_id"` // 是否添加追踪id

	HttpUseMiddlewareRecovery   bool `mapstructure:"http_use_middleware_recovery"`   // 是否使用崩溃恢复
	HttpUseMiddlewareGinLogger  bool `mapstructure:"http_use_middleware_ginlogger"`  // 是否使用gin日志
	HttpUseMiddlewareZapLogger  bool `mapstructure:"http_use_middleware_zaplogger"`  // 是否使用zap日志
	HttpUseMiddlewareHttpHeader bool `mapstructure:"http_use_middleware_httpheader"` // 是否使用http头
	HttpReadTimeOut             int  `mapstructure:"http_read_timeout"`              // 读超时
	HttpWriteTimeOut            int  `mapstructure:"http_write_timeout"`             // 写超时
}
