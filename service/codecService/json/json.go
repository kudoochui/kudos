package json

import (
	"encoding/json"
	"fmt"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/service/msgService"
	"reflect"
	"strings"
)

type JsonCodec struct {

}

func NewCodec() *JsonCodec {
	p := new(JsonCodec)
	return p
}

// goroutine safe
func (p *JsonCodec) Unmarshal(route uint16, data []byte) (interface{}, error) {
	i := msgService.GetMsgService().GetMsgByRouteId(route)
	if i == nil {
		return nil, fmt.Errorf("invalid route id")
	}
	call := new(rpc.Call)
	call.MsgReq = reflect.New(i.MsgReqType.Elem()).Interface()
	call.MsgResp = reflect.New(i.MsgRespType.Elem()).Interface()
	err := json.Unmarshal(data, call.MsgReq)
	if err != nil {
		return nil, err
	}
	rr := strings.Split(i.Route, ".")
	if len(rr) < 2 {
		log.Error("service format error")
	}
	call.ServicePath = rr[0]
	call.ServiceName = rr[1]

	return call, nil
}

// goroutine safe
func (p *JsonCodec) Marshal(msg interface{}) ([]byte, error) {
	//msgType := reflect.TypeOf(msg)
	//if msgType == nil || msgType.Kind() != reflect.Ptr {
	//	return nil, errors.New("json message pointer required")
	//}

	data, err := json.Marshal(msg)
	return data, err
}
