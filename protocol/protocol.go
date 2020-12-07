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
