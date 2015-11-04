// +build freebsd darwin

package main

import (
	"github.com/blabber/go-freebsd-sysctl/sysctl"
)

// #cgo LDFLAGS: -L/usr/lib -lkvm
// #include <kvm.h>
import "C"

var pagesize int64
var pae bool
var kd C.kvm_t
var swapArr C.struct_kvm_swap

func init() {
	var err error
	pagesize, err = sysctl.GetInt64("hw.pagesize")
	if err != nil {
		panic(err)
	}

	_, err = sysctl.GetString("kern.features.pae")
	if err.Error() == "no such file or directory" {
		pae = false
	} else {
		pae = true
	}

}

func (m *memData) Update() error {
	if pae {
		mtemp, err := sysctl.GetInt64("hw.availpages")
		if err != nil {
			panic(err)
		}
		m.memTotal = uint64(mtemp * pagesize)
	}

	var mpage, mtemp int64

	mtemp, err := sysctl.GetInt64("hw.physmem")
	if err != nil {
		panic(err)
	}
	m.memTotal = uint64(mtemp)

	for _, str := range []string{"vm.stats.vm.v_cache_count", "vm.stats.vm.v_free_count"} {

		mpage, err = sysctl.GetInt64(str)
		if err != nil {
			panic(err)
		}
		mtemp = mpage * pagesize
	}
	m.memFree = uint64(mtemp)

	m.memUse = m.memTotal - m.memFree
	m.memPercent = int(m.memUse * 100 / m.memTotal)

	mtemp, err = sysctl.GetInt64("vm.swap_total")
	m.swapTotal = uint64(mtemp)

	err = nil

	i, _ := C.kvm_getswapinfo(&kd, &swapArr, C.int(1), C.int(0))
	if err != nil {
		panic(err)
	}
	if i >= 0 && swapArr.ksw_total != 0 {
		m.swapUse = uint64(swapArr.ksw_used) * uint64(pagesize)
	}
	m.swapPercent = int(m.swapUse * 100 / m.swapTotal)

	return nil
}

func getifnum() (int, error) {
	return 0, nil
}

func (nd *netData) Setup() error {
	return nil
}

func (nd *netData) Update() error {
	return nil
}
