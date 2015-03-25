// memory
package main

import (
	proc "github.com/cespare/goproc"
)

type memData struct {
	memTotal    uint64
	memFree     uint64
	memPercent  float64
	swapTotal   uint64
	swapFree    uint64
	swapPercent float64
}

func (m *memData) Update() error {
	t, err := proc.MemInfo()
	if err != nil {
		return err
	}
	m.memTotal = t["MemTotal"]
	m.memFree = t["MemFree"]
	m.memPercent = float64((m.memTotal - m.memFree) / m.memTotal)

	m.swapTotal = t["SwapTotal"]
	m.swapFree = t["SwapFree"]
	m.swapPercent = float64((m.swapTotal - m.swapFree) / m.swapTotal)

	return nil
}
