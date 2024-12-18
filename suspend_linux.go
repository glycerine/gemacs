package main

import (
	"github.com/glycerine/tcell_old_hacked_up/termbox"
	"syscall"
)

func suspend(g *gemacs) {
	// finalize termbox
	termbox.Close()

	// suspend the process
	pid := syscall.Getpid()
	tid := syscall.Gettid()
	err := syscall.Tgkill(pid, tid, syscall.SIGSTOP)
	if err != nil {
		panic(err)
	}

	// reset the state so we can get back to work again
	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputAlt)
	g.resize()
}
