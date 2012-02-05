package app

import ()

type App struct {
	Name string `json:"name"`
}

func (a *App) Valid() (valid bool) {
	if len(a.Name) == 0 {
		return
	}

	return true
}

type AppList struct {
	Count int   `json:"count"`
	Items []App `json:"items"`
}
