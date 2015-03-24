// tgotop project main.go
package main

import (
	//	proc "github.com/cespare/goproc"
	ui "github.com/gizak/termui"
	tm "github.com/nsf/termbox-go"
	// "time"
)

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	//getting ready to close stuff on command
	evt := make(chan tm.Event)
	go func() {
		for {
			evt <- tm.PollEvent()
		}
	}()

	for {
		select {
		case e := <-evt:
			if e.Type == tm.EventKey && e.Ch == 'q' {
				return
			}
			if e.Type == tm.EventResize {
				ui.Body.Width = ui.TermWidth()
				ui.Body.Align()
			}
		default:
			//	draw(i)
		}
	}
}
