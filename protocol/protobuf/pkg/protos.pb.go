// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos.proto

package pkg

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type EMsgType int32

const (
	EMsgType_INVALID_TYPE        EMsgType = 0
	EMsgType_TYPE_HEARTBEAT      EMsgType = 70000
	EMsgType_TYPE_COMMON_RESULT  EMsgType = 70001
	EMsgType_TYPE_KICK_BY_SERVER EMsgType = 70002
)

var EMsgType_name = map[int32]string{
	0:     "INVALID_TYPE",
	70000: "TYPE_HEARTBEAT",
	70001: "TYPE_COMMON_RESULT",
	70002: "TYPE_KICK_BY_SERVER",
}

var EMsgType_value = map[string]int32{
	"INVALID_TYPE":        0,
	"TYPE_HEARTBEAT":      70000,
	"TYPE_COMMON_RESULT":  70001,
	"TYPE_KICK_BY_SERVER": 70002,
}

func (x EMsgType) String() string {
	return proto.EnumName(EMsgType_name, int32(x))
}

func (EMsgType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5da3cbeb884d181c, []int{0}
}

type EErrorCode int32

const (
	EErrorCode_INVALID_ERROR_CODE   EErrorCode = 0
	EErrorCode_SUCCESS              EErrorCode = 1
	EErrorCode_ERROR_ROUTE_ID       EErrorCode = 2
	EErrorCode_ERROR_KICK_BY_SERVER EErrorCode = 3
)

var EErrorCode_name = map[int32]string{
	0: "INVALID_ERROR_CODE",
	1: "SUCCESS",
	2: "ERROR_ROUTE_ID",
	3: "ERROR_KICK_BY_SERVER",
}

var EErrorCode_value = map[string]int32{
	"INVALID_ERROR_CODE":   0,
	"SUCCESS":              1,
	"ERROR_ROUTE_ID":       2,
	"ERROR_KICK_BY_SERVER": 3,
}

func (x EErrorCode) String() string {
	return proto.EnumName(EErrorCode_name, int32(x))
}

func (EErrorCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5da3cbeb884d181c, []int{1}
}

type ReqHeartbeat struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqHeartbeat) Reset()         { *m = ReqHeartbeat{} }
func (m *ReqHeartbeat) String() string { return proto.CompactTextString(m) }
func (*ReqHeartbeat) ProtoMessage()    {}
func (*ReqHeartbeat) Descriptor() ([]byte, []int) {
	return fileDescriptor_5da3cbeb884d181c, []int{0}
}

func (m *ReqHeartbeat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqHeartbeat.Unmarshal(m, b)
}
func (m *ReqHeartbeat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqHeartbeat.Marshal(b, m, deterministic)
}
func (m *ReqHeartbeat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqHeartbeat.Merge(m, src)
}
func (m *ReqHeartbeat) XXX_Size() int {
	return xxx_messageInfo_ReqHeartbeat.Size(m)
}
func (m *ReqHeartbeat) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqHeartbeat.DiscardUnknown(m)
}

var xxx_messageInfo_ReqHeartbeat proto.InternalMessageInfo

type RespResult struct {
	Code                 int32    `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RespResult) Reset()         { *m = RespResult{} }
func (m *RespResult) String() string { return proto.CompactTextString(m) }
func (*RespResult) ProtoMessage()    {}
func (*RespResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_5da3cbeb884d181c, []int{1}
}

func (m *RespResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RespResult.Unmarshal(m, b)
}
func (m *RespResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RespResult.Marshal(b, m, deterministic)
}
func (m *RespResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RespResult.Merge(m, src)
}
func (m *RespResult) XXX_Size() int {
	return xxx_messageInfo_RespResult.Size(m)
}
func (m *RespResult) XXX_DiscardUnknown() {
	xxx_messageInfo_RespResult.DiscardUnknown(m)
}

var xxx_messageInfo_RespResult proto.InternalMessageInfo

func (m *RespResult) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *RespResult) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterEnum("pkg.EMsgType", EMsgType_name, EMsgType_value)
	proto.RegisterEnum("pkg.EErrorCode", EErrorCode_name, EErrorCode_value)
	proto.RegisterType((*ReqHeartbeat)(nil), "pkg.ReqHeartbeat")
	proto.RegisterType((*RespResult)(nil), "pkg.RespResult")
}

func init() { proto.RegisterFile("protos.proto", fileDescriptor_5da3cbeb884d181c) }

var fileDescriptor_5da3cbeb884d181c = []byte{
	// 258 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xcf, 0x4a, 0xc3, 0x40,
	0x10, 0x87, 0x4d, 0x1b, 0xff, 0x8d, 0x21, 0x2c, 0x63, 0x91, 0x78, 0x2b, 0x3d, 0x95, 0x1e, 0x3c,
	0xe8, 0x13, 0xa4, 0x9b, 0x81, 0x86, 0x36, 0x8d, 0x4c, 0x36, 0x85, 0x9e, 0x96, 0x16, 0x97, 0x1c,
	0x14, 0x12, 0xb3, 0xf1, 0xe0, 0xeb, 0xe4, 0x09, 0xd5, 0x27, 0x90, 0x6c, 0xf1, 0xe2, 0x69, 0x3f,
	0x7e, 0xdf, 0xe1, 0x5b, 0x06, 0x82, 0xa6, 0xad, 0xbb, 0xda, 0x3e, 0xb8, 0x07, 0xc7, 0xcd, 0x6b,
	0x35, 0x0b, 0x21, 0x60, 0xf3, 0xbe, 0x32, 0x87, 0xb6, 0x3b, 0x9a, 0x43, 0x37, 0x7b, 0x04, 0x60,
	0x63, 0x1b, 0x36, 0xf6, 0xe3, 0xad, 0x43, 0x04, 0x5f, 0xd6, 0x2f, 0x26, 0xf2, 0xa6, 0xde, 0xfc,
	0x9c, 0x1d, 0xa3, 0x80, 0x71, 0x66, 0xab, 0x68, 0x34, 0xf5, 0xe6, 0xd7, 0x3c, 0xe0, 0xa2, 0x82,
	0x2b, 0xca, 0x6c, 0xa5, 0x3e, 0x9b, 0xc1, 0x06, 0xe9, 0x76, 0x17, 0x6f, 0xd2, 0x44, 0xab, 0xfd,
	0x33, 0x89, 0x33, 0x9c, 0x40, 0x38, 0x90, 0x5e, 0x51, 0xcc, 0x6a, 0x49, 0xb1, 0x12, 0x5f, 0xbd,
	0x8f, 0x11, 0xa0, 0x5b, 0x65, 0x9e, 0x65, 0xf9, 0x56, 0x33, 0x15, 0xe5, 0x46, 0x89, 0xef, 0xde,
	0xc7, 0x7b, 0xb8, 0x75, 0x66, 0x9d, 0xca, 0xb5, 0x5e, 0xee, 0x75, 0x41, 0xbc, 0x23, 0x16, 0x3f,
	0xbd, 0xbf, 0xd0, 0x00, 0x44, 0x6d, 0x5b, 0xb7, 0xee, 0x23, 0x77, 0x80, 0x7f, 0x29, 0x62, 0xce,
	0x59, 0xcb, 0x3c, 0x19, 0x82, 0x37, 0x70, 0x59, 0x94, 0x52, 0x52, 0x51, 0x08, 0x0f, 0x11, 0xc2,
	0x93, 0xe4, 0xbc, 0x54, 0xa4, 0xd3, 0x44, 0x8c, 0x30, 0x82, 0xc9, 0x69, 0xfb, 0x97, 0x18, 0x1f,
	0x2f, 0xdc, 0x65, 0x9e, 0x7e, 0x03, 0x00, 0x00, 0xff, 0xff, 0x4e, 0x44, 0x86, 0x8d, 0x29, 0x01,
	0x00, 0x00,
}