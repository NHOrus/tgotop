// memory_linux
package main

import (
	proc "github.com/cespare/goproc"
)

func (m *memData) Update() error {
	t, err := proc.MemInfo()
	if err != nil {
		return err
	}
	m.memTotal = t["MemTotal"]
	m.memFree = t["MemFree"]
	m.memUse = m.memTotal - m.memFree
	m.memPercent = int(m.memUse * 100 / m.memTotal)

	m.swapTotal = t["SwapTotal"]
	m.swapFree = t["SwapFree"]
	m.swapUse = m.swapTotal - m.swapFree
	m.swapPercent = int(m.swapUse * 100 / m.swapTotal)

	return nil
}
