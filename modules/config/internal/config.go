package internal

import (
	networkcore "github.com/dan-and-dna/gin-grpc-network/core"
	"github.com/dan-and-dna/gin-grpc-network/modules/grpc"
	"github.com/dan-and-dna/gin-grpc-network/modules/http"
	"github.com/dan-and-dna/gin-grpc-network/utils"
	"github.com/dan-and-dna/singleinstmodule"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"singleinst/cmd"
	"singleinst/core"
	"sync"
	"sync/atomic"
)

var (
	singleInst *Config = nil
	once       sync.Once
)

type Config struct {
	core        *core.ConfigCore
	coreMgr     *viper.Viper
	coreChanged atomic.Bool // 配置是否更新
	zapLogger   *zap.Logger
}

func (config *Config) ModuleConstruct() {
	config.coreMgr = viper.New()
	config.coreMgr.SetConfigName("gin-grpc-network")
	config.coreMgr.SetConfigType("json")
	config.coreMgr.AddConfigPath("./config/modules/")
	config.coreMgr.AddConfigPath(cmd.ConfigPath)
	err := config.coreMgr.ReadInConfig()
	if err != nil {
		panic(err)
	}
	config.coreChanged.Store(false)
	config.core = new(core.ConfigCore)
	config.zapLogger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}

	config.core.Lock()
	defer config.core.Unlock()

	if err := config.coreMgr.Unmarshal(config.core); err != nil {
		panic(err)
	}

	log.Println("[config] constructed")
}

func (config *Config) ModuleDestruct() {
	config.zapLogger.Sync()
	log.Println("[config] destructed")
}

func (config *Config) ModuleLock() singleinstmodule.ModuleCore {
	config.core.Lock()
	return config.core
}

func (config *Config) ModuleUnlock() {
	config.core.Unlock()
	config.CoreChanged()
}

func (config *Config) ModuleAfterRun(method string) {
	log.Printf("[config] %s\n", method)
}

func (config *Config) ModuleRunStartup() {
	grpc_zap.ReplaceGrpcLoggerV2(config.zapLogger)
	config.CoreChanged()
}

func (config *Config) ModuleRunConfigWatcher() {
	// 监听配置变化
	config.coreMgr.OnConfigChange(func(e fsnotify.Event) {
		config.core.Lock()
		defer config.core.Unlock()

		err := config.coreMgr.Unmarshal(config.core)
		if err != nil {
			panic(err)
		}

		config.CoreChanged()
	})

	config.coreMgr.WatchConfig()
}

func (config *Config) ModuleRestart() bool {
	if config.coreChanged.CompareAndSwap(true, false) {
		log.Println("[config] start restart")
		config.Recreate()
		return true
	}

	return false
}

func (config *Config) ModuleShutdown() {
	log.Println("[config] shutdown")
}

func (config *Config) Recreate() {
	// 防止本模块内容读写并行
	config.core.RLock()
	defer config.core.RUnlock()

	configCore := config.core

	if configCore.EnableGrpc {
		cfg := grpc.ModuleLock()
		defer grpc.ModuleUnlockRestart()

		// 只修改grpc模块，不会死锁
		grpcCore := cfg.(*networkcore.GrpcCore)
		grpcCore.Enable = true
		grpcCore.ListenPort = configCore.ListenPort
		grpcCore.ListenIp = configCore.ListenIp
		grpcCore.Middlewares = nil
		grpcCore.Middlewares = append(grpcCore.Middlewares,
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(config.zapLogger, grpc_zap.WithDurationField(grpc_zap.DefaultDurationToField)),
		)

		grpcCore.MiddlewaresStream = nil
		grpcCore.MiddlewaresStream = append(grpcCore.MiddlewaresStream,
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(config.zapLogger, grpc_zap.WithDurationField(grpc_zap.DefaultDurationToField)),
		)

	} else if configCore.EnableHttp {
		cfg := http.ModuleLock()
		defer http.ModuleUnlockRestart()

		// 只修改http模块，不会死锁
		httpCore := cfg.(*networkcore.HttpCore)
		httpCore.Enable = true
		httpCore.ListenPort = configCore.ListenPort
		httpCore.ListenIp = configCore.ListenIp

		httpCore.WriteTimeOut = configCore.HttpWriteTimeOut
		httpCore.ReadTimeOut = configCore.HttpReadTimeOut
		httpCore.Path = "sii/:pkg/:service/:method"
		httpCore.PathToServiceName = func(c *gin.Context) string {
			pkg := c.Param("pkg")
			service := c.Param("service")
			method := c.Param("method")
			return utils.MakeKey(pkg, service, method)
		}
		// gin中间件
		httpCore.Middlewares = nil
		httpCore.Middlewares = append(httpCore.Middlewares,
			gin.Recovery(),
			gin.Logger(),
		)
		// 跟gin的ctx无关，给每个handler使用
		httpCore.CtxOptions = nil
		httpCore.CtxOptions = append(httpCore.CtxOptions,
			&ZapOption{},
		)

	}

}

func (config *Config) CoreChanged() {
	config.coreChanged.Store(true)
	singleinstmodule.RestartModule(config)
}

func GetSingleInst() *Config {
	if singleInst == nil {
		once.Do(func() {
			singleInst = new(Config)
		})
	}

	return singleInst
}

func init() {
	singleinstmodule.Register(GetSingleInst())
}
