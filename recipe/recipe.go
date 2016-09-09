package recipe

import (
	"encoding/json"
	"errors"
	"time"

	"fmt"
	"math/rand"
)

var RecipeFormat = map[string]Operation{

	"SimpleRecipe": new(SimpleOperation),
}

type Recipe struct {
	Name        string
	Description string
	Format      string
	Operation   Operation
}

type Operation interface {
	Init()
	LoadDummyRecipe()
	LoadRecipe([]byte) error
	CurrentOperation(string, time.Time) (float64, error)
	NextOperation(string, time.Time) (time.Duration, float64, error)
}

func NewRecipe(recipeFormat string, dummy bool) (*Recipe, error) {

	if operation, exist := RecipeFormat[recipeFormat]; exist {

		r := Recipe{
			Name:        "",
			Description: "",
			Format:      "",
			Operation:   operation,
		}

		r.Operation.Init()

		if dummy {
			r.Operation.LoadDummyRecipe()
		}
		return &r, nil

	}

	return nil, errors.New(errorMessages[E_RECIPE_TYPE_UNKNOWN])
}

type OpTable struct {
	Time  time.Time
	Value float64
}

type SimpleOperation struct {
	Sensors map[string][]OpTable
}

func (so *SimpleOperation) Init() {
	so.Sensors = make(map[string][]OpTable)
}

// create some dummy data for testing
func (so *SimpleOperation) LoadDummyRecipe() {

	amount := 100
	interval := 120
	so.Sensors = make(map[string][]OpTable)
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < amount; i++ {

		var opT OpTable

		t := new(time.Time)
		opT.Time = t.Add(time.Duration(i*interval) * time.Second)
		opT.Value = rand.Float64() * 100

		so.Sensors["InsideHumidity"] = append(so.Sensors["InsideHumidity"], opT)
	}
}

func (so *SimpleOperation) LoadRecipe(jsonRecipe []byte) error {
	err := json.Unmarshal(jsonRecipe, so)
	if err != nil {
		return err
	}
	return nil
}

// CurrentOperation returns the current operation value
func (so *SimpleOperation) CurrentOperation(sensorName string, elapsedTime time.Time) (float64, error) {

	if sensorMap, exist := so.Sensors[sensorName]; exist {

		for _, v := range sensorMap {
			if elapsedTime.After(v.Time) {
				return v.Value, nil
			}
		}
		return 0.0, errors.New(errorMessages[E_SENSOR_DISABLED])

	}

	errString := fmt.Sprintf("%s : %s ", errorMessages[E_SENSOR_NOT_FOUND], sensorName)
	return 0.0, errors.New(errString)
}

// NextOperation will search through the receipe operations and return the next operation
// based on elapsed time
func (so *SimpleOperation) NextOperation(sensorName string, elapsedTime time.Time) (time.Duration, float64, error) {

	if sensorMap, exist := so.Sensors[sensorName]; exist {

		for _, v := range sensorMap {

			if !elapsedTime.After(v.Time) {

				// return next time value
				return v.Time.Sub(elapsedTime), v.Value, nil
			}
		}
		// no next operation
		return 0.0, 0.0, errors.New(errorMessages[E_NO_NEXT_OPERATION])

	}
	// sensor name unknown in recipe
	errString := fmt.Sprintf("%s : %s ", errorMessages[E_SENSOR_NOT_FOUND], sensorName)
	return 0.0, 0.0, errors.New(errString)
}
