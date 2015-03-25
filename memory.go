// memory
package main

import (
	proc "github.com/cespare/goproc"
)

type memData struct {
	memTotal    uint64
	memFree     uint64
	memPercent  int
	swapTotal   uint64
	swapFree    uint64
	swapPercent int
}

func (m *memData) Update() error {
	t, err := proc.MemInfo()
	if err != nil {
		return err
	}
	m.memTotal = t["MemTotal"]
	m.memFree = t["MemFree"]
	m.memPercent = int((m.memTotal - m.memFree) * 100 / m.memTotal)

	m.swapTotal = t["SwapTotal"]
	m.swapFree = t["SwapFree"]
	m.swapPercent = int((m.swapTotal - m.swapFree) * 100 / m.swapTotal)

	return nil
}
