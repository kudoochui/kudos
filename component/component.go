package component

type Component interface {
	OnInit(ServerImpl)
	OnDestroy()
	OnRun(closeSig chan bool)
}

type ServerImpl interface {
	GetServerId() string
	GetComponent(string) Component
}