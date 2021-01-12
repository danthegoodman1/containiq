package main

import (
	"github.com/containiq/containiq/pkg/setup"
	"github.com/containiq/containiq/pkg/controller"
)



func main() {
	config := setup.Setup()
	controller.Controller(config)

}