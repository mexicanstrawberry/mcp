package gatekeeper

import (
	"github.com/eclipse/paho.mqtt.golang"
	//clog "github.com/morriswinkler/cloudglog"
	"encoding/json"

	"hub.jazz.net/git/ansi/MS-FE/events"
)

type mqttDataT struct {
	Options mqtt.ClientOptions
	Client  mqtt.Client
	Channel chan interface{}
}

func (p *mqttDataT) Dial() error {

	if token := p.Client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	p.Client.Subscribe("iot-2/type/RPi/id/Plant1/evt/+/fmt/json", 0, p.MessageHandler())

	return nil
}

var MqttData mqttDataT

func (m *mqttDataT) SetChannel(c chan interface{}) {
	m.Channel = c
}

func (p *mqttDataT) MessageHandler() func(client mqtt.Client, msg mqtt.Message) {

	var m mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

		var a interface{}

		json.Unmarshal(msg.Payload(), &a)

		for _, v := range a.(map[string]interface{}) {
			//clog.Info("<==============================")

			p.Channel <- events.MqttRecive{
				Map: v.(map[string]interface{}),
			}
		}

	}

	return m
}

func init() {
	MqttData.Options.SetClientID("a:7mqeaj:8xgxdkgi7y")
	MqttData.Options.SetUsername("a-7mqeaj-8xgxdkgi7y")
	MqttData.Options.SetPassword("saiwUQG6n2@uwFbC!o")
	MqttData.Options.AddBroker("tls://7mqeaj.messaging.internetofthings.ibmcloud.com:8883")
	MqttData.Options.SetDefaultPublishHandler(MqttData.MessageHandler())
	MqttData.Client = mqtt.NewClient(&MqttData.Options)
}
