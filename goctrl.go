package main

import (
	"bytes"
	"fmt"
	YamlUtil "github.com/smallfish/simpleyaml"
	parser "goCtrl/parser/yaml"
	"io"
	"os"
	"os/exec"
)

const errInvalidCommand = "Invalid command"

func main() {
	defer handleError()

	parameters := parser.Parse("config/parameters.yaml")
	ctrlCmd := os.Args[1]

	switch ctrlCmd {
	case "start":
		startMinikube()
		createKubernetesObjects()
	case "create":
		if len(os.Args) < 3 {
			panic("Please enter object name")
		}

		objectName := os.Args[2]
		createObject(parameters, objectName)
	case "delete":
		if len(os.Args) < 3 {
			panic("Please enter object name")
		}

		objectName := os.Args[2]
		deleteObject(parameters, objectName)
	case "stop":
		stopMinikube()
	default:
		panic(errInvalidCommand)
	}
}

func startMinikube() {
	runOsCommand("minikube", "start")
}

//TODO: Refactor this func
func createKubernetesObjects() {
	parameters := parser.Parse("config/parameters.yaml")

	configBasePath, err := parameters.Get("config-path").String()
	objects, err := parameters.Get("objects").Array()

	if err != nil {
		panic(err)
	}

	for _, objectName := range objects {
		fullPath := configBasePath + objectName.(string)
		runOsCommand("kubectl", "apply", "-f", fullPath)
	}
}

func createObject(yaml *YamlUtil.Yaml, name string) {
	runOsCommand("kubectl", "apply", "-f", getBasePath(yaml)+name)
}

func deleteObject(yaml *YamlUtil.Yaml, name string) {
	runOsCommand("kubectl", "delete", "-f", getBasePath(yaml)+name)
}

func getBasePath(yaml *YamlUtil.Yaml) string {
	configBasePath, err := yaml.Get("config-path").String()
	if err != nil {
		panic(err)
	}
	return configBasePath
}

func stopMinikube() {
	runOsCommand("minikube", "stop")
}

// TODO: to capture error message
func runOsCommand(rawCmd string, args ...string) {
	cmd := exec.Command(rawCmd, args...)

	var stdoutBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)

	var stderrBuf bytes.Buffer
	stdErrIn, _ := cmd.StderrPipe()
	stdErr := io.MultiWriter(os.Stderr, &stderrBuf)

	cmd.Start()

	go func() {
		io.Copy(stdout, stdoutIn)
	}()

	go func() {
		io.Copy(stdErr, stdErrIn)
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
