// memory
package main

import (
//	proc "github.com/cespare/goproc"
)

type memData struct {
	memTotal    int64
	memFree     int64
	memPercent  float64
	swapTotal   int64
	swapFree    int64
	swapPercent float64
}
