package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	version string = "v0.0.1"
	boxes          = []string{"5991537e7e280a0010421ba7", "61ec1a8478ce14001bc634f1", "5c8d5386922ca90019e2959d"}
)

// types for receiving
// generated using https://jsonlint.com/json-to-go, on a sample output
type APIResponse struct {
	Id                string          `json:"_id,omitempty"`
	CreatedAt         string          `json:"createdAt,omitempty"`
	UpdatedAt         string          `json:"updatedAt,omitempty"`
	Name              string          `json:"name,omitempty"`
	CurrentLocation   CurrentLocation `json:"currentLocation,omitempty"`
	Exposure          string          `json:"exposure,omitempty"`
	Sensors           []Sensors       `json:"sensors,omitempty"`
	Model             string          `json:"model,omitempty"`
	LastMeasurementAt string          `json:"lastMeasurementAt,omitempty"`
	Grouptag          []string        `json:"grouptag,omitempty"`
	Weblink           string          `json:"weblink,omitempty"`
	Loc               []Loc           `json:"loc,omitempty"`
}

type CurrentLocation struct {
	Type        string    `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
	Timestamp   string    `json:"timestamp,omitempty"`
}

type LastMeasurement struct {
	CreatedAt string `json:"createdAt,omitempty"`
	Value     string `json:"value,omitempty"`
}

type Sensors struct {
	Title           string          `json:"title,omitempty"`
	Unit            string          `json:"unit,omitempty"`
	SensorType      string          `json:"sensorType,omitempty"`
	Icon            string          `json:"icon,omitempty"`
	Id              string          `json:"_id,omitempty"`
	LastMeasurement LastMeasurement `json:"lastMeasurement,omitempty"`
}

type Geometry struct {
	Type        string    `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
	Timestamp   string    `json:"timestamp,omitempty"`
}

type Loc struct {
	Geometry Geometry `json:"geometry,omitempty"`
	Type     string   `json:"type,omitempty"`
}

// fetch_api parses the boxes defined above
func fetch_api() []APIResponse {
	var responses []APIResponse

	for _, box := range boxes {
		// 1. Make the HTTP Request
		request := fmt.Sprintf("https://api.opensensemap.org/boxes/%v?format=json", box)
		println("request: ", request)

		resp, err := http.Get(request) // Replace with your API endpoint
		if err != nil {
			log.Fatalf("Error making request: %v", err)
		}
		defer resp.Body.Close()

		// 2. Handle the Response
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Request failed with status code: %d", resp.StatusCode)
		}

		// 3. Decode the JSON
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v", err)
		}

		var apiResponse APIResponse // Create a variable of the struct type

		err = json.Unmarshal(body, &apiResponse) // Decode the JSON into the struct
		if err != nil {
			log.Fatalf("Error decoding JSON: %v", err)
		}

		responses = append(responses, apiResponse)
	}

	return responses
}

func get_temperatures(responses []APIResponse) []string {
	//2. parse the responses, by getting the sensors with title = "Temperature", get the last measurements, get the value string and
	var temperature_strings []string

	for _, response := range responses {
		for _, sensor := range response.Sensors {
			if sensor.Title == "Temperatur" {
				//get the latest measurement
				temperature_strings = append(temperature_strings, sensor.LastMeasurement.Value)
			}
		}
	}

	return temperature_strings
}

func average_temperature(temperature_strings []string) float64 {
	var temperatures []float64
	for _, temperature_string := range temperature_strings {
		float_t, err := strconv.ParseFloat(temperature_string, 64)
		if err != nil {
			println(fmt.Errorf("could not Parse Float: %v", err))
			continue
		}
		temperatures = append(temperatures, float_t)
	}

	//calc the mean
	sum_temperatures := 0.0
	for _, temperature := range temperatures {
		sum_temperatures += temperature
	}
	return (sum_temperatures / float64(len(temperatures)))
}

func fetch_Temperature() float64 {
	//1. gets responses from the API of 3 sensors
	responses := fetch_api()

	//2. parse the responses, by getting the sensors with title = "Temperature", get the last measurements, get the value string and
	temperature_strings := get_temperatures(responses)
	fmt.Printf("%v", temperature_strings)

	//3. convert the value strings to float64, average the values
	avgtemperature := average_temperature(temperature_strings)

	return avgtemperature
}

// ----- methods for presenting endpoints

func start_enpoints() {
	r := gin.Default()

	// Version endpoint
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version": version,
		})
	})

	// Temperature endpoint (simulated)
	r.GET("/temperature", func(c *gin.Context) {
		// Simulate fetching temperature from a sensor or API
		temperature := fetch_Temperature()
		c.JSON(http.StatusOK, gin.H{
			"temperature": temperature,
			"unit":        "Celsius",
		})
	})

	fmt.Println("Server listening on port 8080...")
	r.Run(":8080") // Listen and serve on 0.0.0.0:8080
}

// ----- main methods

func print_version() {
	println(version)
}

func main() {
	print_version()
	start_enpoints()
	// start_enpoints()
}
