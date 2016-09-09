package sensor

import (
	"time"

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
		offsetTargetInside := ih.TargetValue - insideTemperature.(float64)
		offsetOutsideInside := 0.0
		if outsideTemperature, exist := gatekeeper.CurrentData["OutsideTemperature"]; exist {
			if outsideTemperature != nil {
				offsetOutsideInside = outsideTemperature.(float64) - insideTemperature.(float64)
			}
		}
		//clog.Info(offsetTargetInside)
		//clog.Info(offsetOutsideInside)
	}

}

func (ih *InsideTemperature) Run() {

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
