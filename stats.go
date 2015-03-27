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

func newNetData(ifnum int, depth int) *netData {
	var t *netData = new(netData)
	t.name = make([]string, ifnum, ifnum)
	t.upacc = make([]DeltaAcc, ifnum, ifnum)
	t.downacc = make([]DeltaAcc, ifnum, ifnum)

	for i := 0; i < ifnum; i++ {
		t.upacc[i] = *NewDeltaAcc(depth)
		t.downacc[i] = *NewDeltaAcc(depth)
	}
	return t
}
