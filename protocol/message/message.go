package message

import "fmt"

const (
	TYPE_REQUEST = iota
	TYPE_NOTIFY			//1
	TYPE_RESPONSE		//2
	TYPE_PUSH			//3

	MSG_FLAG_BYTES = 1
	MSG_ROUTE_CODE_BYTES = 2
	MSG_ROUTE_CODE_MAX = 0xffff
	MSG_COMPRESS_ROUTE_MASK = 0x1
	MSG_TYPE_MASK = 0x7
)

// Message protocol encode.
func Encode(id int, msgType int, route uint16, msg []byte) [][]byte {
	idBytes := 0
	if msgHasId(msgType) {
		idBytes = caculateMsgIdBytes(id)
	}
	msgLen := MSG_FLAG_BYTES + idBytes;

	if msgHasRoute(msgType) {
		msgLen += MSG_ROUTE_CODE_BYTES
	}
	buffer := make([]byte, msgLen)
	offset := 0

	// add flag
	offset = encodeMsgFlag(msgType, buffer, offset)

	// add message id
	if msgHasId(msgType) {
		offset = encodeMsgId(id, idBytes, buffer, offset)
	}

	// add route
	if msgHasRoute(msgType) {
		offset = encodeMsgRoute(route, buffer, offset)
	}

	return [][]byte{buffer, msg}
}

// Message protocol decode.
func Decode(buffer []byte) (id int, msgType int, route uint16, body []byte){
	offset := 0

	flag := buffer[offset]
	offset++
	//compressRoute := flag & MSG_COMPRESS_ROUTE_MASK
	msgType = int((flag >> 1) & MSG_TYPE_MASK)

	if msgHasId(msgType) {
		var i uint32 = 0
		m := int(buffer[offset])
		id += (m & 0x7f) << (7*i)
		offset++
		i++
		for ;m >= 128; {
			m = int(buffer[offset])
			id += (m & 0x7f) << (7*i)
			offset++
			i++
		}
	}

	if msgHasRoute(msgType) {
		route = uint16(buffer[offset] << 8 | buffer[offset+1])
		offset += 2
	}

	body = buffer[offset:]
	return
}

func msgHasId(msgType int) bool {
	return msgType == TYPE_REQUEST || msgType == TYPE_RESPONSE
}

func msgHasRoute(msgType int) bool {
	return msgType == TYPE_REQUEST || msgType == TYPE_NOTIFY || msgType == TYPE_PUSH
}

func caculateMsgIdBytes(id int) int {
	l := 0
	for ;id>0; {
		l += 1
		id >>= 7
	}
	return l
}

func encodeMsgFlag(msgType int, buffer []byte, offset int) int {
	if msgType != TYPE_REQUEST && msgType != TYPE_NOTIFY && msgType != TYPE_RESPONSE && msgType !=  TYPE_PUSH {
		fmt.Printf("unkonw message type: %d", msgType)
		return offset
	}

	compressRoute := 1
	if !msgHasRoute(msgType) {
		compressRoute = 0
	}
	buffer[offset] = byte((msgType << 1) | compressRoute)

	return offset + MSG_FLAG_BYTES
}

func encodeMsgId(id int, idBytes int, buffer []byte, offset int) int {
	tmp := id % 128
	next := int(id/128)
	if next != 0 {
		tmp += 128
	}
	buffer[offset] = byte(tmp)
	offset++
	id = next
	for ;id != 0; {
		tmp = id % 128
		next = int(id/128)
		if next != 0 {
			tmp += 128
		}
		buffer[offset] = byte(tmp)
		offset++
		id = next
	}
	return offset
}

func encodeMsgRoute(route uint16, buffer []byte, offset int) int {
	if route > MSG_ROUTE_CODE_MAX {
		fmt.Println("route number is overflow")
		return offset
	}
	buffer[offset] = byte((route >> 8) & 0xff)
	offset++
	buffer[offset] = byte(route & 0xff);
	offset++

	return offset;
}