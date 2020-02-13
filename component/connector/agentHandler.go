package connector

import (
	"encoding/json"
	"github.com/kudoochui/kudos/log"
	"github.com/kudoochui/kudos/protocol/message"
	"github.com/kudoochui/kudos/protocol/pkg"
	"github.com/kudoochui/kudos/rpc"
	"github.com/kudoochui/kudos/service/codecService"
	"github.com/kudoochui/kudos/service/msgService"
	"github.com/kudoochui/kudos/utils/timer"
)

const (
	CODE_OK = 200
	CODE_USE_ERROR = 500
	CODE_OLD_CLIENT = 501
)

type agentHandler struct {
	agent 	*agent
	timerHandler *timer.Timer
}

func NewAgentHandler(a *agent) *agentHandler {
	return &agentHandler{agent:a}
}

func (h *agentHandler) Handle(pkgType int, body []byte) {
	switch pkgType {
	case pkg.TYPE_HANDSHAKE:
		h.handleHandshake(pkgType, body)
	case pkg.TYPE_HANDSHAKE_ACK:
		h.handleHandshakeAck(pkgType, body)
	case pkg.TYPE_HEARTBEAT:
		h.handleHeartbeat(pkgType, body)
	case pkg.TYPE_DATA:
		h.handleData(pkgType, body)
	}
}

func (h *agentHandler) handleHandshake(pkgType int, body []byte) {
	var message map[string]json.RawMessage
	err := json.Unmarshal(body, &message)
	if err != nil {
		log.Error("handshake decode error: %v", err)
		h.processError(CODE_USE_ERROR)
		return
	}

	if message["sys"] == nil {
		h.processError(CODE_USE_ERROR)
		return
	}

	sys := make(map[string]interface{})
	res := make(map[string]interface{})
	res["code"] = CODE_OK
	res["sys"] = sys
	sys["heartbeat"] = 10
	sys["useDict"] = true
	sys["dict"] = msgService.GetMsgService().GetMsgMap()

	bin,_ := json.Marshal(res)
	buffer := pkg.Encode(pkg.TYPE_HANDSHAKE, bin)
	h.agent.Write(buffer...)
}

func (h *agentHandler) handleHandshakeAck(pkgType int, body []byte) {
	h.handleHeartbeat(pkgType, body)
}

func (h *agentHandler) handleHeartbeat(pkgType int, body []byte) {
	buffer := pkg.Encode(pkg.TYPE_HEARTBEAT, nil)
	h.agent.Write(buffer...)

	if h.timerHandler != nil {
		h.agent.connector.timers.ClearTimeout(h.timerHandler)
	}

	//heartbeat timeout close the socket
	h.timerHandler = h.agent.connector.timers.AfterFunc(2*h.agent.connector.opts.HeartbeatTimeout, func() {
		h.agent.Close()
	})
}

func (h *agentHandler) handleData(pkgType int, body []byte) {
	msgId, msgType, route, data := message.Decode(body)
	//_ = msgId
	_ = msgType
	m, err := codecService.GetCodecService().Unmarshal(route, data)
	if err != nil {
		log.Error("unmarshal error: %v", err)
		return
	}
	call := m.(*rpc.Call)
	call.MsgId = msgId
	call.Session = h.agent.session
	h.agent.connector.Route.Go(call)
}

func (h *agentHandler) processError(code int){
	r := make(map[string]int)
	r["code"] = code
	bin,_ := json.Marshal(r)
	buffer := pkg.Encode(pkg.TYPE_HANDSHAKE, bin)
	h.agent.Write(buffer...)
}