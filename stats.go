// platform-independent system statistics things

package main

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

type netData struct {
	name    []string
	upacc   []DeltaAcc
	downacc []DeltaAcc
}
