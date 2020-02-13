package protobuf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/siddontang/go/log"
	"hash/crc32"
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/service/msgService"
	"reflect"
	"strings"
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
func (p *ProtobufCodec) Unmarshal(route uint16, data []byte) (interface{}, error) {
	info := msgService.GetMsgService().GetMsgByRouteId(route)
	if info == nil {
		return nil, fmt.Errorf("invalid route id")
	}

	call := new(rpc.Call)
	call.MsgReq = reflect.New(info.MsgReqType.Elem()).Interface()
	call.MsgResp = reflect.New(info.MsgRespType.Elem()).Interface()
	proto.UnmarshalMerge(data, call.MsgReq.(proto.Message))
	rr := strings.Split(info.Route, ".")
	if len(rr) < 2 {
		log.Error("service format error")
	}
	call.ServicePath = rr[0]
	call.ServiceName = rr[1]
	return call, nil
}

// goroutine safe
func (p *ProtobufCodec) Marshal(msg interface{}) ([]byte, error) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return nil, errors.New("pb message pointer required")
	}
	msgID := msgType.Elem().Name()
	_id := crc32.ChecksumIEEE([]byte(msgID))

	id := make([]byte, 4)
	if p.littleEndian {
		binary.LittleEndian.PutUint32(id, _id)
	} else {
		binary.BigEndian.PutUint32(id, _id)
	}

	// data
	data, err := proto.Marshal(msg.(proto.Message))
	return data, err
}
