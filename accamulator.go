// accaverager
package main

import "errors"

var (
	EBIGWINDOW   = errors.New("averaging window is bigger than accumulator")
	ESMALLSAMPLE = errors.New("not enough data saved")
)

type deltaAcc struct {
	size   int
	ptr    int
	last   uint64
	deltas []int64
	full   bool
}

func newAccumil(s int) *deltaAcc {
	return &deltaAcc{
		size:   s,
		deltas: make([]int64, s, s),
	}
}

//Purge cleans up the accumulator and returns it sparkingly clean
func (a *deltaAcc) Purge() {
	a.ptr = 0
	a.last = 0
	for i := 0; i < a.size; i++ {
		a.deltas[i] = 0
	}
	a.full = false
}

//Push adds new delta value into accumulator, going in circles when accumulator is full
func (a *deltaAcc) Push(v uint64) {
	var dlt int64
	if a.last >= v {
		dlt = int64(a.last - v)
	} else { //stupid unsigned math
		dlt = -int64(v - a.last)
	}
	if a.ptr == a.size-1 {
		a.full = true
		a.ptr = 0
	} else {
		a.ptr = a.ptr + 1
	}
	a.deltas[a.ptr] = dlt
}

func (a *deltaAcc) Average(w int) (avg float32, err error) {
	var sum int64 //here we add every delta in the window

	//if window is bigger than our accumulator, we fail horribly
	if w > a.size {
		return 0, EBIGWINDOW
	}
	//if we don't have enough data to fill window, we fail, less horribly
	if !a.full && a.ptr < w {
		return 0, ESMALLSAMPLE
	}
	//if we need to get bits from different ends of our slice, let it be so
	if a.ptr < w {
		for _, v := range a.deltas[:a.ptr] {
			sum += v
		}
		//pointer math, circular buffer, yay!
		for _, v := range a.deltas[a.size+a.ptr-w:] {
			sum += v
		}
	} else { //sane, classic situation - window fell in the middle of the slice
		for _, v := range a.deltas[a.ptr-w : a.ptr] {
			sum += v
		}
	}
	return float32(sum) / float32(w), nil
}
