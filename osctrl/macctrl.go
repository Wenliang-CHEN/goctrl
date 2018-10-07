package osctrl

import (
	"fmt"
	. "github.com/smallfish/simpleyaml"
	config "goCtrl/config"
	oscmd "goCtrl/oscmd"
	parsers "goCtrl/parsers"
	slice "goCtrl/utils/slice"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MacCtrl struct {
	parser    parsers.YamlParser
	cmdRunner oscmd.Runner
}

func (ctrl MacCtrl) StartMinikube(shouldCreateObjs *bool, parameters *Yaml) {
	if ctrl.isServerRunning() {
		panic("Server is already running.")
	}

	ctrl.cmdRunner.Run("minikube", "start")

	if *shouldCreateObjs {
		ctrl.createObjects(parameters)
	}
}

func (ctrl MacCtrl) isServerRunning() bool {
	status := ctrl.cmdRunner.RunForResult("minikube", "status")
	return strings.Contains(status, "minikube: Running")
}

func (ctrl MacCtrl) createObjects(yaml *Yaml) {
	objects := ctrl.getDefinedKubeObjects(yaml)

	objectNames := slice.Map(objects, func(object interface{}) interface{} {
		return slice.ToMap(object)["name"]
	})

	for _, objectName := range objectNames {
		ctrl.BuildYaml(yaml, objectName.(string))
		ctrl.CreateObject(yaml, objectName.(string))
	}
}

func (ctrl MacCtrl) CreateObject(yaml *Yaml, name string) {
	ctrl.assertNameInConfig(yaml, name)
	ctrl.DeleteObject(yaml, name)

	deleted := false
	objectConfigPath := ctrl.getBuiltPath(yaml) + name
	for !deleted {
		deleted = strings.ContainsAny(ctrl.cmdRunner.RunForResult("kubectl", "get", "-f", objectConfigPath), "NotFound Terminating")
		time.Sleep(2 * time.Second)
	}
	ctrl.cmdRunner.Run("kubectl", "apply", "-f", ctrl.getBuiltPath(yaml)+name)
}

func (ctrl MacCtrl) DeleteObject(yaml *Yaml, name string) {
	ctrl.assertNameInConfig(yaml, name)

	result := ctrl.cmdRunner.RunForResult("kubectl", "delete", "-f", ctrl.getBuiltPath(yaml)+name)
	if strings.Contains(result, "NotFound") {
		fmt.Println("Object \"" + name + "\" does not exist.  No deletion executed.")
		return
	}

	fmt.Printf("%v", result)
}

func (ctrl MacCtrl) assertNameInConfig(yaml *Yaml, name string) {
	objects := ctrl.getDefinedKubeObjects(yaml)

	isObjNameDefined := slice.Contains(objects, name, parameterContainsName)
	if !isObjNameDefined {
		panic("Object name " + name + " not defined in config file")
	}
}

func (ctrl MacCtrl) ExecCmdInPod(appName string, cmd string, innerArgs ...string) {
	fullPodName := ctrl.cmdRunner.RunForResult("kubectl", "get", "pod", "-l", "app="+appName, "-o", "name")

	baseKubeArgs := []string{"exec", strings.Trim(fullPodName, "pods/ \n"), cmd}
	ctrl.cmdRunner.Run("kubectl", append(baseKubeArgs, innerArgs...)...)
}

func (ctrl MacCtrl) StopMinikube() {
	ctrl.cmdRunner.Run("minikube", "stop")
}

func (ctrl MacCtrl) BuildYaml(yaml *Yaml, name string) {
	objects := ctrl.getDefinedKubeObjects(yaml)

	parameters, ok := slice.First(objects, name, parameterContainsName)
	if ok == false {
		panic("Object name" + name + "is not defined in config file.")
	}

	buff, err := ioutil.ReadFile(ctrl.getTemplatePath(yaml) + name + "/" + name + ".yaml")
	if err != nil {
		panic("Unable to read template file.  Please make sure it exists.")
	}

	content := string(buff)
	content = strings.Replace(content, "{TIMESTAMP}", time.Now().Format(time.RFC3339), -1)
	configs, exist := slice.ToMap(parameters)["config"]

	if !exist {
		ctrl.writeToBuild(yaml, name, content)
		return
	}

	for key, val := range slice.ToMap(configs) {
		content = strings.Replace(content, "{"+key.(string)+"}", val.(string), -1)
	}
	ctrl.writeToBuild(yaml, name, content)
}

var parameterContainsName = func(object interface{}, name interface{}) bool {
	return slice.ToMap(object)["name"].(string) == name
}

func (ctrl MacCtrl) getDefinedKubeObjects(yaml *Yaml) []interface{} {
	objects, err := yaml.Get("objects").Array()
	if err != nil {
		panic(err)
	}

	return objects
}

func (ctrl MacCtrl) writeToBuild(yaml *Yaml, name string, content string) {
	if !ctrl.isServerRunning() {
		panic("Kubernetes server is not running.  Please start first.")
	}
	buildPath := ctrl.getBuiltPath(yaml) + name
	ctrl.cmdRunner.Run("rm", "-rf", buildPath)
	ctrl.cmdRunner.Run("mkdir", buildPath)
	ioutil.WriteFile(buildPath+"/"+name+".yaml", []byte(content), 0644)
}

func (ctrl MacCtrl) List(yaml *Yaml) {
	var objectTypes = make(map[string]bool)
	err := filepath.Walk(ctrl.getBuiltPath(yaml),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !strings.Contains(path, ".yaml") {
				return nil
			}

			object := ctrl.parser.Parse(path)
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
		ctrl.cmdRunner.Run("kubectl", "get", objectKeyword)
		fmt.Println()
	}
}

func (ctrl MacCtrl) getBuiltPath(yaml *Yaml) string {
	return config.BASE_PATH + "build/"
}

func (ctrl MacCtrl) getTemplatePath(yaml *Yaml) string {
	return config.BASE_PATH + "template/"
}
