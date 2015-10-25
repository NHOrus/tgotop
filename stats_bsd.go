// +build freebsd darwin

package main

import (
	"github.com/blabber/go-freebsd-sysctl/sysctl"
)

var pagesize int64
var pae bool

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
	mtemp, err := sysctl.GetInt64("hw.physmem")
	if err != nil {
		panic(err)
	}
	m.memTotal = uint64(mtemp)

	mtemp, err = sysctl.GetInt64("vm.stats.vm.v_free_count")
	if err != nil {
		panic(err)
	}
	m.memFree = uint64(mtemp * pagesize)

	m.memUse = m.memTotal - m.memFree
	m.memPercent = int(m.memUse * 100 / m.memTotal)
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
