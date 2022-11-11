package mvc

import (
	mvc "singleinst/core/mvc"
	"singleinst/modules/mvc/internal"
)

type Mvc = internal.Mvc

func RegisterCtrl(module string, ctrl mvc.Controller) {
	internal.GetSingleInst().RegisterCtrl(module, ctrl)
}

func RegisterModel(module string, model mvc.Model) {
	internal.GetSingleInst().RegisterModel(module, model)
}

func RegisterView(module string, view mvc.View) {
	internal.GetSingleInst().RegisterView(module, view)
}
