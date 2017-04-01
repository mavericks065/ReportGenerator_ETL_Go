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
	doneChannel := make(chan bool)

	go ExtractHealthCenters(healthCenterChannel)

	go TransformToAreaStatistics(healthCenterChannel, areaStatisticsChannel)
	go TransformToMaternityStatistics(healthCenterChannel, maternityStatisticsChannel)
	<-maternityStatisticsChannel
	go LoadAreaStatistics(areaStatisticsChannel, doneChannel)
	<-doneChannel

	fmt.Println(time.Since(start))
}
