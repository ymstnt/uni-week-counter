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

// Start should be the first Monday of examPeriods, End should be the last Friday
// Periods should be descending, newest period should be at the top
var examPeriods = []Period{
	{Start: date(2025, 12, 15), End: date(2026, 2, 6)},
	{Start: date(2025, 5, 26), End: date(2025, 7, 5)},
	{Start: date(2024, 12, 16), End: date(2025, 2, 8)},
	{Start: date(2024, 5, 20), End: date(2024, 6, 29)},
	{Start: date(2023, 12, 18), End: date(2024, 2, 3)},
}

// Start should be the first Monday of studyPeriods, End should be the last Saturday
// Periods should be descending, newest period should be at the top
var studyPeriods = []Period{
	{Start: date(2025, 9, 8), End: date(2025, 12, 13)},
	{Start: date(2025, 2, 17), End: date(2025, 5, 24)},
	{Start: date(2024, 9, 9), End: date(2024, 12, 14)},
	{Start: date(2024, 2, 12), End: date(2024, 5, 18)},
	{Start: date(2023, 9, 11), End: date(2023, 12, 16)},
}

func isDateInPeriod(date time.Time, period Period) bool {
	periodEnd := period.End.Add(24 * time.Hour) // Include full end day
	return !date.Before(period.Start) && date.Before(periodEnd)
}

func calculateDaysBetween(start, end time.Time) int {
	days := int(end.Sub(start).Hours()/24) + 1
	return days
}

func getPeriods(w http.ResponseWriter, r *http.Request, period []Period) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(period)
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

	numberOnly := r.URL.Query().Get("number-only") != ""

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
			daysLeftBreak := r.URL.Query().Get("days-left-break") != ""

			// make it only work in the summer break
			if daysLeftBreak && int(currentDate.Month()) >= 6 && !numberOnly {
				response = fmt.Sprintf("Break (%d days left)", calculateDaysBetween(currentDate, studyPeriods[0].Start))
				if lang == "hu" {
					response = fmt.Sprintf("Szünet (%d nap van hátra)", calculateDaysBetween(currentDate, studyPeriods[0].Start))
				}
			} else {
				response = "-2";
				if !numberOnly {
					response = "Break"
					if lang == "hu" {
						response = "Szünet"
					}
				}
			}
		}
	} else {
		daysLeftExam := r.URL.Query().Get("days-left-exam") != ""

		if daysLeftExam && !numberOnly {
			response = fmt.Sprintf("Exams - break (%d days left)", calculateDaysBetween(currentDate, studyPeriods[0].Start))
			if lang == "hu" {
				response = fmt.Sprintf("Vizsgaidőszak - szünet (%d nap van hátra)", calculateDaysBetween(currentDate, studyPeriods[0].Start))
			}
		} else {
			response = "-1"
			if !numberOnly {
				response = "Exams - break"
				if lang == "hu" {
					response = "Vizsgaidőszak - szünet"
				}
			}
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
	fmt.Println("Starting...")

	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	} else if len(os.Args) > 1 {
		port = os.Args[1]
	}

	http.HandleFunc("/uwc", getCurrentWeek)
	http.HandleFunc("/study-periods", func(w http.ResponseWriter, r *http.Request) {
		getPeriods(w, r, studyPeriods)
	})
	http.HandleFunc("/exam-periods", func(w http.ResponseWriter, r *http.Request) {
		getPeriods(w, r, examPeriods)
	})

	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()
	fmt.Println("Server is listening on: ", port)

	select {}
}
