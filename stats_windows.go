// memory_windows
// in use only when building for Windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

type memstatex struct {
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
	calledMem   memstatex
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

	m.swapTotal = calledMem.ullTotalPageFile
	m.swapFree = calledMem.ullAvailPageFile
	m.swapUse = m.swapTotal - m.swapFree
	m.swapPercent = int(m.swapUse * 100 / m.swapTotal)

	return nil
}
