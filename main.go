package main

import (
	"log"

	"os"

	"encoding/json"

	"time"

	"github.com/eclipse/paho.mqtt.golang"
	clog "github.com/morriswinkler/cloudglog"
	"hub.jazz.net/git/ansi/MS-FE/recipe"
	"hub.jazz.net/git/ansi/MS-FE/sensor"

	"hub.jazz.net/git/ansi/MS-FE/gatekeeper"
	"hub.jazz.net/git/ansi/MS-FE/events"
)

const (
	DEBUG = false
)

func init() {

	if DEBUG {
		mqtt.ERROR = log.New(os.Stdout, "ERROR: ", 0)
		mqtt.CRITICAL = log.New(os.Stdout, "CRITCIAL: ", 0)
		mqtt.WARN = log.New(os.Stdout, "WARN: ", 0)
		mqtt.DEBUG = log.New(os.Stdout, "DEBUG: ", 0)
	}
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	clog.Infoln("TOPIC: %s\n", msg.Topic())
	clog.Infoln("MSG: %s\n", msg.Payload())
}

var sensorMessage mqtt.MessageHandler = func(clien mqtt.Client, msg mqtt.Message) {
	//clog.Infoln("&#v", msg)

	var reading interface{}

	json.Unmarshal(msg.Payload(), &reading)

	for i, v := range reading.(map[string]interface{}) {

		clog.Infoln(i)
		clog.Infoln(v)

	}
}

func main() {

	clog.Infoln("init mexicanstrawberry")

	gatekeeper.MqttData.Dial()






	//i := 0
	//for _ = range time.Tick(time.Duration(1) * time.Second) {
	//	if i == 100 {
	//		break
	//	}
	//
	//	text := fmt.Sprintf("\"d\":{ msg %d}", i)
	//	mqttClient.Publish("iot-2/type/RPi/id/dummyposter/evt/Plant2/fmt/json", 0, false, text)
	//	i++
	//}

	//clog.LogLevel = 0

	r, err := recipe.NewRecipe("SimpleRecipe", true)
	if err != nil {
		clog.Errorln(err)
	}

	s, err := sensor.NewSensor("InsideHumidity", r)
	if err != nil {
		clog.Errorln(err)
	}

	ctl := make(chan sensor.CtrlChannel)

	go s.Run(ctl)


	for{
		l := <-events.Channel
		clog.Info("===================================")
		switch e := l.(type){
		case events.MqttRecive:
			for k,v := range e {
				clog.Info(k)
				clog.Info(v)
			}
		}
	}

	time.Sleep(time.Duration(60) * time.Second)
	//ctl <- 1

}
