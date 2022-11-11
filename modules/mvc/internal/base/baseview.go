package base

type BaseView struct {
}

func (baseView *BaseView) Initialize() error {
	return nil
}

func (baseView *BaseView) OnRun() {

}

func (baseView *BaseView) OnStop() {

}

func (baseView *BaseView) OnClean() {

}

func (baseView *BaseView) JustAsView() {
}
