package gatekeeper

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"strings"

	"fmt"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/mexicanstrawberry/go-cfenv"
	"github.com/mexicanstrawberry/mcp/events"
	clog "github.com/morriswinkler/cloudglog"
)

var (
	DEBUG = false // used to enable debuging inside this package

	CurrentData     map[string]interface{}
	CurrentCommands = make([]events.MqttCommand, 0)
)

// myttDataT holds the mqtt connection and event channel for communication
type mqttDataT struct {
	Options mqtt.ClientOptions
	Client  mqtt.Client
	Channel chan interface{}
}

var MqttData mqttDataT

// dial mqtt connections
// this function is privat, do not call until mqtt.ClientOptions been set
func (p *mqttDataT) dial() error {

	if token := p.Client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	p.Client.Subscribe("iot-2/type/RPi/id/Plant1/evt/+/fmt/json", 0, MessageHandler)

	return nil
}

// TODO: what is this for
func (m *mqttDataT) SetChannel(c chan interface{}) {
	m.Channel = c
}

// Default mqtt Message handler
var MessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

	var a interface{}

	json.Unmarshal(msg.Payload(), &a)

	for _, v := range a.(map[string]interface{}) {
		// TODO: Update map instead of overwrite
		CurrentData = v.(map[string]interface{})
	}
}

// init reads VCAP_SERCVICES and initialises mqtt connection
func init() {

	if DEBUG {
		mqtt.ERROR = log.New(os.Stdout, "ERROR: ", 0)
		mqtt.CRITICAL = log.New(os.Stdout, "CRITCIAL: ", 0)
		mqtt.WARN = log.New(os.Stdout, "WARN: ", 0)
		mqtt.DEBUG = log.New(os.Stdout, "DEBUG: ", 0)
	}

	//TODO: unify this in a common function for all services
	//	- web/couchdbProxy.go
	// this function is fird if VCAP_SERVICES was not found and refires every 10 Seconds
	var vcapServicesNotFound func() // declare first to access var recursively
	vcapServicesNotFound = func() {

		clog.Errorln("[gatekeeper] ###############################################")
		clog.Errorln("[gatekeeper] CRITICAL ERROR: could not read VCAP_SERVICES   ")
		clog.Errorln("[gatekeeper] ###############################################")
		for _, k := range os.Environ() {
			clog.Errorln(k)
		}

		time.AfterFunc(time.Second*10, vcapServicesNotFound) // recursion

	}

	appEnv, err := cfenv.Current()
	if err != nil {
		// this is critical VCAP_SERVICES was not found
		vcapServicesNotFound()
	} else {

		service, err := appEnv.Services.WithName("MS-IoT")
		if err != nil {
			clog.Errorln("[gatekeeper] ", err)

		} else {

			id := strings.Replace(service.Credentials["apiKey"].(string), "-", ":", -1)
			broker := fmt.Sprintf("tls://%s:%.f", service.Credentials["mqtt_host"].(string), service.Credentials["mqtt_s_port"].(float64))
			MqttData.Options.SetClientID(id)
			MqttData.Options.SetUsername(service.Credentials["apiKey"].(string))
			MqttData.Options.SetPassword(service.Credentials["apiToken"].(string))
			MqttData.Options.AddBroker(broker)
			MqttData.Options.SetDefaultPublishHandler(MessageHandler)
			MqttData.Client = mqtt.NewClient(&MqttData.Options)

			err := MqttData.dial()
			if err != nil {
				clog.Errorln("[mqtt] ", err)
			}
		}
	}
}

func Run() {

	myTicker := time.NewTicker(time.Duration(10) * time.Second)

	for {
		select {

		case <-myTicker.C:
			// TODO: is this a dummy function ?
			for _, command := range CurrentCommands {
				channel := fmt.Sprintf("iot-2/type/RPi/id/Plant1/cmd/%s/fmt/json", command.Actuator)
				payload := fmt.Sprintf("{\"d\":{ \"value\": %f}}", command.Value)
				MqttData.Client.Publish(channel, 0, false, payload)
			}
			CurrentCommands = make([]events.MqttCommand, 0)

		case i := <-events.Channel:
			switch i.(type) {
			case events.MqttCommand:
				newCommand := i.(events.MqttCommand)
				match := false

				for i := range CurrentCommands {
					oldCommand := CurrentCommands[i]
					if newCommand.Actuator == oldCommand.Actuator {
						match = true
						if newCommand.Priority < oldCommand.Priority {
							CurrentCommands[i] = newCommand
						}
					}
				}

				if !match {
					CurrentCommands = append(CurrentCommands, newCommand)
				}

			}
		}
	}

}
