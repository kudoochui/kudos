package protobuf

import (
	"github.com/golang/protobuf/proto"
)

// -------------------------
// | id | protobuf message |
// -------------------------
type ProtobufCodec struct {
	littleEndian bool
}

func NewCodec() *ProtobufCodec {
	p := new(ProtobufCodec)
	p.littleEndian = false
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *ProtobufCodec) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// goroutine safe
func (p *ProtobufCodec) Unmarshal(obj interface{}, data []byte) error {
	err := proto.UnmarshalMerge(data, obj.(proto.Message))
	if err != nil {
		return err
	}

	return nil
}

// goroutine safe
func (p *ProtobufCodec) Marshal(msg interface{}) ([]byte, error) {
	data, err := proto.Marshal(msg.(proto.Message))
	return data, err
}
