package main

import (
	"container/list"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
)

func main() {
	healthCenterChannel := make(chan *HealthCenter)
	areaStatisticsChannel := make(chan *AreaStatistics)
	doneChannel := make(chan bool)

	go extract(healthCenterChannel)
	go transform(healthCenterChannel, areaStatisticsChannel)
	go load(areaStatisticsChannel, doneChannel)
	<-doneChannel
}

// HealthCenter is the record type parsed from the CSV file
type HealthCenter struct {
	Name                    string
	Category                string
	City                    string
	Area                    int
	LegalStatus             string
	Scanner                 bool
	Mri                     bool
	Camera                  bool
	Tomograph               bool
	MaternityLevel1         bool
	MaternityLevel2         bool
	MaternityLevel3         bool
	CesareanLevel1Rate      float64
	CesareanLevel2Rate      float64
	CesareanLevel3Rate      float64
	Deliveries              int
	AverageMaternityStay    float64
	PediatricEmergency      bool
	Emergency               bool
	IntensiveCareUnit       bool
	IntensiveCareUnitLevel2 bool
	BurnCareUnit            bool
	NeuroCareUnit           bool
	HeartCareUnit           bool
	Chimio                  bool
	ChimioSessionsCount     int
	Dialysis                bool
	DialysisSessionsCount   int
	Abortion                int
	AbortionMedicalReasons  int
	AbortionAverageDelay    float64
}

// AreaStatistics represent the interesting statistics of
// a specific area
type AreaStatistics struct {
	Number                    int
	HealthCenters             map[string]*HealthCenter
	MriCount                  int
	Scanner                   int
	MaternityServicesCount    int
	ChimioCenters             int
	BurnCenters               int
	EmergencyCount            int
	AbortionTotalCount        int
	MedicalAbortionTotalCount int
}

func extract(ch chan *HealthCenter) {
	fmt.Println("extract")

	f, _ := os.Open("./health_centers_and_hospitals_statistics_2010.csv")
	defer f.Close()
	r := csv.NewReader(f)

	for record, err := r.Read(); err == nil; record, err = r.Read() {

		for i := 0; i < len(record)-1; i = i + 60 {
			healthCenter := parseRecord(record, i)
			ch <- healthCenter
		}
	}

	close(ch)
}

func transform(extractChannel chan *HealthCenter, areaStatisticsChannel chan *AreaStatistics) {

	numMessages := 0
	areaStatisticsMap := make(map[int]list.List)

	// fill map od area statistics
	for ch := range extractChannel {
		numMessages++
		go func(h *HealthCenter) {
			time.Sleep(3 * time.Millisecond)
			mutex.Lock()
			defer mutex.Unlock()
			if valueList, isPresent := areaStatisticsMap[h.Area]; isPresent {
				valueList.PushBack(h)
			} else {
				var healthCenterList list.List
				healthCenterList.PushBack(h)
				areaStatisticsMap[h.Area] = healthCenterList
			}
			// transformChannel <- h
			numMessages--
		}(ch)
	}

	// wait for the map
	for numMessages > 0 {
		time.Sleep(1 * time.Millisecond)
	}

	// build areaStatistics
	// Fill the area statistics channel
	for areaNumber := range areaStatisticsMap {
		listOfHcs := areaStatisticsMap[areaNumber]
		areaStatistics := buildAreaStatistics(areaNumber, &listOfHcs)

		areaStatisticsChannel <- areaStatistics
	}
	close(areaStatisticsChannel)
}

func load(areaStatisticsChannel chan *AreaStatistics, doneChannel chan bool) {
	destinationFile, _ := os.Create("./dest.txt")
	defer destinationFile.Close()

	numMessages := 0

	for o := range areaStatisticsChannel {
		numMessages++
		go func(o *AreaStatistics) {
			time.Sleep(1 * time.Millisecond)
			fmt.Fprintf(destinationFile, strconv.Itoa(o.Number)+"\n")
			numMessages--
		}(o)
	}
	for numMessages > 0 {
		time.Sleep(1 * time.Millisecond)
	}
	doneChannel <- true
}

// NewAreaStatistics instantiate AreaStatistics with its map of HealthCenter
func NewAreaStatistics() *AreaStatistics {
	var areaStatistics AreaStatistics
	areaStatistics.HealthCenters = make(map[string]*HealthCenter)
	return &areaStatistics
}

func buildAreaStatistics(areaNumber int, listOfHcs *list.List) *AreaStatistics {
	areaStatistics := NewAreaStatistics()
	areaStatistics.Number = areaNumber

	for e := listOfHcs.Front(); e != nil; e = e.Next() {
		if e.Value != nil {
			item := e.Value
			hc := item.(*HealthCenter)

			areaStatistics.HealthCenters[hc.Name] = hc
			if hc.Mri {
				areaStatistics.MriCount++
			}
			if hc.Scanner {
				areaStatistics.Scanner++
			}
			if hc.MaternityLevel1 || hc.MaternityLevel2 || hc.MaternityLevel3 {
				areaStatistics.MaternityServicesCount++
			}
			if hc.Chimio {
				areaStatistics.ChimioCenters++
			}
			if hc.BurnCareUnit {
				areaStatistics.BurnCenters++
			}
			if hc.Emergency {
				areaStatistics.EmergencyCount++
			}
			areaStatistics.AbortionTotalCount += hc.Abortion
			areaStatistics.MedicalAbortionTotalCount += hc.AbortionMedicalReasons
		}
	}
	return areaStatistics
}

func isTrue(condition string) bool {
	if condition == "oui" {
		return true
	}
	return false
}

// parseRecord to create HealthCenter
func parseRecord(record []string, i int) *HealthCenter {
	healthCenter := new(HealthCenter)
	healthCenter.Name = record[i+1]
	healthCenter.Category = record[i+2]
	healthCenter.City = record[i+3]
	healthCenter.Area, _ = strconv.Atoi(record[i+4])
	healthCenter.LegalStatus = record[i+6]
	healthCenter.Scanner = isTrue(record[i+19])
	healthCenter.Mri = isTrue(record[i+20])
	healthCenter.Camera = isTrue(record[i+21])
	healthCenter.Tomograph = isTrue(record[i+22])
	healthCenter.MaternityLevel1 = isTrue(record[i+26])
	healthCenter.MaternityLevel2 = isTrue(record[i+28])
	healthCenter.MaternityLevel3 = isTrue(record[i+30])
	healthCenter.CesareanLevel1Rate, _ = strconv.ParseFloat(record[i+27], 64)
	healthCenter.CesareanLevel2Rate, _ = strconv.ParseFloat(record[i+29], 64)
	healthCenter.CesareanLevel3Rate, _ = strconv.ParseFloat(record[i+31], 64)
	healthCenter.Deliveries, _ = strconv.Atoi(record[i+32])
	if record[i+33] == "sans objet" {
		healthCenter.AverageMaternityStay = 0
	} else {
		healthCenter.AverageMaternityStay, _ = strconv.ParseFloat(record[i+33], 64)
	}
	healthCenter.PediatricEmergency = isTrue(record[i+34])
	healthCenter.Emergency = isTrue(record[i+36])
	healthCenter.IntensiveCareUnit = isTrue(record[i+38])
	healthCenter.IntensiveCareUnitLevel2 = isTrue(record[i+42])
	healthCenter.BurnCareUnit = isTrue(record[i+46])
	healthCenter.HeartCareUnit = isTrue(record[i+47])
	healthCenter.NeuroCareUnit = isTrue(record[i+48])
	healthCenter.Chimio = isTrue(record[i+49])
	healthCenter.ChimioSessionsCount, _ = strconv.Atoi(record[i+50])
	healthCenter.Dialysis = isTrue(record[i+51])
	healthCenter.DialysisSessionsCount, _ = strconv.Atoi(record[i+52])
	healthCenter.Abortion, _ = strconv.Atoi(record[i+55])
	healthCenter.AbortionMedicalReasons, _ = strconv.Atoi(record[i+56])
	if record[i+57] == "sans objet" {
		healthCenter.AbortionAverageDelay = 0
	} else {
		healthCenter.AbortionAverageDelay, _ = strconv.ParseFloat(record[i+57], 64)
	}
	return healthCenter
}
