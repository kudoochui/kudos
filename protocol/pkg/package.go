package pkg

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
func Encode(pkgType int, body ...[]byte) [][]byte {
	length := 0
	for i := 0; i < len(body); i++ {
		length += len(body[i])
	}
	head := make([]byte, PKG_HEAD_BYTES)
	head[0] = byte(pkgType & 0xff)
	head[1] = byte(length >> 16)
	head[2] = byte(length >> 8)
	head[3] = byte(length)

	ret := [][]byte{head}
	ret = append(ret, body...)
	return ret
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