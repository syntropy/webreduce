package sensor

import ()

type Sensor struct {
	Name string `json:"name"`
}

func (s *Sensor) Valid() (valid bool) {
	if len(s.Name) == 0 {
		return
	}

	return true
}

type SensorList struct {
	Count int      `json:"count"`
	Items []Sensor `json:"items"`
}
