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

func (s *ServerDefault) GetComponent(name string) component.Component {
	return s.Components[name]
}