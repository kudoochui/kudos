package log

import (
	logs "github.com/kudoochui/kudos/log/beego"
	"testing"
)

func TestLog(t *testing.T)  {
	LogBeego().SetLogger("console")
	LogBeego().SetLogger(logs.AdapterFile,`{"filename":"test.log"}`)
	LogBeego().EnableFuncCallDepth(true)
	LogBeego().SetLogFuncCallDepth(3)
	//LogBeego().Async()

	Debug("my book is bought in the year of ", 2016)
	Info("this %s cat is %v years old", "yellow", 3)
	Notice("ABCD")
	Warning("json is a type of kv like", map[string]int{"key": 2016})
	Error(1024, "is a very", "good game")
	Critical("oh,crash")
	Alert(1E3)
	Emergency(232.32342233)
}
