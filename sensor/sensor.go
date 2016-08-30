package sensor

import (
	"errors"
	"time"

	"fmt"

	"github.com/mexicanstrawberry/mcp/recipe"
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

type Timer struct {
	Func *time.Timer
	Chan *time.Timer
}
