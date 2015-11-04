// tgotop project main.go

package main

import (
	"bytes"
	"fmt"

	"github.com/davecgh/go-spew/spew"

	ui "gopkg.in/gizak/termui.v2"
	//tm "github.com/nsf/termbox-go"
	"text/tabwriter"
	"time"
)

const (
	mult  = 10
	rtick = time.Second / 60
	atick = time.Second  //averaging interval
	ptick = atick / mult //polling interval
)

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

	gNet := ui.NewList()

	//getting ready to close stuff on command
	//	sig := make(chan os.Signal)
	//	signal.Notify(sig, os.Interrupt, os.Kill)

	nd := new(netData)
	err = nd.Init(3*mult, ptick)

	if err != nil {
		panic(err)
	}

	gNet.Height = nd.size + 3
	var m memData

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, gMem, gSwap),
			ui.NewCol(6, 0, gNet)),
		ui.NewRow(ui.NewCol(12, 0, qMess)))

	ui.Body.Align()
	ui.Render(ui.Body)

	ui.Merge("timer/update", ui.NewTimerCh(ptick))
	ui.Merge("timer/render", ui.NewTimerCh(rtick))

	ui.Handle("/timer/"+ptick.String(), func(ui.Event) {
		err := m.Update()
		if err != nil {
			panic(err)
		}
		gMem.Percent = m.memPercent
		gMem.BorderLabel = fillfmt("Memory", m.memUse, m.memTotal)

		gSwap.Percent = m.swapPercent
		gSwap.BorderLabel = fillfmt("Swap", m.swapUse, m.swapTotal)

		gNet.Items = netf(nd)

		ui.Render(ui.Body)
	})

	ui.Handle("/timer/"+rtick.String(), func(ui.Event) {
		ui.Render(ui.Body)
	})

	ui.Handle("/sys/kdb/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kdb/Q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kdb/C-c", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		if ui.TermWidth() > 20 {
			ui.Body.Width = ui.TermWidth()
			ui.Body.Align()
		}
		ui.Render(ui.Body)
	})

	spew.Dump(ui.DefaultEvtStream)

	ui.Loop()
}

func fillfmt(s string, u uint64, t uint64) string {
	return fmt.Sprintf("%v used: %v / %v", s, humanBytes(float32(u)), humanBytes(float32(t)))
}

func netf(nd *netData) []string {
	strings := make([]string, 0, nd.size+1)
	var b bytes.Buffer
	tb := tabwriter.NewWriter(&b, 10, 8, 0, ' ', tabwriter.AlignRight)
	fmt.Fprintln(tb, "IFace\t Down\t Up\t")
	for i := 0; i < nd.size; i++ {
		fmt.Fprintf(tb, "%v:\t %s/s\t %s/s\t\n",
			nd.name[i],
			humanBytes(nd.GetD(i, mult)),
			humanBytes(nd.GetU(i, mult)))
	}

	err := tb.Flush()

	if err != nil {
		panic(err)
	}
	for i := 0; i <= nd.size; i++ {
		ts, err := b.ReadString('\n')
		if err != nil {
			panic(err)
		}
		strings = append(strings, ts)
	}
	return strings
}
