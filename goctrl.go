package main

import (
	"fmt"
	. "github.com/smallfish/simpleyaml"
	config "goCtrl/config"
	errors "goCtrl/errhandle"
	oscmd "goCtrl/oscmd"
	parser "goCtrl/parser"
	slice "goCtrl/utils/slice"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	defer errors.HandleError()

	parameters := parser.Parse(config.BASE_PATH + "parameters.yaml")
	ctrlCmd := os.Args[1]

	switch ctrlCmd {
	case "start":
		startMinikube()
		createObjects(parameters)
	case "build":
		if len(os.Args) < 3 {
			panic("Please enter a kubernetes config file to build.")
		}
		buildYaml(parameters, os.Args[2])
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
	case "list":
		list(parameters)
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
	objects := getDefinedKubeObjects(yaml)

	objectNames := slice.Map(objects, func(object interface{}) interface{} {
		return slice.ToMap(object)["name"]
	})

	for _, objectName := range objectNames {
		buildYaml(yaml, objectName.(string))
		createObject(yaml, objectName.(string))
	}
}

func createObject(yaml *Yaml, name string) {
	assertNameInConfig(yaml, name)
	deleteObject(yaml, name)
	oscmd.Run("kubectl", "apply", "-f", getBuiltPath(yaml)+name)
}

func deleteObject(yaml *Yaml, name string) {
	assertNameInConfig(yaml, name)

	result := oscmd.RunForResult("kubectl", "delete", "-f", getBuiltPath(yaml)+name)
	if strings.Contains(result, "NotFound") {
		fmt.Println("Object \"" + name + "\" does not exist.  No deletion executed.")
		return
	}

	fmt.Printf("%v", result)
}

var parameterContainsName = func(object interface{}, name interface{}) bool {
	return slice.ToMap(object)["name"].(string) == name
}

func assertNameInConfig(yaml *Yaml, name string) {
	objects := getDefinedKubeObjects(yaml)

	isObjNameDefined := slice.Contains(objects, name, parameterContainsName)

	if !isObjNameDefined {
		panic("Object name " + name + " not defined in config file")
	}
}

func execCmdInPod(appName string, cmd string, innerArgs ...string) {
	fullPodName := oscmd.RunForResult("kubectl", "get", "pod", "-l", "app="+appName, "-o", "name")

	baseKubeArgs := []string{"exec", strings.Trim(fullPodName, "pods/ \n"), cmd}
	oscmd.Run("kubectl", append(baseKubeArgs, innerArgs...)...)
}

func stopMinikube() {
	oscmd.Run("minikube", "stop")
}

func buildYaml(yaml *Yaml, name string) {
	objects := getDefinedKubeObjects(yaml)

	parameters, ok := slice.First(objects, name, parameterContainsName)
	if ok == false {
		panic("Object name" + name + "is not defined in config file.")
	}

	buff, err := ioutil.ReadFile(getTemplatePath(yaml) + name + "/" + name + ".yaml")
	if err != nil {
		panic("Unable to read template file.  Please make sure it exists.")
	}

	content := string(buff)
	content = strings.Replace(content, "{TIMESTAMP}", time.Now().Format(time.RFC3339), -1)
	configs, exist := slice.ToMap(parameters)["config"]

	if !exist {
		writeToBuild(yaml, name, content)
		return
	}

	for key, val := range slice.ToMap(configs) {
		content = strings.Replace(content, "{"+key.(string)+"}", val.(string), -1)
	}
	writeToBuild(yaml, name, content)
}

func getDefinedKubeObjects(yaml *Yaml) []interface{} {
	objects, err := yaml.Get("objects").Array()
	if err != nil {
		panic(err)
	}

	return objects
}

func writeToBuild(yaml *Yaml, name string, content string) {
	buildPath := getBuiltPath(yaml) + name
	oscmd.Run("rm", "-rf", buildPath)
	oscmd.Run("mkdir", buildPath)
	ioutil.WriteFile(buildPath+"/"+name+".yaml", []byte(content), 0644)
}

func list(yaml *Yaml) {
	var objectTypes = make(map[string]bool)
	err := filepath.Walk(getBuiltPath(yaml),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !strings.Contains(path, ".yaml") {
				return nil
			}

			object := parser.Parse(path)
			objectType, parseErr := object.Get("kind").String()
			if parseErr != nil {
				return parseErr
			}

			objectTypes[objectType] = true
			return nil
		})

	if err != nil {
		fmt.Println(err)
	}

	for objectType, _ := range objectTypes {
		objectKeyword := strings.ToLower(objectType + "s")

		fmt.Printf("%v: \n", objectKeyword)
		oscmd.Run("kubectl", "get", objectKeyword)
		fmt.Println()
	}
}

func getBuiltPath(yaml *Yaml) string {
	return config.BASE_PATH + "build/"
}

func getTemplatePath(yaml *Yaml) string {
	return config.BASE_PATH + "template/"
}
