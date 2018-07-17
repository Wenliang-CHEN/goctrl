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
		createObjects(parameters)
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

func createObjects(yaml *Yaml) {
	objects, err := yaml.Get("objects").Array()
	if err != nil {
		panic(err)
	}

	for _, objectName := range objects {
		createObject(yaml, objectName.(string))
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
