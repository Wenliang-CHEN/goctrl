# goctrl
goctrl controls the life circle of your minikube cluster and the kubernetes objects running on it.

Currently goctrl supports MAC machines only.

## Prerequisite 
In order to install and run goctrl, you need the following installed

Minikube: https://kubernetes.io/docs/setup/minikube/

govendor (for installation): https://github.com/kardianos/govendor

## Configuration
You can find a sample of configuration here: https://github.com/Wenliang-CHEN/kubernetes-object-sample

You need to point the BASE_PATH to your configuration folder in this file: https://github.com/Wenliang-CHEN/goctrl/blob/master/config/path.go

## Installation
To install, just pull this repo to your GOPATH and do 

`go install`

## Available commands
`start  [--createobjs=false]`   start minikube cluster and create all kubernetes objects from configuration

`build                      `   build all the configuration files from template

`create [objectname]        `   create the kubernetes object with given name

`delete [objectname]        `   delete the kubernetes object with given name

`exec   [objectname] [command]` execute the command in a given kubernetes object

`list`                          list all the kubernetes objects running in minikube cluster

`stop`                          stop the minikube cluster
