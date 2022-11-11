# singleinst

一些单例模块

## 例子
- [KVStore](https://github.com/DAN-AND-DNA/singleinst-examples)  
- [network](https://github.com/DAN-AND-DNA/singleinst/tree/main/modules/network/internal)
 
## usage
1. [单例模块的实现](https://github.com/DAN-AND-DNA/singleinstmodule)  
2. ModuleCore 代表模块的核心，即数据部分，可以来自配置，也可以为纯内存数据
3. Module 代表模块的运行时
4. 当ModuleCore发生变化时候，就需要重新构造Module模块，并重新运行模块
5. ModuleLock表示锁住ModuleCore，提供给外部访问，ModuleUnlock表示解锁，一般会通知模块管理器来重新构建和重启模块，这时候就会调用ModuleRestart来重建和重启模块运行时
6. 模块管理器会自动运行以ModuleRun开头的函数，代表运行模块
7. 创建完模块，需要注册到模块管理器，来管理生命周期，参考[network](https://github.com/DAN-AND-DNA/singleinst/tree/main/modules/network/internal)

```golang
// 代表单例模块
type Module interface {
	ModuleConstruct()
	ModuleDestruct() 
	ModuleLock() ModuleCore
	ModuleUnlock()
}

type ModuleCanRestart interface {
	ModuleRestart() bool // for module
}

type ModuleCanBeforeRun interface {
	ModuleBeforeRun(string) // for run
}

type ModuleCanAfterRun interface {
	ModuleAfterRun(string) // for run
}

type ModuleCanShutdown interface {
	ModuleShutdown() // for run
}

type ModuleCanAfterRestart interface {
	AfterRestart()
}

type ModuleCore interface {
	Lock()
	Unlock()
}
```
