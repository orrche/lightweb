package main

import (
	"fmt"
	"net/http"
	"strconv"

	"minoris.se/rabbitmq/camq"

	"github.com/gorilla/mux"
)

var amq *camq.Amq

type LightCommand struct {
	Name    string
	Command string
}

type Light struct {
	Name     string
	Id       string
	Commands []LightCommand
}

func (light Light) RunCommand(index int) {
	body := fmt.Sprintf("%s", light.Commands[index].Command)
	amq.Publish(body)
}

func GetAllLights() []Light {
	return config.Lights
}

func failOnError(err error, msg string) {
}

func GetLight(id string) (Light, bool) {
	for _, light := range GetAllLights() {
		if light.Id == id {
			return light, true
		}
	}

	return (*new(Light)), false
}

func lightInitiate(id string, command string) {
	light, ok := GetLight(id)
	if ok {
		index, err := strconv.Atoi(command)
		if err == nil {
			light.RunCommand(index)
		}
	}
}

func lightInit(router *mux.Router) {
	amq = camq.GetAMQChannel(config.AMQ)
	lightTemplate := readTemplateFile("template/lightindex.html")
	lightApikeyCommand := func(w http.ResponseWriter, r *http.Request) {
		if mux.Vars(r)["apikey"] != "lazyasfuck" {
			return
		}

		lightInitiate(mux.Vars(r)["id"], mux.Vars(r)["command"])
	}
	lightCommand := func(w http.ResponseWriter, r *http.Request) {
		lightInitiate(mux.Vars(r)["id"], mux.Vars(r)["command"])
		http.Redirect(w, r, "/", 302)
	}
	lightIndex := func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Title  string
			Lights []Light
		}{
			Title:  "Lights",
			Lights: GetAllLights(),
		}
		err := lightTemplate.Execute(w, data)
		failOnErr(err, w, r)
	}

	router.HandleFunc("/", lightIndex)
	router.HandleFunc("/cmd/{id:[a-z0-9_]*}/{command:[0-9]*}", lightCommand)
	router.HandleFunc("/apikey/cmd/{apikey:[a-z0-9_]*}/{id:[a-z0-9_]*}/{command:[0-9]*}", lightApikeyCommand)
}
