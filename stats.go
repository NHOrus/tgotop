// platform-independent system statistics things

package main

import (
	//spew "github.com/davecgh/go-spew/spew"
	"time"
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

type netData struct {
	size    int
	name    []string
	upacc   []DeltaAcc
	downacc []DeltaAcc
	done    chan bool
}

/*
type cpuData struct {
	//usage = user+system/(user+system+idle)
	CPUhist []int
	CPUusage []int
} */

func (nd *netData) setNetData(ifnum int, depth int) {
	nd.size = ifnum
	nd.name = make([]string, 0, ifnum)
	nd.upacc = make([]DeltaAcc, ifnum, ifnum)
	nd.downacc = make([]DeltaAcc, ifnum, ifnum)
	nd.done = make(chan bool, 3)

	for i := 0; i < ifnum; i++ {
		nd.upacc[i] = *NewDeltaAcc(depth)
		nd.downacc[i] = *NewDeltaAcc(depth)
	}
	return
}

func (nd *netData) Init(depth int, ptk <-chan time.Time) error {
	noi, err := getifnum()

	if err != nil {
		return err
	}
	nd.setNetData(noi, depth)
	err = nd.Setup()
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-nd.done:
				return
			case <-ptk:
				err = nd.Update()
				if err != nil {
					return
				}
			}
		}
	}()
	return err
}

func (nd *netData) Close() {
	nd.done <- true
}

func getnetdata(i int, mult int, acc *DeltaAcc) (r float32) {
	r, err := acc.Average(mult)
	if err != nil && err != ErrSmallSample {
		panic(err)
	}
	return r * float32(mult)
}

func (nd *netData) GetD(i int, mult int) float32 {
	return getnetdata(i, mult, &nd.downacc[i])
}
func (nd *netData) GetU(i int, mult int) float32 {
	return getnetdata(i, mult, &nd.upacc[i])
}
