package msgService

import (
	"github.com/kudoochui/kudos/log"
	"reflect"
	"sync"
)


var msgMgrSingleton *msgMgr
var once sync.Once
var gRouteId uint32 = 0

func GetMsgService() *msgMgr {
	once.Do(func() {
		msgMgrSingleton = &msgMgr{
			msgMap: map[string]uint32{},
			//msgArray: make([]*MsgInfo,0),
			idMap: map[uint32]*MsgInfo{},
		}
	})
	return msgMgrSingleton
}

type msgMgr struct {
	msgMap map[string]uint32
	//msgArray []*MsgInfo
	idMap map[uint32]*MsgInfo
}

type MsgInfo struct {
	Route 		string
	RespId 		uint32
	MsgReqType 	reflect.Type
	MsgRespType reflect.Type
}

func (m *msgMgr) Register(route string, reqId uint32, respId uint32, msgReq interface{}, msgResp interface{}) {
	msgType := reflect.TypeOf(msgReq)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Error("message request pointer required")
	}

	if _, ok := m.msgMap[route]; ok {
		log.Warning("route %s is already registered", route)
		return
	}

	msgRespType := reflect.TypeOf(msgResp)
	if msgRespType == nil || msgRespType.Kind() != reflect.Ptr {
		log.Error("message response pointer required")
	}

	i := new(MsgInfo)
	i.Route = route
	i.RespId = respId
	i.MsgReqType = msgType
	i.MsgRespType = msgRespType

	id := reqId
	if reqId == 0 {
		gRouteId++
		id = gRouteId
	}
	m.msgMap[route] = id
	m.idMap[id] = i
}

func (m *msgMgr) RegisterPush(route string, routeId uint32) {
	if _, ok := m.msgMap[route]; ok {
		log.Warning("route %s is already registered", route)
		return
	}

	i := new(MsgInfo)
	i.Route = route

	id := routeId
	if routeId == 0 {
		gRouteId++
		id = gRouteId
	}
	m.msgMap[route] = id
	m.idMap[id] = i
}

func (m *msgMgr) GetMsgByRouteId(routeId uint32) *MsgInfo {
	return m.idMap[routeId]
}

func (m *msgMgr) GetMsgByRoute(route string) *MsgInfo {
	routeId := m.msgMap[route]
	return m.idMap[routeId]
}

func (m *msgMgr) GetRouteId(route string) uint32 {
	return m.msgMap[route]
}

func (m *msgMgr) GetMsgMap() map[string]uint32 {
	return m.msgMap
}
