package main

import (
	. "github.com/smallfish/simpleyaml"
	errors "goCtrl/errhandle"
	oscmd "goCtrl/oscmd"
	parser "goCtrl/parser"
	slice "goCtrl/utils/slice"
	"os"
	"strings"
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
	case "exec":
		if len(os.Args) < 3 {
			panic("Please enter object name")
		}

		if len(os.Args) < 4 {
			panic("Please enter the command you want to run")
		}

		if len(os.Args) == 4 {
			execCmdInPod(os.Args[2], os.Args[3])
			return
		}

		execCmdInPod(os.Args[2], os.Args[3], os.Args[4:]...)
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
	assertNameInConfig(yaml, name)
	oscmd.Run("kubectl", "apply", "-f", getBasePath(yaml)+name)
}

func deleteObject(yaml *Yaml, name string) {
	assertNameInConfig(yaml, name)
	oscmd.Run("kubectl", "delete", "-f", getBasePath(yaml)+name)
}

func assertNameInConfig(yaml *Yaml, name string) {
	objects, err := yaml.Get("objects").Array()
	if err != nil {
		panic(err)
	}

	if !slice.ContainsName(objects, name) {
		panic("Object name " + name + " not defined in config file")
	}
}

func execCmdInPod(appName string, cmd string, innerArgs ...string) {
	fullPodName := oscmd.RunForResult("kubectl", "get", "pod", "-l", "app="+appName, "-o", "name")

	baseKubeArgs := []string{"exec", strings.Trim(fullPodName, "pods/ \n"), cmd}
	oscmd.Run("kubectl", append(baseKubeArgs, innerArgs...)...)
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
