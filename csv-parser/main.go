package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type IntakeLine struct {
	Date   time.Time
	Time   string
	Intake float64
}

type IntakeDay struct {
	Date   time.Time `json:"date"`
	Intake []float64 `json:"intake"`
}

type ByDate []IntakeDay

func (b ByDate) Len() int      { return len(b) }
func (b ByDate) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByDate) Less(i, j int) bool {
	return b[i].Date.Before(b[j].Date)
}

type FinalData struct {
	Max              float64     `json:"max"`
	Min              float64     `json:"min"`
	PercentileMiddle float64     `json:"percentileMiddle"`
	Data             []IntakeDay `json:"data"`
}

func parseLines(fileName string) ([]IntakeLine, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = 0
	reader.Read() // Ignore first line
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var lines []IntakeLine

	for _, rec := range records {
		intake, err := strconv.ParseFloat(strings.Replace(rec[3], ",", ".", -1), 64)
		if err != nil {
			fmt.Println(err)
			continue
		}

		date, err := time.ParseInLocation("02/01/2006", rec[1], time.Local)

		inLn := IntakeLine{date, rec[2], math.Floor(intake * 1000)}
		lines = append(lines, inLn)
	}

	return lines, nil
}

func calculatePercentile(data []float64, percentile float64) float64 {
	sort.Sort(sort.Float64Slice(data))

	index := int(math.Floor(percentile * float64(len(data)) / 100))

	return data[index]
}

type IntakeStatistics struct {
	max              float64
	min              float64
	percentileMiddle float64
}

func calculateStatistics(intakeValues []float64, percentilePoint float64) IntakeStatistics {
	var min, max float64 = 1000000, 0

	for _, intake := range intakeValues {
		if intake > max {
			max = intake
		}
		if intake < min {
			min = intake
		}
	}

	percentileMiddle := calculatePercentile(intakeValues, percentilePoint)

	return IntakeStatistics{max, min, percentileMiddle}
}

func groupIntakesByDate(lines []IntakeLine) map[string]*IntakeDay {
	intakeByDay := make(map[string]*IntakeDay)

	for _, line := range lines {
		date := line.Date.Format("2006-01-02")

		if _, ok := intakeByDay[date]; !ok {
			intakeByDay[date] = &IntakeDay{line.Date, []float64{}}
		}

		intakeByDay[date].Intake = append(intakeByDay[date].Intake, line.Intake)
	}

	return intakeByDay
}

func getAllIntakeValues(intakeByDay map[string]*IntakeDay) []float64 {
	allIntakes := []float64{}

	for _, intake := range intakeByDay {
		allIntakes = append(allIntakes, intake.Intake...)
	}

	return allIntakes
}

func getFormatedIntakeDays(intakeByDay map[string]*IntakeDay) []IntakeDay {
	formatedIntakesDays := []IntakeDay{}

	for _, intake := range intakeByDay {
		formatedIntakesDays = append(formatedIntakesDays, *intake)
	}

	sort.Sort(ByDate(formatedIntakesDays))

	return formatedIntakesDays
}

func main() {
	inLines, err := parseLines("./raw_data.csv")

	if err != nil {
		fmt.Println(err)
		return
	}

	intakeByDay := groupIntakesByDate(inLines)
	allIntakeValues := getAllIntakeValues(intakeByDay)

	intakeStatistics := calculateStatistics(allIntakeValues, 80)

	formatedIntakesDays := getFormatedIntakeDays(intakeByDay)

	finalData := &FinalData{intakeStatistics.max, intakeStatistics.min, intakeStatistics.percentileMiddle, formatedIntakesDays}

	json_data, err := json.MarshalIndent(finalData, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	json_file, err := os.Create("../data.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer json_file.Close()
	json_file.Write(json_data)
}
