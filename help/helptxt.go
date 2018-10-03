package help

import (
	"fmt"
)

type HelpTextPrinter struct{}

func (printer HelpTextPrinter) Print() {
	text :=
		`
goctrl controls the life circle of your minikube cluster and the kubernetes objects running on it.

Available commands:

start  [--createobjs=false]   start minikube cluster and create all kubernetes objects from configuration
build                         build all the configuration files from template
create [objectname]           create the kubernetes object with given name
delete [objectname]           delete the kubernetes object with given name
exec   [objectname] [command] execute the command in a given kubernetes object
list                          list all the kubernetes objects running in minikube cluster
stop                          stop the minikube cluster
`

	fmt.Print(text)
}
