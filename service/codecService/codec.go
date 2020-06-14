package codecService

import (
	"github.com/kudoochui/kudos/service/codecService/json"
	"github.com/kudoochui/kudos/service/codecService/protobuf"
	"sync"
)

const (
	TYPE_CODEC_JSON = "json"
	TYPE_CODEC_PROTOBUF = "protobuf"
)

var codecType = TYPE_CODEC_JSON
var codec Codec
var once sync.Once

type Codec interface {
	// must goroutine safe
	Unmarshal(obj interface{}, data []byte) error
	// must goroutine safe
	Marshal(msg interface{}) ([]byte, error)
}

// Change codec type in the main
func SetCodecType(t string) {
	codecType = t
}

func GetCodecService() Codec {
	once.Do(func() {
		switch codecType {
		case TYPE_CODEC_JSON:
			codec = json.NewCodec()
		case TYPE_CODEC_PROTOBUF:
			codec = protobuf.NewCodec()
		}
	})
	return codec
}