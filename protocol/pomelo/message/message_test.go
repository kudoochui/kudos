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
	packer := NewMessagePacker()
	assert.DeepEqual(t, packer.Encode(1,0,3,msg), []byte{1, 1, 0, 3, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func Test_Main(t *testing.T) {
	msg := []byte{1,2,3,4,5,6,7,8,9}
	packer := NewMessagePacker()
	buffer := packer.Encode(1,0,3, msg)
	t.Log(buffer)
	t.Log(Decode(buffer))
}

func benchmarkEncode(b *testing.B, buffer []byte) {
	b.ReportAllocs()
	b.SetBytes(int64(len(buffer)))
	packer := NewMessagePacker()
	for i := 0; i < b.N; i++ {
		packer.Encode(1,0,3, buffer)
	}
}

func BenchmarkEncode(b *testing.B) {
	buffer := make([]byte, 1024*1024)

	b.Run("Encode=16", func(bb *testing.B) { benchmarkEncode(bb, buffer[:16])})
	b.Run("Encode=32", func(bb *testing.B) { benchmarkEncode(bb, buffer[:32])})
	b.Run("Encode=64", func(bb *testing.B) { benchmarkEncode(bb, buffer[:64])})
	b.Run("Encode=128", func(bb *testing.B) { benchmarkEncode(bb, buffer[:128])})
	b.Run("Encode=512", func(bb *testing.B) { benchmarkEncode(bb, buffer[:512])})
	b.Run("Encode=1k", func(bb *testing.B) { benchmarkEncode(bb, buffer[:1024])})
	b.Run("Encode=10k", func(bb *testing.B) { benchmarkEncode(bb, buffer[:1024*10])})
	b.Run("Encode=64k", func(bb *testing.B) { benchmarkEncode(bb, buffer[:1024*64])})
	b.Run("Encode=100k", func(bb *testing.B) { benchmarkEncode(bb, buffer[:1024*100])})
	b.Run("Encode=1M", func(bb *testing.B) { benchmarkEncode(bb, buffer[:1024*1024])})
}