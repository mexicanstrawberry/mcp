package main

import (
	"github.com/mexicanstrawberry/mcp/gatekeeper"
	"github.com/mexicanstrawberry/mcp/recipe"
	"github.com/mexicanstrawberry/mcp/sensor"
	"github.com/mexicanstrawberry/mcp/web"
	clog "github.com/morriswinkler/cloudglog"
)

const (
	DEBUG = false
)

func init() {

}

func main() {

	clog.Infoln("init mexicanstrawberry")

	err := gatekeeper.MqttData.Dial()
	if err != nil {
		clog.Errorln("[mqtt] ", err)
	}

	myRecipe, err := recipe.NewRecipe("SimpleRecipe", true) // true for dummy recipe
	if err != nil {
		clog.Errorln(err)
	}

	insideHumidity, err := sensor.NewSensor("InsideHumidity", myRecipe)
	if err != nil {
		clog.Errorln(err)
	}

	insideTemperature, err := sensor.NewSensor("InsideTemperature", myRecipe)
	if err != nil {
		clog.Errorln(err)
	}

	go gatekeeper.Run()
	go insideHumidity.Run()
	go insideTemperature.Run()

	web.StartHttpServer()

}
