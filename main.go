package main

import (
	"os"

	"github.com/CS80-Team/gshell/pkg/gshell"
	"github.com/chzyer/readline"
)

func main() {
	stdin, stdinW := readline.NewFillableStdin(os.Stdin)

	shell := gshell.NewShell(
		stdin,
		stdinW,
		os.Stdout,
		os.Stdout,
		gshell.SHELL_PROMPT,
		".shell_history",
		gshell.NewLogger("shell.log"),
	)

	shell.Run("Debug: Shell started")
}
