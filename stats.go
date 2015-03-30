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

func (nd *netData) Init(depth int, rt time.Duration) error {
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
			default:
				err = nd.Update()
				if err != nil {
					return
				}
				time.Sleep(rt)
			}
		}
	}()
	return err
}

func (nd *netData) Close() {
	nd.done <- true
}

func (nd *netData) GetD(i int, mult int) (d float32) {
	d, err := nd.downacc[i].Average(mult)
	if err != nil && err != ErrSmallSample {
		panic(err)
	}
	return d * float32(mult)
}
func (nd *netData) GetU(i int, mult int) (u float32) {
	u, err := nd.upacc[i].Average(mult)
	if err != nil && err != ErrSmallSample {
		panic(err)
	}
	return u * float32(mult)
}
