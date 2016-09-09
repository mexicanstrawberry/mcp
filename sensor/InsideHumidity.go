package sensor

import (
	"time"

	"github.com/mexicanstrawberry/mcp/events"
	"github.com/mexicanstrawberry/mcp/gatekeeper"
	"github.com/mexicanstrawberry/mcp/recipe"

	"math"

	clog "github.com/morriswinkler/cloudglog"
)

type InsideHumidity struct {
	SensorName  string
	Recipe      *recipe.Recipe
	StartTime   time.Time
	Time        time.Time
	TargetValue float64
	Timer       Timer
	Ticker      *time.Ticker
}

func (ih *InsideHumidity) Init() {
	ih.SensorName = "InsideHumidity"
}

func (ih *InsideHumidity) Name() string {
	return ih.SensorName
}

func (ih *InsideHumidity) SetRecipe(r *recipe.Recipe) {
	ih.Recipe = r
}

func (ih *InsideHumidity) regulate() {

	if insideHumidity, exist := gatekeeper.CurrentData[ih.SensorName]; exist {

		if insideHumidity == nil {
			return
		}

		toDryTowardsTarget := ih.TargetValue - insideHumidity.(float64)
		outsideCompatedToInside := 0.0

		if outsideHumidity, exist := gatekeeper.CurrentData["OutsideHumidity"]; exist {
			if outsideHumidity != nil {
				outsideCompatedToInside = outsideHumidity.(float64) - insideHumidity.(float64)
			}
		}

		// To dry and outside as more humidity: open hatch and activate fan
		if toDryTowardsTarget > 1 && outsideCompatedToInside > 0 {
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Hatch",
				Value:     1, // open
				Priority:  events.Priority(events.P_MIN),
			}
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "OutsideFan",
				Value:     11,
				Priority:  events.Priority(events.P_NORM),
			}
		}
		// To dry and outside as less humidity close hatch and activate humidifier
		if toDryTowardsTarget > 1 && outsideCompatedToInside <= 0 {
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Humidifier",
				Value:     5,
				Priority:  events.Priority(events.P_NORM),
			}
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "InsideFan",
				Value:     8,
				Priority:  events.Priority(events.P_NORM),
			}
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Hatch",
				Value:     0, // close
				Priority:  events.Priority(events.P_NORM),
			}
		}
		// To much humidity and outside drier open the hatch
		if toDryTowardsTarget < -1 && outsideCompatedToInside < 0 {
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Hatch",
				Value:     1, // open
				Priority:  events.Priority(events.P_MIN),
			}
		}
		// Much to much humidity and outside drier open the hatch and use the fan
		if toDryTowardsTarget < -5 && outsideCompatedToInside < 0 {
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Hatch",
				Value:     1, // open
				Priority:  events.Priority(events.P_MIN),
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
		if math.Abs(toDryTowardsTarget) < 1 {
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Hatch",
				Value:     0, // close
				Priority:  events.Priority(events.P_MIN),
			}
		}
	}

}

func (ih *InsideHumidity) Run() {

	ih.StartTime = time.Now() // use StartTime to calc elapsed time
	ih.Ticker = time.NewTicker(time.Duration(defaultTickInterval) * time.Second)
	ih.Timer = ih.nextOpTime()

	for {
		select {
		case t := <-ih.Ticker.C:
			//clog.Infoln("[Ticker] ", t)
			ih.regulate()

		case t := <-ih.Timer.Chan.C:
			//clog.Infoln("[Timer] ", t)
			ih.Timer = ih.nextOpTime()

		}
	}
}

func (ih *InsideHumidity) nextOpTime() Timer {

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
		clog.V(3).Infof("[nextTimer] fire value: %f \n", nextValue)
		ih.TargetValue = nextValue
	})
	// set receiver channel
	t.Chan = time.NewTimer(nextInterval)

	return t
}
