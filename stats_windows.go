// stat_windows
// system stats aquisition function
// in use only when building for Windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	ERROR_INSUFFICIENT_BUFFER uintptr = 122
	ERROR_INVALID_PARAMETER   uintptr = 87
	ERROR_NOT_SUPPORTED       uintptr = 50
	NO_ERROR                  uintptr = 0
)

var (
	kernel32    = syscall.NewLazyDLL("kernel32.dll")
	globMemStat = kernel32.NewProc("GlobalMemoryStatusEx")
	calledMem   memstatex
	iphelper    = syscall.NewLazyDLL("Iphlpapi.dll")
	getIfTable  = iphelper.NewProc("GetIfTable")
	getIfEntry  = iphelper.NewProc("GetIfEntry")
	currIfTable *iftable
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
	currIfTable := new(iftable)
	size := uint32(unsafe.Sizeof(currIfTable))
	dsize := uint32(unsafe.Sizeof(currIfTable.table))
	var ifnum int
	var bOrder int32
	ret, _, callErr := getIfTable.Call(uintptr(unsafe.Pointer(currIfTable)), uintptr(unsafe.Pointer(&size)), uintptr(unsafe.Pointer(&bOrder)))
	if callErr != nil {
		panic(callErr)
	}
	if ret == ERROR_INVALID_PARAMETER || ret == ERROR_NOT_SUPPORTED {
		panic(ret)
	}
	if ret == ERROR_INSUFFICIENT_BUFFER {
		//magic pointer size math!
		var ir_temp ifrow
		rowsize := uint32(unsafe.Sizeof(ir_temp))
		if (size % rowsize) == dsize {
			panic("size mismatch, i fear")
		}
		ifnum = int(size / rowsize)
		currIfTable.table = make([]ifrow, ifnum, ifnum)
		ret, _, callErr := getIfTable.Call(uintptr(unsafe.Pointer(currIfTable)), uintptr(unsafe.Pointer(&size)), uintptr(unsafe.Pointer(&bOrder)))
	}
	if ret == NO_ERROR {
		return int(currIfTable.dwNumEntries), callErr
	}
	return int(currIfTable.dwNumEntries), callErr
}

func (nd *netData) Setup() error {
}

func (nd *netData) Update() error {}
