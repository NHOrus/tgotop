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
const DIVname = "MiB"

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
		gMem.Border.Label = fillfmt("Memory", m.memUse, m.memTotal)

		gSwap.Percent = m.swapPercent
		gSwap.Border.Label = fillfmt("Swap", m.swapUse, m.swapTotal)
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

func fillfmt(s string, u uint64, t uint64) string {
	return fmt.Sprintf("%v used: %d / %d %v", s, u/DIV, t/DIV, DIVname)
}
