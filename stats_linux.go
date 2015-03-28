// stat_linux
// system stats aquisition function
// for Linux and *BSD with mounted procfs

package main

import (
	proc "github.com/cespare/goproc"
	//spew "github.com/davecgh/go-spew/spew"
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

func getifnum() (int, error) {
	r, _, err := proc.NetDevStats()
	if err != nil {
		return 0, err
	}
	return len(r), nil
}

func (nd *netData) Update() error {
	i := 0
	r, t, err := proc.NetDevStats()
	if err != nil {
		return err
	}
	for key, value := range r {
		nd.name[i] = key
		nd.downacc[i].Push(uint64(value["bytes"]))
		nd.upacc[i].Push(uint64(t[key]["bytes"]))
		i++
	}
	return nil
}
