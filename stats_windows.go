// stat_windows
// system stats aquisition function
// in use only when building for Windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32    = syscall.NewLazyDLL("kernel32.dll")
	globMemStat = kernel32.NewProc("GlobalMemoryStatusEx")
	calledMem   memstatex
	iphelper    = syscall.NewLazyDLL("Iphlpapi.dll")
)

func extMemInfo() {
	calledMem.dwLen = uint32(unsafe.Sizeof(calledMem))
	ret, _, callErr := globMemStat.Call(uintptr(unsafe.Pointer(&calledMem)))
	if ret == 0 {
		panic(fmt.Sprintf("%s failed: %v", "GlobalMemoryStatusEx", callErr))
	}
}

func (m *memData) Update() error {
	extMemInfo()

	m.memTotal = calledMem.ullTotalPhys
	m.memFree = calledMem.ullAvailPhys
	m.memUse = m.memTotal - m.memFree
	m.memPercent = int(m.memUse * 100 / m.memTotal)

	//massage, because paged memory includes pagefile and RAM and manual removal is it for now
	//on the plus size, m.swapTotal is now equal to actual pagefile.sys size, so everything is mostly ok
	m.swapTotal = calledMem.ullTotalPageFile - calledMem.ullTotalPhys

	//check in case of unexpected garbage
	if calledMem.ullAvailPhys >= calledMem.ullAvailPageFile {
		m.swapFree = 0
	} else {
		m.swapFree = calledMem.ullAvailPageFile - calledMem.ullAvailPhys
	}
	m.swapUse = m.swapTotal - m.swapFree
	m.swapPercent = int(m.swapUse * 100 / m.swapTotal)

	return nil
}

func getifnum() (int, error) {
	return 1, nil
}

func (nd *netData) Setup() (err error) {
	nd.name = append(nd.name, "Dummy interface")
	for i := 0; i < 2*mult; i++ {
		nd.upacc[0].Push(0)
		nd.downacc[0].Push(0)
	}
	return

}

func (nd *netData) Update() (err error) {
	return
}
