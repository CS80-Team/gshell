package main

import (
	"os"

	"github.com/chzyer/readline"
)

func main() {
	stdin, stdinW := readline.NewFillableStdin(os.Stdin)

	shell := NewShell(
		stdin,
		stdinW,
		os.Stdout,
		os.Stdout,
		SHELL_PROMPT,
		".shell_history",
		NewLogger("shell.log"),
	)

	shell.Run("Debug: Shell started")
}
