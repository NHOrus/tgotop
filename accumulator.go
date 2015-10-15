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

//Acc is a write-only ring buffer of finite and static capacity, with methods that
//provide sum and average of last n values pushed inside
type Acc struct {
	size int
	head int //as Acc is  read-only, there is no need to know position of tail
	vals []int64
	sync.RWMutex
	full bool
}

//NewAcc returns accumulator of given size
func NewAcc(s int) *Acc {
	if s < 1 {
		panic(ErrWrongSize)
	}

	return &Acc{
		size: s,
		vals: make([]int64, s, s),
	}
}

//Purge cleans up the accumulator and returns it sparkingly clean
func (a *Acc) Purge() {
	a.Lock()
	a.head = 0
	for i := 0; i < a.size; i++ {
		a.vals[i] = 0
	}
	a.full = false
	a.Unlock()
}

//Push adds new value into Acc, overwriting old ones and wrapping
//around the ring as needed
func (a *Acc) Push(v int64) {
	a.Lock()
	defer a.Unlock()

	a.vals[a.head] = v

	if a.head == a.size-1 {
		a.full = true
		a.head = 0
	} else {
		a.head = a.head + 1
	}
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
	if !a.full && a.head < w {
		return 0, ErrSmallSample
	}
	//if we need to get bits from different ends of our slice, let it be so
	if a.head < w {
		for _, v := range a.vals[:a.head] {
			sum += v
		}
		//pointer math, circular buffer, yay!
		for _, v := range a.vals[a.size+a.head-w:] {
			sum += v
		}
	} else { //sane, classic situation - window is inside the slice
		for _, v := range a.vals[a.head-w : a.head] {
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

//DeltaAcc accumulates changes between pushed values.
//For DeltaAcc amount of remembered points is either size of underlying Acc or
//pushed amount minus one.
type DeltaAcc struct {
	last  uint64
	initd bool
	Acc
}

//NewDeltaAcc returns delta-accumulator of given size
func NewDeltaAcc(s int) *DeltaAcc {

	return &DeltaAcc{
		Acc: *NewAcc(s),
	}
}

//Purge cleans up the accumulator and returns it sparkingly clean
func (a *DeltaAcc) Purge() {
	a.Lock()
	a.head = 0
	a.last = 0
	for i := 0; i < a.size; i++ {
		a.vals[i] = 0
	}
	a.full = false
	a.Unlock()
}

//Push takes a value and adds difference between it and previous value into
//accumulator, with circular overwrite on full capacity
//First delta happens only after two values were pushed in
func (a *DeltaAcc) Push(v uint64) {
	defer func() { a.last = v }()

	if !a.initd {
		a.initd = true
		return
	}

	if a.last >= v {
		a.Acc.Push(-int64(a.last - v))
	} else { //stupid unsigned math
		a.Acc.Push(int64(v - a.last))
	}

}
