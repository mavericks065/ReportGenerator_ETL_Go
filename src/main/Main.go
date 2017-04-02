package main

import (
	. "extractor"
	"fmt"
	. "loader"
	"log"
	"runtime"
	"time"
	. "transformer"

	"github.com/jessevdk/go-flags"
)

type CommandOptions struct {
	Year string `short:"y" long:"year" description:"year of the file to manipulate" required:"true"`
}

func main() {

	start := time.Now()

	runtime.GOMAXPROCS(4)

	//arguments := os.Args[1:]
	options := CommandOptions{}

	//log.Print(arguments)

	parser := initParser(&options)

	_, err := parser.Parse()

	if err != nil {
		log.Print(err)
	}

	run(&options)

	fmt.Println(time.Since(start))
}

func init() {
	var banner = `
	|---------------------------------------------------------------------|
	|																																			|
	|              ETL Health Centers - Extracting data and report        |
	|                                                                     |
	|---------------------------------------------------------------------|
	`

	log.Print(banner)
}

func initParser(options *CommandOptions) (parser *flags.Parser) {
	//default behaviour is HelpFlag | PrintErrors | PassDoubleDash - we need to override the stderr output
	return flags.NewParser(options, flags.HelpFlag)
}

func run(options *CommandOptions) {
	fmt.Println(options.Year)

	healthCenterChannel := make(chan *HealthCenter)

	areaStatisticsChannel := make(chan *AreaStatistics)
	maternityStatisticsChannel := make(chan *MaternityStatistics)

	doneAreaStatisticsChannel := make(chan bool)
	doneMaternityStatisticsChannel := make(chan bool)

	go ExtractHealthCenters(healthCenterChannel, options.Year)

	go TransformToAreaStatistics(healthCenterChannel, areaStatisticsChannel)
	go TransformToMaternityStatistics(healthCenterChannel, maternityStatisticsChannel)

	go LoadAreaStatistics(areaStatisticsChannel, doneAreaStatisticsChannel)
	<-doneAreaStatisticsChannel
	go LoadMaternityStatistics(maternityStatisticsChannel, doneMaternityStatisticsChannel)
	<-doneMaternityStatisticsChannel
}
