package app

import ()

type App struct {
	Name string `json:"name"`
}

type AppList struct {
	Count int   `json:"count"`
	Items []App `json:"items"`
}

func Get(name string) App {
	return App{name}
}

func GetList() AppList {
	return AppList{1, []App{App{"testapp"}}}
}

func Put(name string) (err error) {
	return
}

func Post(name string, signal string) (err error) {
	return
}

func Delete(name string) (err error) {
	return
}
