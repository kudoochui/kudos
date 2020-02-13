package component

type Component interface {
	OnInit()
	OnDestroy()
	Run(closeSig chan bool)
}

type ServerImpl interface {
	GetComponent() Component
}