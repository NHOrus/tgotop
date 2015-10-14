package main

import "testing"

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
}
