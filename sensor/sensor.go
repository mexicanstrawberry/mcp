package sensor

import (
	"errors"
	"time"

	"fmt"

	"github.com/mexicanstrawberry/mcp/events"
	"github.com/mexicanstrawberry/mcp/gatekeeper"
	"github.com/mexicanstrawberry/mcp/recipe"
	clog "github.com/morriswinkler/cloudglog"
)

var SensorType = map[string]Sensor{
	"InsideHumidity": new(InsideHumidity),
}

const (
	defaultTickInterval = 2
)

type Sensor interface {
	Init()
	Name() string
	Set(float64) error
	Get() error
	SetRecipe(*recipe.Recipe)
	Run()
}

func NewSensor(sensorType string, r *recipe.Recipe) (Sensor, error) {

	if sensor, exist := SensorType[sensorType]; exist {
		sensor.Init()
		sensor.SetRecipe(r)
		return sensor, nil
	}

	erroString := fmt.Sprintf("sensor unknown %s", sensorType)
	return nil, errors.New(erroString)
}

type InsideHumidity struct {
	SensorName   string
	Recipe       *recipe.Recipe
	StartTime    time.Time
	Time         time.Time
	TargetValue  float64
	CurrentValue float64
	Timer        Timer
	Ticker       *time.Ticker
}

type Timer struct {
	Func *time.Timer
	Chan *time.Timer
}

func (ih *InsideHumidity) Init() {
	ih.SensorName = "InsideHumidity"
}

func (ih *InsideHumidity) Name() string {
	return ih.SensorName
}

func (ih *InsideHumidity) Set(val float64) error {

	clog.Infoln("set Inside Huminity")
	return nil
}

func (ih *InsideHumidity) Get() error {

	clog.Infoln("get Inside Huminity")
	return nil
}

func (ih *InsideHumidity) SetRecipe(r *recipe.Recipe) {
	ih.Recipe = r
}

func (ih *InsideHumidity) regulate() {
	if insideHumidity, exist := gatekeeper.CurrentData[ih.SensorName]; exist {
		if insideHumidity == nil {
			return
		}
		offsetTargetInside := ih.TargetValue - insideHumidity.(float64)
		offsetOutsideInside := 0.0
		if outsideHumidity, exist := gatekeeper.CurrentData["OutsideHumidity"]; exist {
			if outsideHumidity != nil {
				offsetOutsideInside = outsideHumidity.(float64) - insideHumidity.(float64)
			}
		}

		// To dry and outside as more humidity
		if offsetTargetInside > 5 && offsetOutsideInside > 0 {
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Hatch",
				Value:     0,
				Priority:  events.Priority(events.P_MIN),
			}
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "OutsideFan",
				Value:     5,
				Priority:  events.Priority(events.P_NORM),
			}
		}
		// Much to dry
		if offsetTargetInside > 10 {
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Humidifier",
				Value:     5,
				Priority:  events.Priority(events.P_NORM),
			}
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "InsideFan",
				Value:     7,
				Priority:  events.Priority(events.P_NORM),
			}
		}
		// To much humidity
		if offsetTargetInside < -5 && offsetOutsideInside < 0 {
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Hatch",
				Value:     1,
				Priority:  events.Priority(events.P_MIN),
			}
		}
		// Much to much humidity
		if offsetTargetInside < -10 && offsetOutsideInside < 0 {
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "Hatch",
				Value:     1,
				Priority:  events.Priority(events.P_MIN),
			}
			events.Channel <- events.MqttCommand{
				CommandID: 1,
				Actuator:  "OutsideFan",
				Value:     7,
				Priority:  events.Priority(events.P_NORM),
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
			clog.Infoln("[Ticker] ", t)
			ih.regulate()

		case t := <-ih.Timer.Chan.C:
			clog.Infoln("[Timer] ", t)
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
		clog.V(3).Infoln("[nextTimer] fire")
		ih.TargetValue = nextValue
	})
	// set receiver channel
	t.Chan = time.NewTimer(nextInterval)

	return t
}
