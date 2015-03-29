// tgotop project main.go

package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	tm "github.com/nsf/termbox-go"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//DIV is size of divider, in this case - MiB
const DIV uint64 = 1024 * 1024

//DIVname is a name of unit that we ends up after dividing bytes by DIV
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
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

	go func() {
		for {
			var m memData
			err := m.Update()
			if err != nil {
				panic(err)
			}

			gMem.Percent = m.memPercent
			gMem.Border.Label = fillfmt("Memory", m.memUse, m.memTotal)

			gSwap.Percent = m.swapPercent
			gSwap.Border.Label = fillfmt("Swap", m.swapUse, m.swapTotal)
			time.Sleep(time.Second / 2)
		}
	}()

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, gMem, gSwap)),
		ui.NewRow(ui.NewCol(12, 0, qMess)))

	ui.Body.Align()

	for {
		select {
		case e := <-evt:
			if dealwithevents(e) {
				return
			}
		case <-sig:
			return
		default:
			ui.Render(ui.Body)
		}
	}
}

func fillfmt(s string, u uint64, t uint64) string {
	return fmt.Sprintf("%v used: %d / %d %v", s, u/DIV, t/DIV, DIVname)
}

func dealwithevents(e tm.Event) bool {
	if e.Type == tm.EventKey && e.Ch == 'q' {
		return true
	}
	if e.Type == tm.EventKey && e.Key == tm.KeyCtrlC {
		return true
	}
	if e.Type == tm.EventResize {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
	}
	return false
}
