package pkg

import (
	"encoding/binary"
	"github.com/kudoochui/kudos/protocol"
)

const (
	PKG_SIZE_BYTES = 4
	PKG_TYPE_BYTES = 4
)

var gLittleEndian bool = false

func SetByteOrder(littleEndian bool) {
	gLittleEndian = littleEndian
}

func GetByteOrder() bool {
	return gLittleEndian
}
/**
 * Package protocol encode.
 *
 * +------+------+------------------+
 * | size | type |       body       |
 * +------+------+------------------+
 *
 */
func Encode(pkgType uint32, body []byte) *[]byte {
	length := 0
	length = len(body) + PKG_TYPE_BYTES
	buffer := protocol.GetPoolBuffer(length + PKG_SIZE_BYTES)
	if gLittleEndian {
		binary.LittleEndian.PutUint32(*buffer, uint32(length))
		binary.LittleEndian.PutUint32((*buffer)[PKG_SIZE_BYTES:], pkgType)
	} else {
		binary.BigEndian.PutUint32(*buffer, uint32(length))
		binary.BigEndian.PutUint32((*buffer)[PKG_SIZE_BYTES:], pkgType)
	}
	copy((*buffer)[PKG_SIZE_BYTES + PKG_TYPE_BYTES:], body)
	return buffer
}

/**
 * Package protocol decode.
 * See encode for package format.
 */
func Decode(buffer []byte) (length int, pkgType uint32, body []byte) {
	if gLittleEndian {
		length = int(binary.LittleEndian.Uint32(buffer))
		pkgType = binary.LittleEndian.Uint32(buffer[PKG_SIZE_BYTES:])
	} else {
		length = int(binary.BigEndian.Uint32(buffer))
		pkgType = binary.BigEndian.Uint32(buffer[PKG_SIZE_BYTES:])
	}
	body = buffer[PKG_SIZE_BYTES+PKG_TYPE_BYTES:]
	return
}