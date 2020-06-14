package json

import (
	"github.com/json-iterator/go"
)

type JsonCodec struct {

}

func NewCodec() *JsonCodec {
	p := new(JsonCodec)
	return p
}

// goroutine safe
func (p *JsonCodec) Unmarshal(obj interface{}, data []byte) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	return nil
}

// goroutine safe
func (p *JsonCodec) Marshal(msg interface{}) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(msg)
	return data, err
}
