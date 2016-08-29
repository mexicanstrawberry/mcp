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

	opts := mqtt.NewClientOptions()

	opts.SetClientID("a:7mqeaj:8xgxdkgi7y")
	opts.SetUsername("a-7mqeaj-8xgxdkgi7y")
	opts.SetPassword("saiwUQG6n2@uwFbC!o")
	opts.AddBroker("tls://7mqeaj.messaging.internetofthings.ibmcloud.com:8883")

	//opts.SetClientID("ssl-sample").SetTLSConfig(tlsconfig)
	opts.SetDefaultPublishHandler(f)

	mqttClient := mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//mqttClient.Subscribe("iot-2/type/RPi/id/dummyposter/evt/Plant1/fmt/json", 0, f)

	mqttClient.Subscribe("iot-2/type/RPi/id/Plant1/evt/+/fmt/json", 0, sensorMessage)

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

	time.Sleep(time.Duration(60) * time.Second)
	//ctl <- 1

	mqttClient.Disconnect(250)

}
