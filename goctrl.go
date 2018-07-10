package main

import (
	. "github.com/smallfish/simpleyaml"
	errors "goCtrl/errhandle"
	oscmd "goCtrl/oscmd"
	parser "goCtrl/parser"
	"os"
)

func main() {
	defer errors.HandleError()

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
		panic(errors.InvalidCommand)
	}
}

func startMinikube() {
	oscmd.Run("minikube", "start")
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
		oscmd.Run("kubectl", "apply", "-f", fullPath)
	}
}

func createObject(yaml *Yaml, name string) {
	oscmd.Run("kubectl", "apply", "-f", getBasePath(yaml)+name)
}

func deleteObject(yaml *Yaml, name string) {
	oscmd.Run("kubectl", "delete", "-f", getBasePath(yaml)+name)
}

func getBasePath(yaml *Yaml) string {
	configBasePath, err := yaml.Get("config-path").String()
	if err != nil {
		panic(err)
	}
	return configBasePath
}

func stopMinikube() {
	oscmd.Run("minikube", "stop")
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
>>>>>>> Add vendor handling and refactor package structure
}
