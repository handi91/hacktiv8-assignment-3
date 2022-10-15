package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type InputData struct {
	Status Status `json:"status"`
}

func reload() {
	const (
		min = 1
		max = 100
	)

	for {
		water := rand.Intn(max-min) + min
		wind := rand.Intn(max-min) + min

		status := Status{
			Water: water,
			Wind:  wind,
		}

		input := InputData{
			Status: status,
		}

		inputJson, err := json.MarshalIndent(&input, "", "  ")
		if err != nil {
			log.Fatalln("error marshal input data", err)
		}

		err = ioutil.WriteFile("data.json", inputJson, 0644)
		if err != nil {
			log.Fatalln("error auto reload data.json file", err)
		}

		time.Sleep(15 * time.Second)
	}
}

func main() {
	go reload()

	http.HandleFunc("/", renderTemplate)

	http.ListenAndServe(":8080", nil)
}

func renderTemplate(w http.ResponseWriter, r *http.Request) {
	jsonData, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Fprintln(w, "error load data from data.json file")
		return
	}

	var data InputData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Fprintln(w, "error unmarshal data")
		return
	}

	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Fprintln(w, "error parse index.html")
		return
	}

	var waterStatus string
	var windStatus string
	water := data.Status.Water
	wind := data.Status.Wind

	if water > 8 {
		waterStatus = "bahaya"
	} else if water > 5 {
		waterStatus = "siaga"
	} else {
		waterStatus = "aman"
	}

	if wind > 15 {
		windStatus = "bahaya"
	} else if wind > 6 {
		windStatus = "siaga"
	} else {
		windStatus = "aman"
	}

	response := map[string]interface{}{
		"water":       water,
		"wind":        wind,
		"waterStatus": waterStatus,
		"windStatus":  windStatus,
	}

	t.Execute(w, response)
}
