package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	go func() {
		for {
			select {
			case <-time.After(time.Second * 5):
				UpdateStatus()
			}
		}
	}()
	http.HandleFunc("/", Index)
	http.ListenAndServe(":8000", nil)
}

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func UpdateStatus() {
	file, err := os.OpenFile("./status.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("error when opening json file")
		return
	}

	rand.Seed(time.Now().UnixNano())
	water := rand.Intn(100) + 1
	wind := rand.Intn(100) + 1

	status := Status{
		Water: water,
		Wind:  wind,
	}

	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")
	encoder.Encode(status)
	if err != nil {
		fmt.Println("error when convert object to bytes")
		return
	}

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("error when writes file")
		return
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/index.html", "./static/layout.html")
	if err != nil {
		panic(err)
		return
	}

	buffer, err := os.ReadFile("./status.json")
	if err != nil {
		fmt.Println("error when opening json file")
		return
	}

	status := Status{}
	reader := bytes.NewReader(buffer)
	err = json.NewDecoder(reader).Decode(&status)
	if err != nil {
		fmt.Println("error when decode json")
		return
	}

	var waterStatus, waterStatusClass string
	switch {
	case status.Water < 5:
		waterStatus = "Aman"
		waterStatusClass = "btn btn-success"
	case status.Water >= 6 && status.Water <= 8:
		waterStatus = "Siaga"
		waterStatusClass = "btn btn-warning"
	case status.Water > 8:
		waterStatus = "Bahaya"
		waterStatusClass = "btn btn-danger"
	}

	var windStatus, windStatusClass string
	switch {
	case status.Wind < 6:
		windStatus = "Aman"
		windStatusClass = "btn btn-success"
	case status.Wind >= 7 && status.Wind <= 15:
		windStatus = "Siaga"
		windStatusClass = "btn btn-warning"
	case status.Wind > 15:
		windStatus = "Bahaya"
		windStatusClass = "btn btn-danger"
	}

	data := map[string]interface{}{
		"Title":            "Wind and Water",
		"Water":            status.Water,
		"Wind":             status.Wind,
		"WaterStatus":      waterStatus,
		"WindStatus":       windStatus,
		"WaterStatusClass": waterStatusClass,
		"WindStatusClass":  windStatusClass,
	}
	tmpl.Execute(w, data)
}
