package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var(
	version string = "v0.0.1"
)

// types for receiving
//generated using https://jsonlint.com/json-to-go, on a sample output
type APIResponse struct {
	Id string `json:"_id,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Name string `json:"name,omitempty"`
	CurrentLocation CurrentLocation `json:"currentLocation,omitempty"`
	Exposure string `json:"exposure,omitempty"`
	Sensors []Sensors `json:"sensors,omitempty"`
	Model string `json:"model,omitempty"`
	LastMeasurementAt string `json:"lastMeasurementAt,omitempty"`
	Grouptag []string `json:"grouptag,omitempty"`
	Weblink string `json:"weblink,omitempty"`
	Loc []Loc `json:"loc,omitempty"`
}

type CurrentLocation struct {
	Type string `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type LastMeasurement struct {
	CreatedAt string `json:"createdAt,omitempty"`
	Value string `json:"value,omitempty"`
}

type Sensors struct {
	Title string `json:"title,omitempty"`
	Unit string `json:"unit,omitempty"`
	SensorType string `json:"sensorType,omitempty"`
	Icon string `json:"icon,omitempty"`
	Id string `json:"_id,omitempty"`
	LastMeasurement LastMeasurement `json:"lastMeasurement,omitempty"`
}

type Geometry struct {
	Type string `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type Loc struct {
	Geometry Geometry `json:"geometry,omitempty"`
	Type string `json:"type,omitempty"`
}

func parse_Sensor()  {
	// parse an example query
}

// ----- methods for presenting endpoints

func start_enpoints(){
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
		temperature := 25.5 + float64(time.Now().UnixNano()%100) // Add some randomness
		c.JSON(http.StatusOK, gin.H{
			"temperature": temperature,
			"unit":        "Celsius",
		})
	})

	fmt.Println("Server listening on port 8080...")
	r.Run(":8090") // Listen and serve on 0.0.0.0:8080
}

// ----- main methods

func print_version()  {
	println(version)	
}

func main()  {
	print_version()
	start_enpoints()
}
