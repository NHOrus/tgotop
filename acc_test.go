package main

import "testing"

func TestZeroDeltaAcc(t *testing.T) {
	defer func() {
		if err := recover(); err != ErrWrongSize {
			t.FailNow()
		}
	}()
	_ = NewDeltaAcc(0)
}

func TestNewDeltaAcc(t *testing.T) {

}
