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
	name    []string
	upacc   []DeltaAcc
	downacc []DeltaAcc
	done    chan bool
}

func (t *netData) setNetData(ifnum int, depth int) {

	t.name = make([]string, ifnum, ifnum)
	t.upacc = make([]DeltaAcc, ifnum, ifnum)
	t.downacc = make([]DeltaAcc, ifnum, ifnum)
	t.done = make(chan bool, 3)

	for i := 0; i < ifnum; i++ {
		t.upacc[i] = *NewDeltaAcc(depth)
		t.downacc[i] = *NewDeltaAcc(depth)
	}
	return
}

func (nd *netData) Init(depth int, rt time.Duration) error {
	noi, err := getifnum()

	if err != nil {
		return err
	}
	nd.setNetData(noi, depth)
	nd.Setup()
	go func() {
		for {
			select {
			case <-nd.done:
				return
			default:
				nd.Update()
				time.Sleep(rt)
			}
		}
	}()
	return nil
}

func (nd *netData) Close() {
	nd.done <- true
}
