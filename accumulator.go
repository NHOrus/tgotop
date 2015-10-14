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
	last uint64
	Acc
}

type Acc struct {
	size int
	ptr  int
	vals []int64
	sync.RWMutex
	full bool
}

//NewDeltaAcc returns delta-accumulator of given size

func NewAcc(s int) *Acc {
	if s < 1 {
		panic(ErrWrongSize)
	}
	return &Acc{
		size: s,
		//ptr increments or drops down to zero, it pains me to make "-1" special value
		//but other way is to add more fields in struct
		vals: make([]int64, s, s),
	}
}

func NewDeltaAcc(s int) *DeltaAcc {

	return &DeltaAcc{
		Acc: *NewAcc(s),
	}
}

//Purge cleans up the accumulator and returns it sparkingly clean
func (a *DeltaAcc) Purge() {
	a.Lock()
	a.ptr = 0
	a.last = 0
	for i := 0; i < a.size; i++ {
		a.vals[i] = 0
	}
	a.full = false
	a.Unlock()
}

//Purge cleans up the accumulator and returns it sparkingly clean
func (a *Acc) Purge() {
	a.Lock()
	a.ptr = 0
	for i := 0; i < a.size; i++ {
		a.vals[i] = 0
	}
	a.full = false
	a.Unlock()
}

//Code duplicate due to mutexes. Do not want to mess with mutexes

func (a *Acc) Push(v int64) {
	a.Lock()
	defer a.Unlock()
	a.vals[a.ptr] = v
}

//Push takes a value and adds difference between it and previous value into
//accumulator, with circular overwrite on full capacity
//First delta happens only after two values are pushed in
func (a *DeltaAcc) Push(v uint64) {
	var dlt int64 //temporary delta for calculations

	//if accumulator freshly initialized, we can't put delta in it
	//because we know not enough to calculate it
	//so pushed value gets saved and accumulator goes waiting for
	//next value, everything proceeding as it ought
	if a.ptr == -1 {
		a.last = v
		a.ptr = 0
		return
	}

	if a.last >= v {
		dlt = -int64(a.last - v)
	} else { //stupid unsigned math
		dlt = int64(v - a.last)
	}
	a.last = v
	if a.ptr == a.size-1 {
		a.full = true
		a.ptr = 0
	} else {
		a.ptr = a.ptr + 1
	}
	a.Acc.Push(dlt)
}

//Sum returns sum total of deltas in latest window of given size
func (a *Acc) Sum(w int) (sum int64, err error) {
	//if window is bigger than our accumulator, we fail horribly
	if w > a.size {
		return 0, ErrBigWindow
	}

	//Critical section - DeltaAcc should not chage while we are in it
	a.RLock()
	defer a.RUnlock()

	//if we don't have enough data to fill window, we fail, less horribly
	if !a.full && a.ptr < w {
		return 0, ErrSmallSample
	}
	//if we need to get bits from different ends of our slice, let it be so
	if a.ptr < w {
		for _, v := range a.vals[:a.ptr] {
			sum += v
		}
		//pointer math, circular buffer, yay!
		for _, v := range a.vals[a.size+a.ptr-w:] {
			sum += v
		}
	} else { //sane, classic situation - window fell in the middle of the slice
		for _, v := range a.vals[a.ptr-w : a.ptr] {
			sum += v
		}
	}
	return sum, nil
}

//Average returns sum average of deltas in latest window of given size
func (a *Acc) Average(w int) (avg float32, err error) {
	sum, err := a.Sum(w)
	avg = float32(sum) / float32(w)
	return
}
