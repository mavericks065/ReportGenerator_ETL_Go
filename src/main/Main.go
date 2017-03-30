package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

func main() {
	healthCenterChannel := make(chan *HealthCenter)
	transformChannel := make(chan *HealthCenter)
	doneChannel := make(chan bool)

	go extract(healthCenterChannel)
	go transform(healthCenterChannel, transformChannel)
	go load(transformChannel, doneChannel)
	<-doneChannel
}

//https://www.dotnetperls.com/csv-go

// HealthCenter is the record type parsed from the CSV file
type HealthCenter struct {
	Name                   string
	Category               string
	City                   string
	Area                   int
	LegalStaus             string
	Scanner                bool
	Mri                    bool
	Camera                 bool
	Tomograph              bool
	MaternityLevel1        bool
	MaternityLevel2        bool
	MaternityLevel3        bool
	CesareanLevel1Rate     float64
	CesareanLevel2Rate     float64
	CesareanLevel3Rate     float64
	Deliveries             int
	AverageMaternityStay   float64
	PediatricEmergency     bool
	Emergency              bool
	IntensiveCareUnit      bool
	BurnCareUnit           bool
	NeuroCareUnit          bool
	Chimio                 bool
	ChimioSessionsCount    int
	Dialysis               bool
	DialysisSessionsCount  int
	Abortion               int
	AbortionMedicalReasons int
	AbortionAverageDelay   float64
}

func extract(ch chan *HealthCenter) {
	fmt.Println("extract")

	f, _ := os.Open("./health_centers_and_hospitals_statistics_2010.csv")
	defer f.Close()
	r := csv.NewReader(f)

	for record, err := r.Read(); err == nil; record, err = r.Read() {

		for i := 0; i < len(record)-1; i = i + 60 {
			healthCenter := new(HealthCenter)
			healthCenter.Name = record[i+1]
			healthCenter.Category = record[i+2]
			fmt.Println(healthCenter.Name)
			fmt.Println(healthCenter.Category)

			ch <- healthCenter
		}
	}

	close(ch)
}

func transform(extractChannel, transformChannel chan *HealthCenter) {
	numMessages := 0

	for ch := range extractChannel {
		numMessages++
		go func(h *HealthCenter) {
			time.Sleep(3 * time.Millisecond)
			h.Name = h.Name + "test"
			// fmt.Println(h.Name)
			transformChannel <- h
			numMessages--
		}(ch)
	}

	for numMessages > 0 {
		time.Sleep(1 * time.Millisecond)
	}

	close(transformChannel)
}

func load(transformChannel chan *HealthCenter, doneChannel chan bool) {
	fmt.Println("transform")
	fmt.Println(len(transformChannel))
	numMessages := 0

	for o := range transformChannel {
		numMessages++
		go func(o *HealthCenter) {
			time.Sleep(1 * time.Millisecond)
			numMessages--
		}(o)
	}
	for numMessages > 0 {
		time.Sleep(1 * time.Millisecond)
	}
	doneChannel <- true
}
