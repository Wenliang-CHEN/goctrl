package provider

import (
	inject "github.com/codegangsta/inject"
	errhandle "goCtrl/errhandle"
	osctrl "goCtrl/osctrl"
	parsers "goCtrl/parsers"
)

type DependencyContainer struct {
	ErrorHandler *errhandle.ErrorHandler `inject`
	Parser       *parsers.YamlParser     `inject`
	OsCtrl       *osctrl.MacCtrl         `inject`
}

func GetDependencyContainer() DependencyContainer {
	container := DependencyContainer{}
	injector := inject.New()
	injector.Map(&errhandle.ErrorHandler{})
	injector.Map(&parsers.YamlParser{})
	injector.Map(&osctrl.MacCtrl{})
	err := injector.Apply(&container)
	if err != nil {
		panic(err)
	}
	return container
}
