package sensor

import ()

type Sensor struct {
	Name   string `json:"name"`
	Sensor string `json:"sensor"`
}

func (s *Sensor) Valid() (valid bool) {
	if len(s.Name) == 0 {
		return
	}

	if s.Sensor != "POST" {
		return
	}

	return true
}

type SensorList struct {
	Count int      `json:"count"`
	Items []Sensor `json:"items"`
}
