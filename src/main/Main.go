package main

import (
	. "extractor"
	"fmt"
	. "loader"
	"runtime"
	"time"
	. "transformer"
)

func main() {

	start := time.Now()

	runtime.GOMAXPROCS(4)

	healthCenterChannel := make(chan *HealthCenter)

	areaStatisticsChannel := make(chan *AreaStatistics)
	maternityStatisticsChannel := make(chan *MaternityStatistics)

	doneAreaStatisticsChannel := make(chan bool)
	doneMaternityStatisticsChannel := make(chan bool)

	go ExtractHealthCenters(healthCenterChannel)

	go TransformToAreaStatistics(healthCenterChannel, areaStatisticsChannel)
	go TransformToMaternityStatistics(healthCenterChannel, maternityStatisticsChannel)

	go LoadAreaStatistics(areaStatisticsChannel, doneAreaStatisticsChannel)
	<-doneAreaStatisticsChannel
	go LoadMaternityStatistics(maternityStatisticsChannel, doneMaternityStatisticsChannel)
	<-doneMaternityStatisticsChannel

	fmt.Println(time.Since(start))
}
