package gatekeeper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"strings"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/mexicanstrawberry/go-cfenv"
	clog "github.com/morriswinkler/cloudglog"
)

// myttDataT holds the mqtt connection and event channel for communication
type mqttDataT struct {
	Options mqtt.ClientOptions
	Client  mqtt.Client
	Channel chan interface{}
}

var MqttData mqttDataT // holds mqtt connection, subscription data

// dial connections to the mqtt broaker set's up subscriptions
// do not call until mqtt.ClientOptions is set
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

	Debug(true) // enable debugging during mqtt connect

	//TODO: unify this in a common function for all services
	//	- web/couchdbProxy.go
	// this function is fired if VCAP_SERVICES was not found and refires every 10 Seconds
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

		service, err := appEnv.Services.WithLabel("iotf-service")
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

	Debug(false) // disable debuggig
}

// Debug takes a bool and en / disables mqtt debugging
func Debug(on bool) {

	if on {
		clog.Infoln("[mqtt] enable debugging")
		mqtt.ERROR = log.New(os.Stdout, "ERROR: ", 0)
		mqtt.CRITICAL = log.New(os.Stdout, "CRITCIAL: ", 0)
		mqtt.WARN = log.New(os.Stdout, "WARN: ", 0)
		mqtt.DEBUG = log.New(os.Stdout, "DEBUG: ", 0)
	} else {
		clog.Infoln("[mqtt] disabling debugging")
		mqtt.ERROR = log.New(ioutil.Discard, "", 0)
		mqtt.CRITICAL = log.New(ioutil.Discard, "", 0)
		mqtt.WARN = log.New(ioutil.Discard, "", 0)
		mqtt.DEBUG = log.New(ioutil.Discard, "", 0)
	}
}
