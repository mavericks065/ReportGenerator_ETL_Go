package main

import (
	. "extractor"
	. "loader"
	. "transformer"
)

func main() {
	healthCenterChannel := make(chan *HealthCenter)
	areaStatisticsChannel := make(chan *AreaStatistics)
	doneChannel := make(chan bool)

	go ExtractHealthCenters(healthCenterChannel)
	go TransformToAreaStatistics(healthCenterChannel, areaStatisticsChannel)
	go LoadAreaStatistics(areaStatisticsChannel, doneChannel)
	<-doneChannel
}
