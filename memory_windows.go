// memory_windows
// in use only when building for Windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

type memData struct {
	memTotal    uint64
	memFree     uint64
	memUse      uint64
	memPercent  int
	swapTotal   uint64
	swapFree    uint64
	swapUse     uint64
	swapPercent int
}

type MEMORYSTATUSEX struct {
	dwLen                   uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

var (
	kernel32    = syscall.NewLazyDLL("kernel32.dll")
	globMemStat = kernel32.NewProc("GlobalMemoryStatusEx")
	calledMem   MEMORYSTATUSEX
)

func callforMem() {
	calledMem.dwLen = uint32(unsafe.Sizeof(calledMem))
	ret, _, callErr := syscall.Syscall(globMemStat.Addr(), 1, uintptr(unsafe.Pointer(&calledMem)), 0, 0)
	if ret == 0 {
		panic(fmt.Sprintf("%s failed: %v", "GlobalMemoryStatusEx", callErr))
	}
}

func (m *memData) Update() error {
	callforMem()

	m.memTotal = calledMem.ullTotalPhys
	m.memFree = calledMem.ullAvailPhys
	m.memUse = m.memTotal - m.memFree
	m.memPercent = int(m.memUse * 100 / m.memTotal)

	m.swapTotal = calledMem.ullTotalPageFile
	m.swapFree = calledMem.ullAvailPageFile
	m.swapUse = m.swapTotal - m.swapFree
	m.swapPercent = int(m.swapUse * 100 / m.swapTotal)

	return nil
}
