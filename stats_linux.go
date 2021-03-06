// stat_linux
// system stats aquisition function
// for Linux and *BSD with mounted procfs

package main

import (
	proc "github.com/cespare/goproc"
	//	spew "github.com/davecgh/go-spew/spew"
)

/*
var (
	CPUuser DeltaAcc
	CPUsuser []DeltaAcc
	CPUsys DeltaAcc
	CPUssys []DeltaAcc
	CPUidle DeltaAcc
	CPUsidle []DeltaAcc)
*/

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

func (nd *netData) Setup() error {
	r, t, err := proc.NetDevStats()
	if err != nil {
		return err
	}
	for key := range r {
		nd.name = append(nd.name, key)
	}

	for i := 0; i < nd.size; i++ {
		nd.downacc[i].Push(uint64(r[nd.name[i]]["bytes"]))
		nd.upacc[i].Push(uint64(t[nd.name[i]]["bytes"]))
	}
	return nil
}

func (nd *netData) Update() error {
	r, t, err := proc.NetDevStats()
	if err != nil {
		return err
	}
	for i, key := range nd.name {
		nd.downacc[i].Push(uint64(r[key]["bytes"]))
		nd.upacc[i].Push(uint64(t[key]["bytes"]))
	}
	return nil
}
