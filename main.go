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

	web.StartHttpServer()

}
