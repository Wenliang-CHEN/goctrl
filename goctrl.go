package main

import (
	"bytes"
	"fmt"
	"github.com/smallfish/simpleyaml"
	"io"
	"io/ioutil"
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
		createKubernetesObjects()
	case "stop":
		stopMinikube()
	default:
		panic(errInvalidCommand)
	}
}

func startMinikube() {
	runOsCommand("minikube", "start")
}

//TODO: add separate create and delete for each object
func createKubernetesObjects() {
	parameters := parseYaml("parameters.yaml")

	configBasePath, err := parameters.Get("config-path").String()
	objects, err := parameters.Get("objects").Array()

	if err != nil {
		panic(err)
	}

	for _, objectName := range objects {
		fullPath := configBasePath + objectName.(string) + ".yaml"
		runOsCommand("kubectl", "apply", "-f", fullPath)
	}
}

func stopMinikube() {
	runOsCommand("minikube", "stop")
}

func parseYaml(filePath string) *simpleyaml.Yaml {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic("unable to read file: " + filePath)
	}

	yaml, err := simpleyaml.NewYaml(content)
	if err != nil {
		panic(err)
	}

	return yaml
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

// TODO: Add help text
func printHelpText() {
	fmt.Println("this is help text")
}
