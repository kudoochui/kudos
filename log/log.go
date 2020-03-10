package log

import (
	beegolog "github.com/kudoochui/kudos/log/beego"
	"sync"
)

var beego *beegolog.BeeLogger
var once sync.Once

func LogBeego() *beegolog.BeeLogger {
	once.Do(func() {
		beego = beegolog.NewLogger()
	})
	return beego
}


// Emergency logs a message at emergency level.
func Emergency(f interface{}, v ...interface{}) {
	LogBeego().Emergency(beegolog.FormatLog(f, v...))
}

// Alert logs a message at alert level.
func Alert(f interface{}, v ...interface{}) {
	LogBeego().Alert(beegolog.FormatLog(f, v...))
}

// Critical logs a message at critical level.
func Critical(f interface{}, v ...interface{}) {
	LogBeego().Critical(beegolog.FormatLog(f, v...))
}

// Error logs a message at error level.
func Error(f interface{}, v ...interface{}) {
	LogBeego().Error(beegolog.FormatLog(f, v...))
}

// Warning logs a message at warning level.
func Warning(f interface{}, v ...interface{}) {
	LogBeego().Warn(beegolog.FormatLog(f, v...))
}

// Notice logs a message at notice level.
func Notice(f interface{}, v ...interface{}) {
	LogBeego().Notice(beegolog.FormatLog(f, v...))
}

// Info compatibility alias for Warning()
func Info(f interface{}, v ...interface{}) {
	LogBeego().Info(beegolog.FormatLog(f, v...))
}

// Debug logs a message at debug level.
func Debug(f interface{}, v ...interface{}) {
	LogBeego().Debug(beegolog.FormatLog(f, v...))
}
