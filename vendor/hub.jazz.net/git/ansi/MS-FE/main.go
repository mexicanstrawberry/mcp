package main

import (
	"time"

	clog "github.com/morriswinkler/cloudglog"
	"hub.jazz.net/git/ansi/MS-FE/recipe"
	"hub.jazz.net/git/ansi/MS-FE/sensor"

	"hub.jazz.net/git/ansi/MS-FE/gatekeeper"
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

	time.Sleep(time.Duration(60) * time.Second)

}
