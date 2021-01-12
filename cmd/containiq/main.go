package main

import (
	"github.com/mclenhard/containiq/pkg/setup"
	"github.com/mclenhard/containiq/pkg/controller"
)



func main() {
	config := setup.Setup()
	controller.Controller(config)

}