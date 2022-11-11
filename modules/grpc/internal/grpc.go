package internal

import (
	"context"
	"fmt"
	"github.com/dan-and-dna/singleinstmodule"
	"google.golang.org/grpc"
	"log"
	"singleinst/cmd"
	"singleinst/core"
	"singleinst/modules/network"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	singleInst *Grpc = nil
	once       sync.Once
)

type Grpc struct {
	core           *core.GrpcCore
	cfgMgr         *viper.Viper
	zapLogger      *zap.Logger
	traceBaseId    uint64
	coreChanged    atomic.Bool // 配置是否更新
	isModuleLoaded bool        // 模块是否加载
	isModuleRun    bool        // 模块是否运行
}

func (grpc *Grpc) ModuleConstruct() {
	//加载配置文件
	grpc.cfgMgr = viper.New()
	grpc.cfgMgr.SetConfigName("grpc")
	grpc.cfgMgr.SetConfigType("json")
	grpc.cfgMgr.AddConfigPath("./config/modules/")
	grpc.cfgMgr.AddConfigPath(cmd.ConfigPath)
	err := grpc.cfgMgr.ReadInConfig()
	if err != nil {
		panic(err)
	}
	grpc.coreChanged.Store(false)
	grpc.core = new(core.GrpcCore)

	grpc.core.Lock()
	defer grpc.core.Unlock()
	if err := grpc.cfgMgr.Unmarshal(grpc.core); err != nil {
		panic(err)
	}

	// FIXME 测试用的追踪ID，需要比较好唯一ID
	grpc.traceBaseId = 10000
}

func (grpc *Grpc) ModuleDestruct() {
}

func (grpc *Grpc) ModuleLock() singleinstmodule.ModuleCore {
	grpc.core.Lock()
	return grpc.core
}

func (grpc *Grpc) ModuleUnlock() {
	grpc.core.Unlock()
	grpc.CoreChanged()
}

func (grpc *Grpc) ModuleShutdown() {
	if grpc.zapLogger != nil {
		grpc.zapLogger.Sync()
	}
	log.Println("[grpc] shutdown")
}

func (grpc *Grpc) ModuleAfterRun(method string) {
	log.Printf("[grpc] %s\n", method)
	time.Sleep(50 * time.Millisecond)
}

func (grpc *Grpc) ModuleRunConfigWatcher() {
	// 监听配置变化
	grpc.cfgMgr.OnConfigChange(func(e fsnotify.Event) {
		grpc.core.Lock()
		defer grpc.core.Unlock()

		err := grpc.cfgMgr.Unmarshal(grpc.core)
		if err != nil {
			panic(err)
		}

		grpc.CoreChanged()
	})

	grpc.cfgMgr.WatchConfig()
}

func (grpc *Grpc) ModuleRunStartup() {
	grpc.CoreChanged()
}

func (grpc *Grpc) ModuleRestart() bool {
	if grpc.coreChanged.CompareAndSwap(true, false) {
		log.Println("[grpc] start restart")
		grpc.Recreate()
		return true
	}

	return false
}

func (grpc *Grpc) CoreChanged() {
	grpc.coreChanged.Store(true)
	singleinstmodule.RestartModule(grpc)
}

func (grpc *Grpc) Recreate() {
	cfg := network.ModuleLock().(*core.NetworkCore)
	defer network.ModuleUnlock()

	grpc.core.RLock()
	defer grpc.core.RUnlock()

	cfg.GrpcMiddlewares = nil

	if grpc.core.UseMiddleware {
		// 启动tags
		if grpc.core.UseMiddlewareTags {
			cfg.GrpcMiddlewares = append(cfg.GrpcMiddlewares, grpc_ctxtags.UnaryServerInterceptor())
		}

		// 启动zap
		if grpc.core.UseMiddlewareZap {
			opts := []grpc_zap.Option{
				grpc_zap.WithDurationField(grpc_zap.DefaultDurationToField),
			}
			var err error
			grpc.zapLogger, err = zap.NewProduction()
			if err != nil {
				panic(err)
			}

			grpc_zap.ReplaceGrpcLoggerV2(grpc.zapLogger)
			cfg.GrpcMiddlewares = append(cfg.GrpcMiddlewares, grpc_zap.UnaryServerInterceptor(grpc.zapLogger, opts...))
		}

		// 追踪Id
		if grpc.core.UseMiddlewareTraceId {
			cfg.GrpcMiddlewares = append(cfg.GrpcMiddlewares, grpc.appendTraceId())
		}
	}

	cfg.ListenGrpc = grpc.core.Enable
	if cfg.ListenGrpc {
		cfg.ListenIp = grpc.core.ListenIp
		cfg.ListenPort = grpc.core.ListenPort
	}

}

func (grpc1 *Grpc) appendTraceId() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		grpc1.traceBaseId++
		grpc_ctxtags.Extract(ctx).Set("trace_id", fmt.Sprintf("trace_%d", grpc1.traceBaseId))
		return handler(ctx, req)
	}
}

func GetSingleInst() *Grpc {
	if singleInst == nil {
		once.Do(func() {
			singleInst = new(Grpc)
		})
	}

	return singleInst
}

// 注册模块
func init() {
	singleinstmodule.Register(GetSingleInst())
}
