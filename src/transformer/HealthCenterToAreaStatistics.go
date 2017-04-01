package transformer

import (
	"container/list"
	. "extractor"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
)

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

// TransformToAreaStatistics transform HealthCenter extracted previously into
// AreaStatistics and puts them in the channel
func TransformToAreaStatistics(extractChannel chan *HealthCenter,
	areaStatisticsChannel chan *AreaStatistics) {

	numMessages := 0
	areaStatisticsMap := make(map[int]list.List)

	// fill map and area statistics
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

// newAreaStatistics instantiate AreaStatistics with its map of HealthCenter
func newAreaStatistics() *AreaStatistics {
	var areaStatistics AreaStatistics
	areaStatistics.HealthCenters = make(map[string]*HealthCenter)
	return &areaStatistics
}

func buildAreaStatistics(areaNumber int, listOfHcs *list.List) *AreaStatistics {
	areaStatistics := newAreaStatistics()
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
