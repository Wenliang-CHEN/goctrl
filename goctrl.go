package main

import (
	"fmt"
	"os"
	"os.exec"
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
}

func stopMinikube() {
	fmt.Println("stoping minikube")
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
		fmt.Println("Unhandled Error")
	}
}

# TODO
func printHelpText() {
	fmt.Println("this is help text")
}
