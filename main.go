package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Period struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

var examPeriods = []Period{
	{Start: date(2023, 12, 18), End: date(2024, 2, 3)},
	{Start: date(2024, 5, 20), End: date(2024, 6, 29)},
	{Start: date(2024, 12, 16), End: date(2025, 2, 8)},
	{Start: date(2025, 5, 26), End: date(2025, 7, 5)},
}

var studyPeriods = []Period{
	{Start: date(2023, 9, 11), End: date(2023, 12, 16)},
	{Start: date(2024, 2, 12), End: date(2024, 5, 18)},
	{Start: date(2024, 9, 9), End: date(2024, 12, 14)},
	{Start: date(2025, 2, 17), End: date(2025, 5, 24)},
}

func isDateInPeriod(date time.Time, period Period) bool {
	periodEnd := period.End.Add(24 * time.Hour) // Include full end day
	return !date.Before(period.Start) && date.Before(periodEnd)
}

func getCurrentWeek(w http.ResponseWriter, r *http.Request) {
	currentDate := time.Now()
	//currentDate := date(2025, 2, 17)
	var response string

	isInExamPeriod := false
	var firstStudyPeriodStart time.Time

	lang := r.URL.Query().Get("lang")
	if lang != "hu" {
		lang = "en"
	}

	numberOnly := r.URL.Query().Get("numberOnly") == "true"

	for _, period := range examPeriods {
		if isDateInPeriod(currentDate, period) {
			isInExamPeriod = true
			break
		}
	}

	if !isInExamPeriod {
		for _, period := range studyPeriods {
			if isDateInPeriod(currentDate, period) {
				firstStudyPeriodStart = period.Start
				break
			}
		}

		if !firstStudyPeriodStart.IsZero() {
			weeksPassed := int(currentDate.Sub(firstStudyPeriodStart).Hours()/(24*7)) + 1
			if numberOnly {
				response = fmt.Sprintf("%d", weeksPassed)
			} else {
				suffix := getSuffix(weeksPassed)
				if lang == "hu" {
					suffix = "."
				}
				response = fmt.Sprintf("%d%s", weeksPassed, suffix)
			}
		} else {
			response = "Break"
			if lang == "hu" {
				response = "Szünet"
			}
		}
	} else {
		response = "Exams - break"
		if lang == "hu" {
			response = "Vizsgaidőszak - szünet"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": response})
}

func getSuffix(weekNum int) string {
	var suffix string
	switch weekNum {
	case 1:
		suffix = "st"
	case 2:
		suffix = "nd"
	case 3:
		suffix = "rd"
	default:
		suffix = "th"
	}

	return suffix
}

func main() {
	port := "8080"

	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	} else if len(os.Args) > 1 {
		port = os.Args[1]
	}

	http.HandleFunc("/uniWeekCount", getCurrentWeek)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	fmt.Println("Server is listening on:" + port)
}
