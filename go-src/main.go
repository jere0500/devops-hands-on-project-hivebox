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

// APIResponse main API response of opensensemap
type APIResponse struct {
	ID                string          `json:"_id,omitempty"`
	CreatedAt         string          `json:"createdAt,omitempty"`
	UpdatedAt         string          `json:"updatedAt,omitempty"`
	Name              string          `json:"name,omitempty"`
	CurrentLocation   currentLocation `json:"currentLocation,omitempty"`
	Exposure          string          `json:"exposure,omitempty"`
	Sensors           []sensors       `json:"sensors,omitempty"`
	Model             string          `json:"model,omitempty"`
	LastMeasurementAt string          `json:"lastMeasurementAt,omitempty"`
	Grouptag          []string        `json:"grouptag,omitempty"`
	Weblink           string          `json:"weblink,omitempty"`
	Loc               []loc           `json:"loc,omitempty"`
}

type currentLocation struct {
	Type        string    `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
	Timestamp   string    `json:"timestamp,omitempty"`
}

type lastMeasurement struct {
	CreatedAt string `json:"createdAt,omitempty"`
	Value     string `json:"value,omitempty"`
}

type sensors struct {
	Title           string          `json:"title,omitempty"`
	Unit            string          `json:"unit,omitempty"`
	SensorType      string          `json:"sensorType,omitempty"`
	Icon            string          `json:"icon,omitempty"`
	ID              string          `json:"_id,omitempty"`
	LastMeasurement lastMeasurement `json:"lastMeasurement,omitempty"`
}

type geometry struct {
	Type        string    `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
	Timestamp   string    `json:"timestamp,omitempty"`
}

type loc struct {
	Geometry geometry `json:"geometry,omitempty"`
	Type     string   `json:"type,omitempty"`
}

// fetchAPI parses the boxes defined above
func fetchAPI() []APIResponse {
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

func getTemperatures(responses []APIResponse) []string {
	//2. parse the responses, by getting the sensors with title = "Temperature", get the last measurements, get the value string and
	var temperatureStrings []string

	for _, response := range responses {
		for _, sensor := range response.Sensors {
			if sensor.Title == "Temperatur" {
				//get the latest measurement
				temperatureStrings = append(temperatureStrings, sensor.LastMeasurement.Value)
			}
		}
	}

	return temperatureStrings
}

func averageTemperature(temperatureStrings []string) float64 {
	var temperatures []float64
	for _, temperatureString := range temperatureStrings {
		floatTemperatur, err := strconv.ParseFloat(temperatureString, 64)
		if err != nil {
			println(fmt.Errorf("could not Parse Float: %v", err))
			continue
		}
		temperatures = append(temperatures, floatTemperatur)
	}

	//calc the mean
	sumTemperatures := 0.0
	for _, temperature := range temperatures {
		sumTemperatures += temperature
	}
	return (sumTemperatures / float64(len(temperatures)))
}

func fetchTemperature() float64 {
	//1. gets responses from the API of 3 sensors
	responses := fetchAPI()

	//2. parse the responses, by getting the sensors with title = "Temperature", get the last measurements, get the value string and
	temperatureStrings := getTemperatures(responses)
	fmt.Printf("%v", temperatureStrings)

	//3. convert the value strings to float64, average the values
	avgtemperature := averageTemperature(temperatureStrings)

	return avgtemperature
}

// ----- methods for presenting endpoints

func startEnpoints() {
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
		temperature := fetchTemperature()
		c.JSON(http.StatusOK, gin.H{
			"temperature": temperature,
			"unit":        "Celsius",
		})
	})

	fmt.Println("Server listening on port 8080...")
	r.Run(":8080") // Listen and serve on 0.0.0.0:8080
}

// ----- main methods

func printVersion() {
	println(version)
}

func main() {
	printVersion()
	startEnpoints()
	// start_enpoints()
}
