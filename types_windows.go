package main

//structure used by GlobalMemoryStatusEx from kernel32.dll , information about memory and things
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
