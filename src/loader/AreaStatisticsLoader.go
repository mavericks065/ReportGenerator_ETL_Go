package loader

import (
	"fmt"
	"os"
	"strconv"
	. "transformer"
)

// LoadAreaStatistics write area numbers i a txt file for now
func LoadAreaStatistics(areaStatisticsChannel chan *AreaStatistics, doneChannel chan bool) {
	destinationFile, _ := os.Create("./areaStatisticsDest.txt")
	defer destinationFile.Close()

	for a := range areaStatisticsChannel {
		fmt.Fprintf(destinationFile, "Area : "+strconv.Itoa(a.Number)+"\n\n")
		fmt.Fprintf(destinationFile, "Number of Mris : "+strconv.Itoa(a.MriCount)+"\n")
		fmt.Fprintf(destinationFile, "Number of Scanners : "+strconv.Itoa(a.Scanner)+"\n")
		fmt.Fprintf(destinationFile, "Number of Maternity services : "+strconv.Itoa(a.MaternityServicesCount)+"\n")
		fmt.Fprintf(destinationFile, "Number of Chimio services : "+strconv.Itoa(a.ChimioCenters)+"\n")
		fmt.Fprintf(destinationFile, "Number of Burn services : "+strconv.Itoa(a.BurnCenters)+"\n")
		fmt.Fprintf(destinationFile, "Abortions : "+strconv.Itoa(a.AbortionTotalCount)+"\n")
		fmt.Fprintf(destinationFile, "Medical Abortions : "+strconv.Itoa(a.MedicalAbortionTotalCount)+"\n")
		fmt.Fprintf(destinationFile, "List of Health centers in this area : "+strconv.Itoa(len(a.HealthCenters))+"\n")
		for k := range a.HealthCenters {
			fmt.Fprintf(destinationFile, k+"\n")
		}
		fmt.Fprintf(destinationFile, "\n")

	}
	doneChannel <- true
}
