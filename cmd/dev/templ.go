package main

import (
	"os/exec"
	"sync"
)

func templ(wg *sync.WaitGroup, cmds chan<- *exec.Cmd) {
	runCommand("TEMPL", colorGreen, []string{"air", "-c", ".air.templ.toml"}, wg, cmds, false)
}
