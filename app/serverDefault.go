package app

import "github.com/kudoochui/kudos/component"

type ServerDefault struct {
	ServerId string
	//component
	Components map[string]component.Component
}

func NewServerDefault(serverId string) *ServerDefault {
	return &ServerDefault{
		ServerId:   serverId,
		Components: map[string]component.Component{},
	}
}

func (s *ServerDefault) GetServerId() string {
	return s.ServerId
}

func (s *ServerDefault) GetComponent(name string) component.Component {
	return s.Components[name]
}

// Initialize components
func (s *ServerDefault) OnInit() {
	for _,com := range s.Components {
		com.OnInit(s)
	}
}

// Destroy components
func (s *ServerDefault) OnDestroy() {
	for _,com := range s.Components {
		com.OnDestroy()
	}
}

// Run components
func (s *ServerDefault) OnRun(closeSig chan bool) {
	for _,com := range s.Components {
		com.OnRun(closeSig)
	}
}