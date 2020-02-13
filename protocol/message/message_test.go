package message

import (
	"gotest.tools/assert"
	"testing"
)

func TestMsgIdLen(t *testing.T) {
	assert.Equal(t, caculateMsgIdBytes(1), 1)
	assert.Equal(t, caculateMsgIdBytes(127), 1)
	assert.Equal(t, caculateMsgIdBytes(128), 2)
	var b []byte = nil
	t.Log(len(b))
}

func Test_EncodeMsgId(t *testing.T) {
	buffer := make([]byte, 10)
	assert.Equal(t, encodeMsgId(1,1,buffer,1), 2)
	assert.DeepEqual(t, buffer, []byte{0, 1, 0, 0, 0, 0, 0, 0, 0, 0})

	buffer1 := make([]byte, 10)
	assert.Equal(t, encodeMsgId(128,2,buffer1,1), 3)
	assert.DeepEqual(t, buffer1, []byte{0, 129, 0, 0, 0, 0, 0, 0, 0, 0})
}

func Test_Encode(t *testing.T) {
	msg := []byte{1,2,3,4,5,6,7,8,9}
	assert.DeepEqual(t, Encode(1,0,3,msg), []byte{1, 1, 0, 3, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func Test_Main(t *testing.T) {
	msg := []byte{1,2,3,4,5,6,7,8,9}
	buffer := Encode(1,0,3, msg)
	t.Log(buffer)
	t.Log(Decode(buffer))
}