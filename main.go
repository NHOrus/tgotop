// tgotop project main.go

package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	tm "github.com/nsf/termbox-go"
	"time"
)

//DIV is size of divider, in this case - MiB
const DIV uint64 = 1024 * 1024

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	qMess := ui.NewPar(":PRESS q TO QUIT")
	qMess.Height = 3

	gSwap := ui.NewGauge()
	gSwap.Height = 3

	gMem := ui.NewGauge()
	gMem.Height = 3

	//getting ready to close stuff on command
	evt := make(chan tm.Event)
	go func() {
		for {
			evt <- tm.PollEvent()
		}
	}()

	draw := func() {
		var m memData
		err := m.Update()
		if err != nil {
			panic(err)
		}

		gMem.Percent = m.memPercent
		gMem.Border.Label = fmt.Sprintf("Memory used: %d / %d MiB", m.memUse/DIV, m.memTotal/DIV)

		gSwap.Percent = m.swapPercent
		gSwap.Border.Label = fmt.Sprintf("Swap used: %d / %d MiB", m.swapUse/DIV, m.swapTotal/DIV)
		ui.Render(ui.Body)
	}

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, gMem, gSwap)),
		ui.NewRow(ui.NewCol(12, 0, qMess)))

	ui.Body.Align()

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
			go draw()
			time.Sleep(time.Second / 2)
		}
	}
}
