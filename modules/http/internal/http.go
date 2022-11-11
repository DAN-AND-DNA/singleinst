package internal

import (
	"github.com/dan-and-dna/singleinstmodule"
	"log"
	"singleinst/cmd"
	"singleinst/core"
	"singleinst/modules/network"
	"singleinst/utils"
	"sync"
	"sync/atomic"
	"time"

	//"github.com/axiaoxin-com/logging"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	singleInst *Http = nil
	once       sync.Once
)

type Http struct {
	core           *core.HttpCore
	cfgMgr         *viper.Viper
	coreChanged    atomic.Bool // core是否更新
	isModuleLoaded bool        // 模块是否加载
	isModuleRun    bool        // 模块是否运行
}

func (http *Http) ModuleConstruct() {
	//加载配置文件
	http.cfgMgr = viper.New()
	http.cfgMgr.SetConfigName("http")
	http.cfgMgr.SetConfigType("json")
	http.cfgMgr.AddConfigPath("./config/modules/")
	http.cfgMgr.AddConfigPath(cmd.ConfigPath)
	if err := http.cfgMgr.ReadInConfig(); err != nil {
		panic(err)
	}

	http.coreChanged.Store(false)
	http.core = new(core.HttpCore)
	http.core.Lock()
	defer http.core.Unlock()

	if err := http.cfgMgr.Unmarshal(http.core); err != nil {
		panic(err)
	}

	log.Println("[http] constructed")
}

func (http *Http) ModuleAfterRun(method string) {
	log.Printf("[http] %s\n", method)
}

func (http *Http) ModuleRunConfigWatcher() {
	// 监听配置变化
	http.cfgMgr.OnConfigChange(func(e fsnotify.Event) {
		http.core.Lock()

		err := http.cfgMgr.Unmarshal(http.core)
		if err != nil {
			panic(err)
		}
		http.core.Unlock()

		http.CoreChanged()
	})
	http.cfgMgr.WatchConfig()
}

func (http *Http) ModuleRunStartup() {
	http.CoreChanged()
}

func (http *Http) ModuleDestruct() {
	http.coreChanged.Store(false)
	//close(http.cfgChangedChan)
	//close(http.done)
	log.Println("[http] destructed")
}

func (http *Http) ModuleRestart() bool {
	//http.done <- struct{}{}
	if http.coreChanged.CompareAndSwap(true, false) {
		log.Println("[http] start restart")
		http.Recreate()
		return true
	}

	return false
}

func (http *Http) ModuleAfterRestart() {
	//http.done <- struct{}{}
	time.Sleep(50 * time.Millisecond)
	log.Println("[http] after restart")
}

func (http *Http) ModuleShutdown() {
	//http.done <- struct{}{}
	log.Println("[http] shutdown")
}

func (http *Http) ModuleLock() singleinstmodule.ModuleCore {
	http.core.Lock()
	return http.core
}

func (http *Http) ModuleUnlock() {
	http.core.Unlock()
	http.CoreChanged()
}

func (http *Http) CoreChanged() {
	// 启动模块
	http.coreChanged.Store(true)
	singleinstmodule.RestartModule(http)
	//http.cfgChangedChan <- struct{}{}
}

func (http *Http) Recreate() {
	http.core.RLock()
	defer http.core.RUnlock()

	cfg := network.ModuleLock().(*core.NetworkCore)
	defer network.ModuleUnlock()

	// 添加中间件
	cfg.HttpMiddlewares = nil
	if http.core.UseMiddleware {
		// 崩溃恢复
		if http.core.UseMiddlewareRecovery {

			cfg.HttpMiddlewares = append(cfg.HttpMiddlewares, gin.Recovery())
		}

		// gin 日志
		if http.core.UseMiddlewareGinLogger {

			cfg.HttpMiddlewares = append(cfg.HttpMiddlewares, gin.Logger())
		}

		// zap 日志
		if http.core.UseMiddlewareZapLogger {

			/*
				conf := logging.GinLoggerConfig{
					Formatter: func(c context.Context, m logging.GinLogDetails) string {
						return "http"
					},
					EnableRequestBody: false,
					EnableContextKeys: false,
				}
				cfg.HttpMiddlewares = append(cfg.HttpMiddlewares, logging.GinLoggerWithConfig(conf))

			*/
		}
	}

	cfg.ListenHttp = http.core.Enable
	if cfg.ListenHttp {
		cfg.ListenIp = http.core.ListenIp
		cfg.ListenPort = http.core.ListenPort
		cfg.HttpReadTimeOut = http.core.ReadTimeOut
		cfg.HttpWriteTimeOut = http.core.WriteTimeOut
		cfg.PathToServiceName = func(c *gin.Context) string {
			pkg := c.Param("pkg")
			service := c.Param("service")
			method := c.Param("method")
			return utils.MakeKey(pkg, service, method)
		}
	}
}

func GetSingleInst() *Http {
	if singleInst == nil {
		once.Do(func() {
			singleInst = new(Http)
		})
	}

	return singleInst
}

// 注册到modules
func init() {
	singleinstmodule.Register(GetSingleInst())
}
