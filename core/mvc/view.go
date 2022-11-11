package mvc

type View interface {
	Initialize() error
	OnRun()
	OnStop()
	OnClean()
	JustAsView()
}
