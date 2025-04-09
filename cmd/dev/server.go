package main

import (
	"os/exec"
	"sync"
)

func server(wg *sync.WaitGroup, cmds chan<- *exec.Cmd) {
	runCommand("SERVER", colorPurple, []string{"air", "-c", ".air.server.toml"}, wg, cmds, true)
}
