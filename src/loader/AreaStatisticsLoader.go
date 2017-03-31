package loader

import (
	"fmt"
	"os"
	"strconv"
	"time"
	. "transformer"
)

// LoadAreaStatistics write area numbers i a txt file for now
func LoadAreaStatistics(areaStatisticsChannel chan *AreaStatistics, doneChannel chan bool) {
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
