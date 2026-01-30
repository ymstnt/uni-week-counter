package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Period struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Semester string    `json:"semester"`
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// Start should be the first Monday of exam periods, End should be the last Friday
// Periods should be descending, newest period should be at the top
var examPeriods = []Period{
	{Start: date(2026, 5, 26), End: date(2026, 7, 4), Semester: "2025/26/2"},
	{Start: date(2025, 12, 15), End: date(2026, 2, 6), Semester: "2025/26/1"},
	{Start: date(2025, 5, 26), End: date(2025, 7, 5), Semester: "2024/25/2"},
	{Start: date(2024, 12, 16), End: date(2025, 2, 8), Semester: "2024/25/1"},
	{Start: date(2024, 5, 20), End: date(2024, 6, 29), Semester: "2023/24/2"},
	{Start: date(2023, 12, 18), End: date(2024, 2, 3), Semester: "2023/24/1"},
}

// Start should be the first Monday of study periods (including registration week), End should be the last Saturday
// Periods should be descending, newest period should be at the top
var studyPeriods = []Period{
	{Start: date(2026, 2, 9), End: date(2026, 5, 23), Semester: "2025/26/2"},
	{Start: date(2025, 9, 1), End: date(2025, 12, 13), Semester: "2025/26/1"},
	{Start: date(2025, 2, 10), End: date(2025, 5, 24), Semester: "2024/25/2"},
	{Start: date(2024, 9, 2), End: date(2024, 12, 14), Semester: "2024/25/1"},
	{Start: date(2024, 2, 5), End: date(2024, 5, 18), Semester: "2023/24/2"},
	{Start: date(2023, 9, 4), End: date(2023, 12, 16), Semester: "2023/24/1"},
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
	//currentDate := date(2026, 2, 8)

	var response string

	isInExamPeriod := false
	var firstStudyPeriodStart time.Time
	isRegWeek := false

	lang := r.URL.Query().Get("lang")
	if lang != "hu" {
		lang = "en"
	}

	verbose := r.URL.Query().Has("verbose")
	countdown := r.URL.Query().Has("countdown")

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
			regWeekEnd := firstStudyPeriodStart.Add(7 * 24 * time.Hour)
			if currentDate.Before(regWeekEnd) {
				isRegWeek = true
			}
		}

		if isRegWeek {
			if !verbose {
				response = "0"
			} else {
				if lang == "hu" {
					response = "Regisztrációs hét"
				} else {
					response = "Registration week"
				}
			}
		} else if !firstStudyPeriodStart.IsZero() {
			weeksPassed := int(currentDate.Sub(firstStudyPeriodStart).Hours() / (24 * 7))
			if verbose {
				suffix := getSuffix(weeksPassed)
				if lang == "hu" {
					suffix = "."
				}
				response = fmt.Sprintf("%d%s", weeksPassed, suffix)
			} else {
				response = fmt.Sprintf("%d", weeksPassed)
			}
		} else {
			// Breaks
			response = "-2"

			if verbose {
				if lang == "hu" {
					if countdown {
						days := calculateDaysBetween(currentDate, studyPeriods[0].Start)
						response = fmt.Sprintf("Szünet (%d nap van hátra)", days)
					} else {
						response = "Szünet"
					}
				} else {
					if countdown {
						days := calculateDaysBetween(currentDate, studyPeriods[0].Start)
						response = fmt.Sprintf("Break (%d days left)", days)
					} else {
						response = "Break"
					}
				}
			} else {
				if countdown {
					days := calculateDaysBetween(currentDate, studyPeriods[0].Start)
					response = fmt.Sprintf("%d", days)
				}
			}
		}
	} else {
		// Exams
		response = "-1"

		if verbose {
			if lang == "hu" {
				if countdown {
					days := calculateDaysBetween(currentDate, studyPeriods[0].Start)
					response = fmt.Sprintf("Vizsgaidőszak - szünet (%d nap van hátra)", days)
				} else {
					response = "Vizsgaidőszak - szünet"
				}
			} else {
				if countdown {
					days := calculateDaysBetween(currentDate, studyPeriods[0].Start)
					response = fmt.Sprintf("Exams - break (%d days left)", days)
				} else {
					response = "Exams - break"
				}
			}
		} else {
			if countdown {
				days := calculateDaysBetween(currentDate, studyPeriods[0].Start)
				response = fmt.Sprintf("%d", days)
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
