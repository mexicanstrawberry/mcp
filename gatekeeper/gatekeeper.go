package gatekeeper

import (
	"time"

	"fmt"

	"github.com/mexicanstrawberry/mcp/events"
)

var (
	CurrentData     map[string]interface{}          // current data from mqtt broker
	CurrentCommands = make([]events.MqttCommand, 0) // TODO: documet variable
)

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
