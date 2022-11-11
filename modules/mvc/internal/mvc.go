package internal

import (
	"github.com/dan-and-dna/singleinstmodule"
	"log"
	cmvc "singleinst/core/mvc"
	"sync"
)

var (
	mvcSingleInst *Mvc = nil
	once          sync.Once
)

type Mvc struct {
	order  map[string]struct{}
	ctrls  map[string]cmvc.Controller
	models map[string]cmvc.Model
	views  map[string]cmvc.View
}

func (mvc *Mvc) ModuleConstruct() {
	if mvc.order == nil {
		mvc.order = make(map[string]struct{})
	}

	if mvc.ctrls == nil {
		mvc.ctrls = make(map[string]cmvc.Controller)
	}

	if mvc.models == nil {
		mvc.models = make(map[string]cmvc.Model)
	}

	if mvc.views == nil {
		mvc.views = make(map[string]cmvc.View)
	}
	log.Println("[mvc] constructed")
}

func (mvc *Mvc) ModuleDestruct() {
	for k := range mvc.order {
		if model, ok := mvc.models[k]; ok {
			model.OnClean()
		}

		if view, ok := mvc.views[k]; ok {
			view.OnClean()
		}

		if ctrl, ok := mvc.ctrls[k]; ok {
			ctrl.OnClean()
		}
		log.Printf("[mvc] %s: ok\n", k)
	}
	mvc.order = nil
	mvc.ctrls = nil
	mvc.models = nil
	mvc.views = nil
	log.Println("[mvc] destructed")
}

func (mvc *Mvc) ModuleLock() singleinstmodule.ModuleCore {
	return nil
}

func (mvc *Mvc) ModuleUnlock() {
}

func (mvc *Mvc) RegisterCtrl(module string, ctrl cmvc.Controller) {
	if err := ctrl.Initialize(); err != nil {
		panic(err)
	}

	mvc.ctrls[module] = ctrl
	mvc.order[module] = struct{}{}
}

func (mvc *Mvc) RegisterModel(module string, model cmvc.Model) {
	if err := model.Initialize(); err != nil {
		panic(err)
	}

	mvc.models[module] = model
	mvc.order[module] = struct{}{}
}

func (mvc *Mvc) RegisterView(module string, view cmvc.View) {

	if err := view.Initialize(); err != nil {
		panic(err)
	}

	mvc.views[module] = view
	mvc.order[module] = struct{}{}
}

func (mvc *Mvc) ModuleRunStartup() {
	for k := range mvc.order {
		if model, ok := mvc.models[k]; ok {
			model.OnRun()
		}

		if view, ok := mvc.views[k]; ok {
			view.OnRun()
		}

		if ctrl, ok := mvc.ctrls[k]; ok {
			ctrl.OnRun()
		}

		log.Printf("[mvc] %s: ok\n", k)
	}
}

func (mvc *Mvc) ModuleShutdown() {
	for k := range mvc.order {
		if model, ok := mvc.models[k]; ok {
			model.OnStop()
		}

		if view, ok := mvc.views[k]; ok {
			view.OnStop()
		}

		if ctrl, ok := mvc.ctrls[k]; ok {
			ctrl.OnStop()
		}
		log.Printf("[mvc] %s: ok\n", k)
	}
}

func GetSingleInst() *Mvc {
	if mvcSingleInst == nil {
		once.Do(func() {
			mvcSingleInst = new(Mvc)
		})
	}

	return mvcSingleInst
}

func init() {
	// 注册模块
	singleinstmodule.Register(GetSingleInst())
}
