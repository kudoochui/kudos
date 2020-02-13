package idService

import (
	"github.com/kudoochui/kudos/config"
	"sync"
)

var node *Node
var once sync.Once

func GenerateID() ID {
	once.Do(func() {
		node, _ = NewNode(config.NodeId)
	})
	return node.Generate()
}