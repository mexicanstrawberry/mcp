package main

import (
	"github.com/mexicanstrawberry/mcp/recipe"
	"github.com/mexicanstrawberry/mcp/sensor"
	clog "github.com/morriswinkler/cloudglog"

	"github.com/mexicanstrawberry/mcp/gatekeeper"
)

const (
	DEBUG = false
)

func init() {

}

func main() {

	clog.Infoln("init mexicanstrawberry")

	gatekeeper.MqttData.Dial()

	r, err := recipe.NewRecipe("SimpleRecipe", true)
	if err != nil {
		clog.Errorln(err)
	}

	s, err := sensor.NewSensor("InsideHumidity", r)
	if err != nil {
		clog.Errorln(err)
	}

	go s.Run()
	go gatekeeper.Run()

	StartHttpServer()

}
