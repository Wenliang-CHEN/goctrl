package main

import (
	"flag"
	config "goCtrl/config"
	errors "goCtrl/errhandle"
	provider "goCtrl/provider"
	"os"
)

func main() {
	serviceContainer := provider.GetDependencyContainer()
	defer serviceContainer.ErrorHandler.HandleError()

	if len(os.Args) < 2 {
		panic(errors.InvalidCommand)
	}

	parameters := serviceContainer.Parser.Parse(config.BASE_PATH + "parameters.yaml")
	ctrlCmd := os.Args[1]

	switch ctrlCmd {
	case "start":
		createObjs := flag.Bool("createobjs", false, "create all defined kubernetes objects")
		serviceContainer.OsCtrl.StartMinikube(createObjs, parameters)
	case "build":
		if len(os.Args) < 3 {
			panic("Please enter a kubernetes config file to build.")
		}
		serviceContainer.OsCtrl.BuildYaml(parameters, os.Args[2])
	case "create":
		if len(os.Args) < 3 {
			panic("Please enter object name")
		}

		objectName := os.Args[2]
		serviceContainer.OsCtrl.CreateObject(parameters, objectName)
	case "delete":
		if len(os.Args) < 3 {
			panic("Please enter object name")
		}

		objectName := os.Args[2]
		serviceContainer.OsCtrl.DeleteObject(parameters, objectName)
	case "exec":
		if len(os.Args) < 3 {
			panic("Please enter object name")
		}

		if len(os.Args) < 4 {
			panic("Please enter the command you want to run")
		}

		if len(os.Args) == 4 {
			serviceContainer.OsCtrl.ExecCmdInPod(os.Args[2], os.Args[3])
			return
		}

		serviceContainer.OsCtrl.ExecCmdInPod(os.Args[2], os.Args[3], os.Args[4:]...)
	case "list":
		serviceContainer.OsCtrl.List(parameters)
	case "stop":
		serviceContainer.OsCtrl.StopMinikube()
	default:
		panic(errors.InvalidCommand)
	}
}
