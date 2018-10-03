package provider

import (
	inject "github.com/codegangsta/inject"
	errhandle "goCtrl/errhandle"
)

type DependencyContainer struct {
	ErrorHandler *errhandle.ErrorHandler `inject`
}

func GetDependencyContainer() DependencyContainer {
	container := DependencyContainer{}
	injector := inject.New()
	injector.Map(&errhandle.ErrorHandler{})
	err := injector.Apply(&container)
	if err != nil {
		panic(err)
	}
	return container
}
