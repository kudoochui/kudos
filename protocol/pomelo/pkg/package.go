package pkg

import "github.com/kudoochui/kudos/protocol"

const (
	TYPE_NULL = iota
	TYPE_HANDSHAKE			//1
	TYPE_HANDSHAKE_ACK		//2
	TYPE_HEARTBEAT			//3
	TYPE_DATA				//4
	TYPE_KICK				//5

	PKG_HEAD_BYTES = 4
)

/**
 * Package protocol encode.
 *
 * Pomelo package format:
 * +------+-------------+------------------+
 * | type | body length |       body       |
 * +------+-------------+------------------+
 *
 * Head: 4bytes
 *   0: package type,
 *      1 - handshake,
 *      2 - handshake ack,
 *      3 - heartbeat,
 *      4 - data
 *      5 - kick
 *   1 - 3: big-endian body length
 * Body: body length bytes
 */
func Encode(pkgType int, body []byte) *[]byte {
	length := 0
	var buffer *[]byte
	if pkgType == TYPE_DATA {
		length = len(body) - PKG_HEAD_BYTES
		buffer = &body
		(*buffer)[0] = byte(pkgType & 0xff)
		(*buffer)[1] = byte(length >> 16)
		(*buffer)[2] = byte(length >> 8)
		(*buffer)[3] = byte(length)
	} else {
		length = len(body)
		buffer = protocol.GetPoolBuffer(length + PKG_HEAD_BYTES)
		(*buffer)[0] = byte(pkgType & 0xff)
		(*buffer)[1] = byte(length >> 16)
		(*buffer)[2] = byte(length >> 8)
		(*buffer)[3] = byte(length)
		copy((*buffer)[PKG_HEAD_BYTES:], body)
	}
	return buffer
}

/**
 * Package protocol decode.
 * See encode for package format.
 */
func Decode(buffer []byte) (pkgType int, body []byte) {
	pkgType = int(buffer[0])
	body = buffer[PKG_HEAD_BYTES:]
	return
}