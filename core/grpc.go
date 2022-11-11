package core

import (
	"github.com/dan-and-dna/singleinstmodule"
)

type GrpcCore struct {
	singleinstmodule.SingleInstModuleCore `mapstructure:"-"`

	Enable               bool   `mapstructure:"enable"`                  // 是否启动模块（走配置）
	ListenIp             string `mapstructure:"listen_ip"`               // http 监听ip
	ListenPort           int    `mapstructure:"listen_port"`             // http 监听端口
	UseMiddleware        bool   `mapstructure:"use_middleware"`          // 是否使用中间件
	UseMiddlewareTags    bool   `mapstructure:"use_middleware_tags"`     // 是否使用tags
	UseMiddlewareRoute   bool   `mapstructure:"use_middleware_route"`    // 是否路由给其他模块
	UseMiddlewareZap     bool   `mapstructure:"use_middleware_zap"`      // 是否使用zap日志
	UseMiddlewareTraceId bool   `mapstructure:"use_middleware_trace_id"` // 是否添加追踪id
}
