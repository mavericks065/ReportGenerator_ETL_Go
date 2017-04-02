package extractor

import (
	"encoding/csv"
	"os"
	"strconv"
)

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

const (
	yes         = "oui"
	emptyRecord = "sans objet"
)

// ExtractHealthCenters extract data from CSV file
func ExtractHealthCenters(ch chan *HealthCenter, year string) {

	f, _ := os.Open("./health_centers_and_hospitals_statistics_" + year + ".csv")
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
	if record[i+33] == emptyRecord {
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
	if record[i+57] == emptyRecord {
		healthCenter.AbortionAverageDelay = 0
	} else {
		healthCenter.AbortionAverageDelay, _ = strconv.ParseFloat(record[i+57], 64)
	}
	return healthCenter
}

func isTrue(condition string) bool {
	if condition == yes {
		return true
	}
	return false
}
