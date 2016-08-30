package gatekeeper

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"strings"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/eclipse/paho.mqtt.golang"
	clog "github.com/morriswinkler/cloudglog"
)

const (
	DEBUG = false
)

type mqttDataT struct {
	Options mqtt.ClientOptions
	Client  mqtt.Client
	Channel chan interface{}
}

var CurrentData map[string]interface{}

func (p *mqttDataT) Dial() error {

	if token := p.Client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	p.Client.Subscribe("iot-2/type/RPi/id/Plant1/evt/+/fmt/json", 0, MessageHandler)

	return nil
}

var MqttData mqttDataT

func (m *mqttDataT) SetChannel(c chan interface{}) {
	m.Channel = c
}

var MessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

	var a interface{}

	json.Unmarshal(msg.Payload(), &a)

	for _, v := range a.(map[string]interface{}) {
		// TODO: UPdate map instead of overwrite
		CurrentData = v.(map[string]interface{})
	}
}

func init() {

	if DEBUG {
		mqtt.ERROR = log.New(os.Stdout, "ERROR: ", 0)
		mqtt.CRITICAL = log.New(os.Stdout, "CRITCIAL: ", 0)
		mqtt.WARN = log.New(os.Stdout, "WARN: ", 0)
		mqtt.DEBUG = log.New(os.Stdout, "DEBUG: ", 0)
	}

	appEnv, _ := cfenv.Current()

	//clog.Info(appEnv.Services)
	service, _ := appEnv.Services.WithName("MS-IoT")

	clog.Info(service.Credentials)

	//	MqttData.Options.SetClientID("a:7mqeaj:dxczjv6qsw")
	//	MqttData.Options.SetUsername("a-7mqeaj-dxczjv6qsw")
	//	MqttData.Options.SetPassword("ysXxLYoF)cNB*zy5N+")
	//	MqttData.Options.AddBroker("tls://7mqeaj.messaging.internetofthings.ibmcloud.com:8883")

	id := strings.Replace(service.Credentials["apiKey"].(string), "-", ":", -1)

	MqttData.Options.SetClientID(id)
	MqttData.Options.SetUsername(service.Credentials["apiKey"].(string))
	MqttData.Options.SetPassword(service.Credentials["apiToken"].(string))
	MqttData.Options.AddBroker("tls://" + service.Credentials["mqtt_host"].(string) + ":" + service.Credentials["mqtt_s_port"].(string))
	MqttData.Options.SetDefaultPublishHandler(MessageHandler)
	MqttData.Client = mqtt.NewClient(&MqttData.Options)
}

func Run() {

	for {
		select {
		case <-time.After(10 * time.Second):
			clog.Info("Gatekeeper JOB")
		}
	}

}

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
