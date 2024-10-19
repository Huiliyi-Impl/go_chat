package utils

import "github.com/nsf/termbox-go"

func ListenKeyboardEsc(c chan bool) {
	err := termbox.Init()
	if err != nil {
		c <- false
	}
	defer termbox.Close()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				c <- true
			}
		default:
			c <- false
		}
	}
}
