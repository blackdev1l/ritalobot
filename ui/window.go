package main

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
	file, _ := os.Open("../logo")
	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		tbprint(0, y, coldef, coldef, scanner.Text())
		y++

	}
}

func main() {
	var current string
	var curev t.Event
	data := make([]byte, 0, 64)
	err := t.Init()
	if err != nil {
		panic(err)
	}
	defer t.Close()
	t.Clear(coldef, coldef)

	showLogo()
	t.Flush()
mainloop:
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
				break mainloop
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
