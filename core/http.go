package core

import (
	"github.com/dan-and-dna/singleinstmodule"
)

type HttpCore struct {
	singleinstmodule.SingleInstModuleCore `mapstructure:"-"`

	Enable                  bool   `mapstructure:"enable"`                    // 是否启动模块
	ListenIp                string `mapstructure:"listen_ip"`                 // http 监听ip
	ListenPort              int    `mapstructure:"listen_port"`               // http 监听端口
	UseMiddleware           bool   `mapstructure:"use_middleware"`            // 是否使用中间件
	UseMiddlewareRecovery   bool   `mapstructure:"use_middleware_recovery"`   // 是否使用崩溃恢复
	UseMiddlewareGinLogger  bool   `mapstructure:"use_middleware_ginlogger"`  // 是否使用gin日志
	UseMiddlewareZapLogger  bool   `mapstructure:"use_middleware_zaplogger"`  // 是否使用zap日志
	UseMiddlewareHttpHeader bool   `mapstructure:"use_middleware_httpheader"` // 是否使用http头
	UseMiddlewareRoute      bool   `mapstructure:"use_middleware_route"`      // 是否路由给其他模块
	ReadTimeOut             int    `mapstructure:"read_timeout"`              // 读超时
	WriteTimeOut            int    `mapstructure:"write_timeout"`             // 写超时
}
