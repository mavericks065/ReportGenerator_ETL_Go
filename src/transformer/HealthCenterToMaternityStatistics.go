package transformer

import (
	. "extractor"
	"time"
)

// MaternityStatistics represent the maternity statistics of
// a health center
type MaternityStatistics struct {
	AreaNumber         int
	HealthCenterName   string
	MaternityLevel1    bool
	MaternityLevel2    bool
	MaternityLevel3    bool
	CesareanLevel1Rate float64
	CesareanLevel2Rate float64
	CesareanLevel3Rate float64
}

// TransformToMaternityStatistics transform HealthCenter extracted previously into
// AreaStatistics and puts them in the channel
func TransformToMaternityStatistics(extractChannel chan *HealthCenter,
	maternityStatisticsChannel chan *MaternityStatistics) {

	numMessages := 0

	// fill maternity statistics
	for ch := range extractChannel {
		numMessages++
		go func(h *HealthCenter) {
			if h.MaternityLevel1 || h.MaternityLevel2 || h.MaternityLevel3 {
				maternityStatistics := new(MaternityStatistics)
				maternityStatistics.AreaNumber = h.Area
				maternityStatistics.HealthCenterName = h.Name
				maternityStatistics.MaternityLevel1 = h.MaternityLevel1
				maternityStatistics.MaternityLevel2 = h.MaternityLevel2
				maternityStatistics.MaternityLevel3 = h.MaternityLevel3
				maternityStatistics.CesareanLevel1Rate = h.CesareanLevel1Rate
				maternityStatistics.CesareanLevel2Rate = h.CesareanLevel2Rate
				maternityStatistics.CesareanLevel3Rate = h.CesareanLevel3Rate

				maternityStatisticsChannel <- maternityStatistics
			}
			numMessages--
		}(ch)
	}

	// wait for the map
	for numMessages > 0 {
		time.Sleep(1 * time.Millisecond)
	}
	close(maternityStatisticsChannel)
}
