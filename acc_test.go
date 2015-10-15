package main

import "testing"
import "math"

func TestZeroNewAcc(t *testing.T) {
	defer func() {
		if err := recover(); err != ErrWrongSize {
			t.FailNow()
		}
	}()
	_ = NewAcc(0)
}

func TestPushAcc(t *testing.T) {
	a := NewAcc(10)
	a.Push(100)
	if a.lastval() != 100 {
		t.Error("Not pushing right")
	}
}

func TestPushDeltaAcc(t *testing.T) {
	a := NewDeltaAcc(10)
	a.Push(100)
	if a.lastval() != 100 {
		t.Error("Not pushing right")
	}
	a.Push(200)
	if a.lastval() != 200 {
		t.Error("Not pushing right")
	}
	if a.lastdelta() != 100 {
		t.Error("Calculating delta wrong")
	}

	a.Push(100)
	if a.lastval() != 100 {
		t.Error("Not pushing right")
	}
	if a.lastdelta() != -100 {
		t.Error("Calculating negative delta wrong")
	}

}

func TestSum(t *testing.T) {
	a := NewAcc(10)

	_, err := a.Sum(15)
	if err != ErrBigWindow {
		t.Error("Dealing with request too big when shouldn't")
	}

	var tv int64 = 10 //test value
	a.Push(tv)
	b, err := a.Sum(1)
	if err != nil {
		t.Error("Erroring while everything is correct")
	}
	if b != tv {
		t.Error("Failing to get sum of 1 component, result is ", b, " last value is ", a.lastval())
	}

	_, err = a.Sum(5)
	if err != ErrSmallSample {
		t.Error("Succeeding in getting sum of underwritten accumulator")
	}
	for i := 0; i < 25; i++ {
		a.Push(tv)
	}
	if i, _ := a.Sum(10); i != 100 {
		t.Error("Something wrong with bigger summs")
	}

}

func TestAverage(t *testing.T) {
	a := NewAcc(10)
	for _, i := range []int64{2, 4, 6, 8} {
		a.Push(i)
	}
	if i, err := a.Average(4); err != nil && i != 10. {
		t.Error("Something wrong when taking an average")
	}
}

func TestDeltaAvg(t *testing.T) {
	a := NewAcc(6) //less than total write, so if second  half fails, something wrong on buffer wrapping
	for _, i := range []int64{2, 4, 6, 8} {
		a.Push(i)
	}
	if i, err := a.Average(4); err != nil && i != 2. {
		t.Error("Something wrong when taking an average")
	}

	for _, i := range []int64{8, 6, 4, 2} {
		a.Push(i)
	}
	if i, err := a.Average(4); err != nil && i != -2. {
		t.Error("Something wrong when taking an average of negatives")
	}

	if i, err := a.Average(0); err != nil && math.IsNaN(float64(i)) {
		t.Error("Got no stupid answer to stupid question of 0/0")
	}
}

func (a *Acc) lastval() int64 {
	return a.vals[a.head-1]
}

func (a *DeltaAcc) lastval() uint64 {
	return a.last
}

func (a *DeltaAcc) lastdelta() int64 {
	return a.Acc.lastval()
}
