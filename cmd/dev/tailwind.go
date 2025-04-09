package main

import (
	"os/exec"
	"sync"
)

func tailwind(wg *sync.WaitGroup, cmds chan<- *exec.Cmd) {
	runCommand("TAILWIND", colorCyan, []string{"mise", "x", "--", "bunx", "tailwindcss", "-i", "src/input.css", "-o", "assets/styles.css", "--watch"}, wg, cmds, true)
}
