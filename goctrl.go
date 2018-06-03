package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

const errInvalidCommand = "Invalid command"

func main() {
	defer handleError()

	ctrlCmd := os.Args[1]
	switch ctrlCmd {
	case "start":
		startMinikube()
	case "stop":
		stopMinikube()
	default:
		panic(errInvalidCommand)
	}
}

func startMinikube() {
	fmt.Println("starting minikube")
	runOsCommand("minikube", "start")
}

func stopMinikube() {
	fmt.Println("stoping minikube")
}

func runOsCommand(rawCmd string, args ...string) {
	var stdoutBuf bytes.Buffer
	cmd := exec.Command(rawCmd, args...)

	stdoutIn, _ := cmd.StdoutPipe()
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	var errStdout error
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
	if errStdout != nil {
		panic("failed to capture stdout")
	}
	outStr := string(stdoutBuf.Bytes())
	fmt.Printf("\nout:\n%s", outStr)
}

func handleError() {
	err := recover()
	if err == nil {
		return
	}

	switch err {
	case errInvalidCommand:
		fmt.Println(err)
		printHelpText()
	default:
		fmt.Println(err)
	}
}

// TODO
func printHelpText() {
	fmt.Println("this is help text")
}
