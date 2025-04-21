package gshell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (sh *Shell) registerBuiltInCommands() {
	sh.RegisterCommand(
		NewCommand(
			"exit",
			"Exit the shell",
			"exit",
			[]Argument{},
			[]string{},
			func(s *Shell, args []string) (Status, error) {
				return EXIT, nil
			},
			func(args []string) (bool, error) {
				return true, nil
			},
		),
	)

	sh.RegisterCommand(
		NewCommand(
			"help",
			"List all available commands",
			"help",
			[]Argument{},
			[]string{},
			func(s *Shell, args []string) (Status, error) {
				for _, cmd := range sh.GetCommands() {
					sh.WriteColored(COLOR_YELLOW, cmd.Name+": ")
					sh.Write(cmd.Description + "\n")
					sh.Write("    Aliases: ")
					if len(cmd.Aliases) > 0 {
						sh.Write(strings.Join(cmd.Aliases, ", ") + "\n")
					} else {
						sh.Write("No aliases found.\n")
					}
					sh.WriteColored(COLOR_CYAN, "    Usage: "+cmd.Usage+"\n\n")
				}
				return OK, nil
			},
			func(args []string) (bool, error) {
				return true, nil
			},
		),
	)

	sh.RegisterCommand(
		NewCommand(
			"clear",
			"Clear the screen",
			"clear",
			[]Argument{},
			[]string{"cls"},
			func(s *Shell, args []string) (Status, error) {
				s.clearScreen()
				return OK, nil
			},
			func(args []string) (bool, error) {
				return true, nil
			},
		),
	)

	sh.RegisterCommand(
		NewCommand(
			"history",
			"Display the shell history",
			"history",
			[]Argument{},
			[]string{"hist"},
			func(s *Shell, args []string) (Status, error) {
				file, err := os.ReadFile(sh.historyFile)
				if err != nil {
					return FAIL, fmt.Errorf("error reading history file: %s", err)
				}

				sh.Write(string(file))
				return OK, nil
			},
			func(args []string) (bool, error) {
				return true, nil
			},
		),
	)

	sh.RegisterCommand(
		NewCommand(
			"alias",
			"Create an alias for a command",
			"alias <alias> <command>",
			[]Argument{
				{
					Name:        "Alias",
					Description: "The alias to create",
					Required:    true,
					Type:        "string",
					Default:     "",
				},
				{
					Name:        "Command",
					Description: "The command to create an alias for",
					Required:    true,
					Type:        "string",
					Default:     "",
				},
			},
			[]string{},
			func(s *Shell, args []string) (Status, error) {
				cmd, _ := sh.findCommandByNameOrAlias(args[1])
				sh.addAlias(args[0], cmd.Name)
				return OK, nil
			},
			func(args []string) (bool, error) {
				if len(args) != 2 {
					return false, fmt.Errorf("invalid number of arguments")
				}

				if _, ok := sh.findCommandByNameOrAlias(args[1]); !ok {
					return false, fmt.Errorf("command or alias not found")
				}
				return true, nil
			},
		),
	)

	sh.RegisterCommand(
		NewCommand(
			"exec",
			"Execute a an external command",
			"exec <command>",
			[]Argument{
				{
					Name:        "Execute command",
					Description: "The command to execute",
					Required:    true,
					Type:        "string",
					Default:     "",
				},
			},
			[]string{},
			func(s *Shell, args []string) (Status, error) {
				var ar []string

				for i := 0; i < len(args); i++ {
					s := args[i]
					if strings.HasPrefix(args[i], "\"") {
						s = strings.TrimPrefix(s, "\"")
						i++
						for i < len(args) && !strings.HasSuffix(args[i], "\"") {
							s += " " + args[i]
							i++
						}
						if i < len(args) {
							s += " " + strings.TrimSuffix(args[i], "\"")
						}
					}
					ar = append(ar, s)
				}

				cmd := exec.Command(ar[0], ar[1:]...)

				cmd.Stdout = sh.outStream
				cmd.Stderr = sh.errStream
				cmd.Stdin = sh.inStream

				err := cmd.Run()

				if err != nil {
					return FAIL, fmt.Errorf("error executing command: %s", err)
				}

				return OK, nil
			},
			func(args []string) (bool, error) {
				if len(args) < 1 {
					return false, fmt.Errorf("no command provided")
				}

				return true, nil
			},
		),
	)

	sh.RegisterCommand(
		NewCommand(
			"run",
			"Run a script",
			"run <script_path.shell>",
			[]Argument{
				{
					Name:        "Script Path",
					Description: "The path to the script to run",
					Required:    true,
					Type:        "string",
					Default:     "",
				},
			},
			[]string{},
			func(s *Shell, args []string) (Status, error) {
				file, err := os.ReadFile(args[0])
				if err != nil {
					return FAIL, fmt.Errorf("error reading script file: %s", err)
				}

				lines := strings.Split(string(file), "\n")
				for _, line := range lines {
					sh.execute(&line)
				}

				return OK, nil
			},
			func(args []string) (bool, error) {
				if len(args) != 1 {
					return false, fmt.Errorf("invalid number of arguments")
				}

				if info, err := os.Stat(args[0]); os.IsNotExist(err) {
					return false, fmt.Errorf("file does not exist")
				} else if info.IsDir() || filepath.Ext(args[0]) != ".shell" {
					return false, fmt.Errorf("invalid file type")
				}

				return true, nil
			},
		),
	)

	sh.RegisterCommand(
		NewCommand(
			"sleep",
			"idle the shell for a specified time",
			"sleep <seconds>",
			[]Argument{
				{
					Name:        "Seconds",
					Description: "The number of seconds to idle",
					Required:    true,
					Type:        "int",
					Default:     "",
				},
			},
			[]string{},
			func(s *Shell, args []string) (Status, error) {
				sec, _ := strconv.Atoi(args[0])
				time.Sleep(time.Duration(sec) * time.Second)
				return OK, nil
			},
			func(args []string) (bool, error) {
				if len(args) != 1 {
					return false, fmt.Errorf("invalid number of arguments")
				}

				sec, err := strconv.Atoi(args[0])
				if err != nil {
					return false, fmt.Errorf("invalid number of seconds: %s", err)
				}

				if sec < 0 {
					return false, fmt.Errorf("number of seconds must be positive")
				}

				return true, nil
			},
		),
	)
}
