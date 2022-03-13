package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/crhntr/giocharts"
)

func main() {
	go func() {
		w := app.NewWindow()
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	var ops op.Ops

	c := make(chan []float64)
	go func() {
		for {
			array := make([]float64, 10)
			for i := range array {
				array[i] = rand.Float64() * 100
			}
			c <- array
			time.Sleep(time.Second * 5)
		}
	}()
	var data []float64
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				giocharts.Bar{
					Data: data,
				}.Layout(gtx)
				e.Frame(gtx.Ops)
			}
		case array := <-c:
			data = array
			w.Invalidate()
		}
	}
}
