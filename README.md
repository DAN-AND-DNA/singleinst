# singleinst

Some single instance modules

##  example

- [KVStore](https://github.com/DAN-AND-DNA/singleinst-examples)  
- [network](https://github.com/DAN-AND-DNA/singleinst/tree/main/modules/network/internal)
 
## usage

1. [Implementation of a single instance module](https://github.com/DAN-AND-DNA/singleinstmodule)  
2. ModuleCore represents the core of the module, i.e. the data part, either from configuration or pure memory data
3. Module represents the runtime of the module. 
4. when the ModuleCore changes, you need to reconstruct the Module module and re-run the module
5. ModuleLock means locking the ModuleCore to provide external access, ModuleUnlock means unlocking it, which usually notifies the module manager to reconstruct and restart the module, then ModuleRestart is called to rebuild and restart the module runtime.
6. the module manager will automatically run the function starting with ModuleRun, which means run the module
7. After creating the module, you need to register it to the module manager to manage the life cycle, refer to [network](https://github.com/DAN-AND-DNA/singleinst/tree/main/modules/network/internal)

```golang
// module runtime
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

// module data
type ModuleCore interface {
	Lock()
	Unlock()
}
```