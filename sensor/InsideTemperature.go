package sensor

import (
	"time"

	"github.com/mexicanstrawberry/mcp/events"
	"github.com/mexicanstrawberry/mcp/gatekeeper"
	"github.com/mexicanstrawberry/mcp/recipe"

	clog "github.com/morriswinkler/cloudglog"
)

type InsideTemperature struct {
	SensorName  string
	Recipe      *recipe.Recipe
	StartTime   time.Time
	Time        time.Time
	TargetValue float64
	Timer       Timer
	Ticker      *time.Ticker
}

func (ih *InsideTemperature) Init() {
	ih.SensorName = "InsideTemperature"
}

func (ih *InsideTemperature) Name() string {
	return ih.SensorName
}

func (ih *InsideTemperature) SetRecipe(r *recipe.Recipe) {
	ih.Recipe = r
}

func (ih *InsideTemperature) regulate() {
	if insideTemperature, exist := gatekeeper.CurrentData[ih.SensorName]; exist {
		if insideTemperature == nil {
			return
		}

		toColdTowardsTarget := ih.TargetValue - insideTemperature.(float64)

		outsideComparedToInside := 0.0

		if outsideTemperature, exist := gatekeeper.CurrentData["OutsideTemperature"]; exist {
			if outsideTemperature != nil {
				outsideComparedToInside = outsideTemperature.(float64) - insideTemperature.(float64)
			}
		}

		// To Cold and outside warmer (can't imagine how this can happen)
		if toColdTowardsTarget > 2 && outsideComparedToInside > 5 {
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Hatch",
				Value:     1, // open
				Priority:  events.Priority(events.P_NORM),
			}
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "OutsideFan",
				Value:     11,
				Priority:  events.Priority(events.P_NORM),
			}
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "InsideFan",
				Value:     11,
				Priority:  events.Priority(events.P_NORM),
			}
		}

		// To cold and outside not warmer
		if toColdTowardsTarget > 2 && outsideComparedToInside <= 5 {
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Hatch",
				Value:     0, // close
				Priority:  events.Priority(events.P_NORM),
			}
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Heater",
				Value:     11,
				Priority:  events.Priority(events.P_NORM),
			}
		}

		// To Warm
		if toColdTowardsTarget < -1 && outsideComparedToInside <= 5 {
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Hatch",
				Value:     1, // open
				Priority:  events.Priority(events.P_MIN),
			}
		}
		if toColdTowardsTarget < -5 && outsideComparedToInside <= 5 {
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Hatch",
				Value:     1, // open
				Priority:  events.Priority(events.P_NORM),
			}
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "OutsideFan",
				Value:     11,
				Priority:  events.Priority(events.P_NORM),
			}
		}
	}
}

func (ih *InsideTemperature) Run() {

	ih.StartTime = time.Now() // use StartTime to calc elapsed time
	ih.Ticker = time.NewTicker(time.Duration(defaultTickInterval) * time.Second)
	ih.Timer = ih.nextOpTime()

	for {
		select {
		case t := <-ih.Ticker.C:
			clog.V(3).Infoln("[Ticker] ", t)
			ih.regulate()

		case t := <-ih.Timer.Chan.C:
			clog.V(3).Infoln("[Timer] ", t)
			ih.Timer = ih.nextOpTime()

		}
	}
}

func (ih *InsideTemperature) nextOpTime() Timer {

	var t Timer

	// get elapsed time
	elapsedTime := time.Since(ih.StartTime)

	// convert elapsed time into time.Time
	nextTime := new(time.Time).Add(elapsedTime)

	// find next operaton and arm time.Timer
	nextInterval, nextValue, err := ih.Recipe.Operation.NextOperation(ih.Name(), nextTime)
	if err != nil {

		switch err {
		case recipe.E_NO_NEXT_OPERATION:
			clog.Infoln(err)
		case recipe.E_SENSOR_NOT_FOUND:
			clog.Infoln(err)
		default:
			clog.Errorln(err)
		}
		t.Func = new(time.Timer)
		t.Chan = new(time.Timer)

		return t
	}

	clog.V(3).Infoln("[nextTimer] ", nextInterval)

	// set func timer
	t.Func = time.AfterFunc(nextInterval, func() {
		clog.V(3).Infoln("[nextTimer] fire")
		ih.TargetValue = nextValue
	})
	// set receiver channel
	t.Chan = time.NewTimer(nextInterval)

	return t
}
