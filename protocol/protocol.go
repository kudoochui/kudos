package protocol

import (
	"bytes"
	"github.com/kudoochui/rpcx/util"
	"sync"
)

var (
	bufferPool = util.NewLimitedPool(512, 4096)
	sPool = &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

func GetPoolMsg() *bytes.Buffer {
	return sPool.Get().(*bytes.Buffer)
}

func FreePoolMsg(buf *bytes.Buffer)  {
	sPool.Put(buf)
}

func GetPoolBuffer(size int) *[]byte {
	return bufferPool.Get(size)
}

func FreePoolBuffer(buf *[]byte)  {
	bufferPool.Put(buf)
}

var poolUint32Data = sync.Pool{
	New: func() interface{} {
		data := make([]byte, 4)
		return &data
	},
}

func GetUint32PoolData() *[]byte {
	return poolUint32Data.Get().(*[]byte)
}

func PutUint32PoolData(buffer *[]byte)  {
	poolUint32Data.Put(buffer)
}