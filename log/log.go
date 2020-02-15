// Copyright 2014 mqant Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
