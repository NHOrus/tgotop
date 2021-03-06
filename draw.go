// tgotop project main.go

package main

import (
	"bytes"
	"fmt"
	//	spew "github.com/davecgh/go-spew/spew"
	ui "gopkg.in/gizak/termui.v1"
	//tm "github.com/nsf/termbox-go"
	"os"
	"os/signal"
	"text/tabwriter"
	"time"
)

const (
	mult  = 10
	dtick = time.Second / 2  //data update interval
	rtick = time.Second / 60 //redrawint interval
	atick = time.Second      //averaging interval
	ptick = atick / mult     //polling interval
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
	evt := ui.EventCh()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill, os.Interrupt)

	//initing tickers for polling and redrawing and updating data
	rtk := time.Tick(rtick)
	ptk := time.Tick(ptick)

	nd := new(netData)
	err = nd.Init(3*mult, ptk)

	if err != nil {
		panic(err)
	}

	gNet.Height = nd.size + 3

	go dataupd(gMem, gSwap, gNet, nd)

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, gMem, gSwap),
			ui.NewCol(6, 0, gNet)),
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
		case <-rtk:
			ui.Render(ui.Body)
		}
	}
}

func fillfmt(s string, u uint64, t uint64) string {
	return fmt.Sprintf("%v used: %v / %v", s, humanBytes(float32(u)), humanBytes(float32(t)))
}

func dealwithevents(e ui.Event) bool {
	if e.Type == ui.EventKey && (e.Ch == 'q' || e.Ch == 'Q') {
		return true
	}
	if e.Type == ui.EventKey && e.Key == ui.KeyCtrlC {
		return true
	}
	if e.Type == ui.EventResize {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
	}
	return false
}

// formatting and massaging net data into pretty tabs with units
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

// Thus function updates data each dtick
func dataupd(gMem, gSwap *ui.Gauge, gNet *ui.List, nd *netData) {

	var m memData
	for {
		err := m.Update()
		if err != nil {
			panic(err)
		}

		gMem.Percent = m.memPercent
		gMem.Border.Label = fillfmt("Memory", m.memUse, m.memTotal)

		gSwap.Percent = m.swapPercent
		gSwap.Border.Label = fillfmt("Swap", m.swapUse, m.swapTotal)

		gNet.Items = netf(nd)
		time.Sleep(dtick)
	}
}
