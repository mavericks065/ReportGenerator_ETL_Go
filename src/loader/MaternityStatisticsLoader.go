package loader

import (
	"fmt"
	"os"
	"strconv"
	. "transformer"
)

// LoadMaternityStatistics write area numbers i a txt file for now
func LoadMaternityStatistics(maternityStatisticsChannel chan *MaternityStatistics, doneChannel chan bool) {
	destinationFile, _ := os.Create("./maternity-statistics.txt")
	defer destinationFile.Close()

	for a := range maternityStatisticsChannel {
		fmt.Fprintf(destinationFile, "Maternity statistics for center : "+a.HealthCenterName+" in the area : "+strconv.Itoa(a.AreaNumber)+"\n\n")

		fmt.Fprintf(destinationFile, "Maternity level 1 : "+strconv.FormatBool(a.MaternityLevel1)+"\n")
		fmt.Fprintf(destinationFile, "Cesarean level 1 rate : "+strconv.FormatFloat(a.CesareanLevel1Rate, 'f', 2, 64)+"\n")
		fmt.Fprintf(destinationFile, "Maternity level 2 : "+strconv.FormatBool(a.MaternityLevel2)+"\n")
		fmt.Fprintf(destinationFile, "Cesarean level 2 rate : "+strconv.FormatFloat(a.CesareanLevel2Rate, 'f', 2, 64)+"\n")
		fmt.Fprintf(destinationFile, "Maternity level 3 : "+strconv.FormatBool(a.MaternityLevel3)+"\n")
		fmt.Fprintf(destinationFile, "Cesarean level 3 rate : "+strconv.FormatFloat(a.CesareanLevel3Rate, 'f', 2, 64)+"\n")
		fmt.Fprintf(destinationFile, "\n")

	}
	doneChannel <- true
}
