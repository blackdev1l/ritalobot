package ui

import (
	"bufio"
	"fmt"
	t "github.com/nsf/termbox-go"
	"os"
)

const coldef = t.ColorDefault

/*
 function used to print a string on backbuffer
 input:
 x y coordinates, foreground and background color, msg to print
*/
func tbprint(x, y int, fg, bg t.Attribute, msg string) {
	for _, c := range msg {
		t.SetCell(x, y, c, fg, bg)
		x++
	}
}

func showLogo() {
	file, _ := os.Open("logo") // weird shit happens here
	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		tbprint(0, y, coldef, coldef, scanner.Text())
		y++

	}
}

func redraw() {
	_, h := t.Size()
	showLogo()
	tbprint(0, h-1, coldef, coldef, "press q to exit")
	t.Flush()

}

// This fucntion handles keyboard keys
// for now it just listen for the "q" key
// and exit the program after the release
func pollKeyboard() {
	var current string
	var curev t.Event
	data := make([]byte, 0, 64)
	for {
		if cap(data)-len(data) < 32 {
			newdata := make([]byte, len(data), len(data)+32)
			copy(newdata, data)
			data = newdata
		}
		beg := len(data)
		d := data[beg : beg+32]

		switch ev := t.PollRawEvent(d); ev.Type {
		case t.EventRaw:
			data = data[:beg+ev.N]
			current = fmt.Sprintf("%q", data)
			if current == `"q"` {
				os.Exit(0)
			}

			for {
				ev := t.ParseEvent(data)
				if ev.N == 0 {
					break
				}
				curev = ev
				copy(data, data[curev.N:])
				data = data[:len(data)-curev.N]
			}
		case t.EventError:
			panic(ev.Err)
		}
	}

}

// window loop, this should get events from the main process
func Show(ch <-chan int) {

	err := t.Init()
	if err != nil {
		panic(err)
	}
	defer t.Close()
	t.Clear(coldef, coldef)
	go pollKeyboard()

	t.Flush()
	for {
		redraw()
		select {
		case number := <-ch:
			if number == 0 {
				tbprint(0, 8, coldef, coldef, "redis ")
				tbprint(6, 8, t.ColorGreen, coldef, "OK")
			}
		}
	}
}
