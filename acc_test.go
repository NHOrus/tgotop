package main

import "testing"

func TestZeroDeltaAcc(t *testing.T) {
	defer func() {
		if err := recover(); err != ErrWrongSize {
			t.FailNow()
		}
	}()
	_ = NewAcc(0)
}

func TestNewAcc(t *testing.T) {

}
