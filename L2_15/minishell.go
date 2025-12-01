package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt) // ловим Ctrl+C

	for {
		fmt.Print("> ")

		inputChan := make(chan string)
		errChan := make(chan error)

		// Чтение ввода в отдельной горутине
		go func() {
			line, err := reader.ReadString('\n')
			if err != nil {
				errChan <- err
			} else {
				inputChan <- line
			}
		}()

		select {
		case <-sigChan:
			// Ctrl+C во время ввода: просто выводим новую строку
			fmt.Println()
			continue
		case err := <-errChan:
			if err == io.EOF {
				fmt.Println("\nExiting shell.")
				return
			}
			fmt.Println("Error reading input:", err)
			continue
		case line := <-inputChan:
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// Здесь запускаем выполнение команды
			executePipelineOrCommand(line, sigChan)
		}
	}
}

// ================= Conditional execution && and ||
type condPart struct {
	cmd string
	op  string
}

func splitConditional(line string) []condPart {
	var parts []condPart
	var buf strings.Builder
	i := 0
	for i < len(line) {
		if strings.HasPrefix(line[i:], "&&") {
			parts = append(parts, condPart{cmd: buf.String(), op: "&&"})
			buf.Reset()
			i += 2
		} else if strings.HasPrefix(line[i:], "||") {
			parts = append(parts, condPart{cmd: buf.String(), op: "||"})
			buf.Reset()
			i += 2
		} else {
			buf.WriteByte(line[i])
			i++
		}
	}
	parts = append(parts, condPart{cmd: buf.String(), op: ""})
	return parts
}

func executeConditional(line string, sigChan chan os.Signal) {
	parts := splitConditional(line)
	lastExit := 0
	for i, p := range parts {
		cmd := strings.TrimSpace(p.cmd)
		if cmd == "" {
			continue
		}

		if i > 0 {
			switch parts[i-1].op {
			case "&&":
				if lastExit != 0 {
					continue
				}
			case "||":
				if lastExit == 0 {
					continue
				}
			}
		}

		lastExit = executePipelineOrCommand(cmd, sigChan)
	}
}

// ================= Variable substitution
func expandEnv(line string) string {
	fields := strings.Fields(line)
	for i, f := range fields {
		if strings.HasPrefix(f, "$") {
			fields[i] = os.Getenv(f[1:])
		}
	}
	return strings.Join(fields, " ")
}

// ================= Pipeline execution
func executePipelineOrCommand(line string, sigChan chan os.Signal) int {
	line = expandEnv(line)
	commands := strings.Split(line, "|")
	if len(commands) > 1 {
		return executePipeline(commands, sigChan)
	} else {
		return executeCommand(line, sigChan)
	}
}

func executePipeline(commands []string, sigChan chan os.Signal) int {
	n := len(commands)
	var prevReader *os.File
	var wg sync.WaitGroup
	exitCode := 0

	for i, cmdStr := range commands {
		cmdStr = strings.TrimSpace(cmdStr)
		args := strings.Fields(cmdStr)

		// Последняя команда выводит в stdout
		var out *os.File
		var r, w *os.File
		if i < n-1 {
			r, w, _ = os.Pipe()
			out = w
		} else {
			out = os.Stdout
		}

		// Сохраняем reader для следующей команды
		currReader := r

		wg.Add(1)
		go func(args []string, in *os.File, out *os.File, lastReader *os.File) {
			defer wg.Done()
			defer func() {
				if out != os.Stdout {
					out.Close()
				}
				if lastReader != nil {
					lastReader.Close()
				}
			}()

			if isBuiltin(args[0]) {
				// встроенная команда
				executeBuiltin(args, "", "", in, out)
			} else {
				// внешняя команда
				cmd := exec.Command(args[0], args[1:]...)
				if in != nil {
					cmd.Stdin = in
				}
				cmd.Stdout = out
				cmd.Stderr = os.Stderr
				cmd.Run()
			}
		}(args, prevReader, out, prevReader)

		// Следующая команда читает из r
		prevReader = currReader
	}

	wg.Wait()
	return exitCode
}

// ================= Parse redirections
func parseRedirection(line string) (args []string, stdinFile, stdoutFile string) {
	parts := strings.Fields(line)
	args = []string{}
	i := 0
	for i < len(parts) {
		if parts[i] == ">" && i+1 < len(parts) {
			stdoutFile = parts[i+1]
			i += 2
		} else if parts[i] == "<" && i+1 < len(parts) {
			stdinFile = parts[i+1]
			i += 2
		} else {
			args = append(args, parts[i])
			i++
		}
	}
	return
}

// ================= Check builtin
func isBuiltin(cmd string) bool {
	switch cmd {
	case "cd", "pwd", "echo", "kill", "ps":
		return true
	}
	return false
}

// ================= Run a command (builtin or external)
func runCommand(args []string, stdinFile, stdoutFile string, input *os.File, output *os.File, sigChan chan os.Signal) int {
	if len(args) == 0 {
		return 0
	}

	if isBuiltin(args[0]) {
		return executeBuiltin(args, stdinFile, stdoutFile, input, output)
	}

	cmd := exec.Command(args[0], args[1:]...)

	// Настраиваем stdin
	if stdinFile != "" {
		f, err := os.Open(stdinFile)
		if err != nil {
			fmt.Println("Failed to open input file:", err)
			return 1
		}
		defer f.Close()
		cmd.Stdin = f
	} else if input != nil {
		cmd.Stdin = input
	} else {
		cmd.Stdin = os.Stdin
	}

	// Настраиваем stdout
	if stdoutFile != "" {
		f, err := os.Create(stdoutFile)
		if err != nil {
			fmt.Println("Failed to open output file:", err)
			return 1
		}
		defer f.Close()
		cmd.Stdout = f
	} else if output != nil {
		cmd.Stdout = output
	} else {
		cmd.Stdout = os.Stdout
	}

	cmd.Stderr = os.Stderr

	// Канал для завершения команды
	done := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- fmt.Errorf("command panicked: %v", r)
			}
		}()
		done <- cmd.Run()
	}()

	// Ждём либо сигнал Ctrl+C, либо завершение команды
	select {
	case <-sigChan:
		if cmd.Process != nil {
			// Windows: Kill() единственный способ прервать команду
			cmd.Process.Kill()
		}
		fmt.Println("\nCommand interrupted")
		return 1
	case err := <-done:
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return exitErr.ExitCode()
			}
			fmt.Println("Error running command:", err)
			return 1
		}
	}

	return 0
}

// ================= Execute builtin (supports pipes)
func executeBuiltin(args []string, stdinFile, stdoutFile string, input *os.File, output *os.File) int {
	// Настроим корректный вывод
	var out *os.File
	if output != nil {
		out = output
	} else {
		out = os.Stdout
	}

	// Настроим корректный ввод
	var in *os.File
	if input != nil {
		in = input
	} else if stdinFile != "" {
		f, err := os.Open(stdinFile)
		if err != nil {
			fmt.Println("Error opening input file:", err)
			return 1
		}
		defer f.Close()
		in = f
	} else {
		in = os.Stdin
	}

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			fmt.Println("cd: missing argument")
			return 1
		}
		if err := os.Chdir(args[1]); err != nil {
			fmt.Println("cd error:", err)
			return 1
		}
	case "pwd":
		writeOutput(getCwd()+"\n", out)
	case "echo":
		writeOutput(strings.Join(args[1:], " ")+"\n", out)
	case "kill":
		if len(args) < 2 {
			fmt.Println("kill: missing PID")
			return 1
		}
		pid := atoi(args[1])
		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Println("kill error:", err)
			return 1
		}
		process.Kill()
	case "ps":
		cmd := exec.Command("ps", "aux")
		cmd.Stdout = out
		cmd.Stderr = os.Stderr
		cmd.Stdin = in
		if err := cmd.Run(); err != nil {
			fmt.Println("ps error:", err)
			return 1
		}
	default:
		fmt.Println("Unknown builtin:", args[0])
		return 1
	}

	return 0
}

func getCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

// ================= Helper to write output
func writeOutput(s string, output *os.File) {
	if output != nil {
		output.WriteString(s)
	} else {
		fmt.Print(s)
	}
}

// ================= Command execution helper
func executeCommand(line string, sigChan chan os.Signal) int {
	args, stdinFile, stdoutFile := parseRedirection(line)
	if len(args) == 0 {
		return 0
	}
	if isBuiltin(args[0]) {
		return executeBuiltin(args, stdinFile, stdoutFile, nil, nil)
	}
	return runCommand(args, stdinFile, stdoutFile, nil, nil, sigChan)
}

// ================= String to int
func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
