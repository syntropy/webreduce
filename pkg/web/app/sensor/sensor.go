package sensor

import (
	"fmt"
	"wr/interpreter/lua"
)

// Sensors represent a persistable behaviour that collect and/or emit data
type Sensor struct {
	Name     string `json:"name"`
	Language string `json:"language"`
	Code     string `json:"code"`
}

// Validates the sensor.
func (a *Sensor) Valid() bool {
	if len(a.Name) < 1 {
		return false
	}

	if a.Language != "lua" {
		return false
	}

	if _, err := lua.New().Eval(a.Code); err != nil {
		return false
	}

	return true
}

// Calls the sensor with data
func (a *Sensor) Call(data interface{}) (err error) {
	lctx := lua.New()
	lctx.RegisterEmitCallback(func(data []byte) { fmt.Printf("EMIT: %v\n", string(data)) })

	fn, err := lctx.Eval(a.Code)
	if err != nil {
		return
	}

	fn(data.([]byte), []byte{})
	return
}

// SensorList represents a list of persisted Sensors
type SensorList struct {
	Count int      `json:"count"`
	Items []Sensor `json:"items"`
}
