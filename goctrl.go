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

// TODO: handle error if i have time
func runOsCommand(rawCmd string, args ...string) {
	var stdoutBuf bytes.Buffer
	cmd := exec.Command(rawCmd, args...)

	stdoutIn, _ := cmd.StdoutPipe()
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Start()

	go func() {
		io.Copy(stdout, stdoutIn)
	}()

	cmd.Wait()
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
