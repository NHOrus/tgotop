// accaverager
package main

import (
	"errors"
	"sync"
)

var (
	//ErrBigWindow happens when we try to get more datapoints to average what
	//we can put in our accumulator. Good place to panic
	ErrBigWindow = errors.New("averaging window is bigger than accumulator")

	//ErrSmallSample happens when we just began and got not enough datapoints
	//to fill our window. Normal working situation.
	ErrSmallSample = errors.New("not enough data saved")

	//ErrWrongSize happens when signed int was expected to be positive but
	//ended up zero or negative. Sad, sad situation. Panic.
	ErrWrongSize = errors.New("size can not be less than one")
)

//DeltaAcc accumulates changes between pushed values
type DeltaAcc struct {
	size       int
	ptr        int
	last       uint64
	deltas     []int64
	full       bool
	sync.Mutex //If stuff gets pushed in while we are getting average, BAD THINGS
}

//NewDeltaAcc returns delta-accumulator of given size
func NewDeltaAcc(s int) (*DeltaAcc, error) {

	if s < 1 {
		return nil, ErrWrongSize
	}
	return &DeltaAcc{
		size: s,
		//ptr increments or drops down to zero, it pains me to make "-1"special value
		//but other way is to add more fields in struct
		ptr:    -1,
		deltas: make([]int64, s, s),
	}, nil
}

//Purge cleans up the accumulator and returns it sparkingly clean
func (a *DeltaAcc) Purge() {
	a.Lock()
	a.ptr = -1
	a.last = 0
	for i := 0; i < a.size; i++ {
		a.deltas[i] = 0
	}
	a.full = false
	a.Unlock()
}

//Push takes a value and adds difference between it and previous value into accumulator
// going in circles when accumulator is full.
//First delta happens only after two values are pushed in
func (a *DeltaAcc) Push(v uint64) {
	var dlt int64 //temporary delta for calculations
	a.Lock()
	//if
	if a.ptr == -1 {
		a.last = v
		return
	}
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
	a.Unlock()
}

//Average returns average of deltas in latest window of given size
func (a *DeltaAcc) Average(w int) (avg float32, err error) {
	a.Lock()
	var sum int64 //here we add every delta in the window

	//if window is bigger than our accumulator, we fail horribly
	if w > a.size {
		return 0, ErrBigWindow
	}
	//if we don't have enough data to fill window, we fail, less horribly
	if !a.full && a.ptr < w {
		return 0, ErrSmallSample
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
	a.Unlock()
	return float32(sum) / float32(w), nil
}
