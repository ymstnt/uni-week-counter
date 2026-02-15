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
	{Start: date(2026, 5, 26), End: date(2026, 7, 4)},
	{Start: date(2025, 12, 15), End: date(2026, 2, 6)},
	{Start: date(2025, 5, 26), End: date(2025, 7, 5)},
	{Start: date(2024, 12, 16), End: date(2025, 2, 8)},
	{Start: date(2024, 5, 20), End: date(2024, 6, 29)},
	{Start: date(2023, 12, 18), End: date(2024, 2, 3)},
}

// Start should be the first Monday of study periods (including registration week), End should be the last Saturday
// Periods should be descending, newest period should be at the top
var studyPeriods = []Period{
	{Start: date(2026, 2, 9), End: date(2026, 5, 23)},
	{Start: date(2025, 9, 1), End: date(2025, 12, 13)},
	{Start: date(2025, 2, 10), End: date(2025, 5, 24)},
	{Start: date(2024, 9, 2), End: date(2024, 12, 14)},
	{Start: date(2024, 2, 5), End: date(2024, 5, 18)},
	{Start: date(2023, 9, 4), End: date(2023, 12, 16)},
}

func isDateInPeriod(date time.Time, period Period) bool {
	periodEnd := period.End.Add(24 * time.Hour) // Include full end day
	return !date.Before(period.Start) && date.Before(periodEnd)
}

func calculateDaysBetween(start, end time.Time) int {
	days := int(end.Sub(start).Hours()/24) + 1
	return days
}

type ResponseData struct {
	Week            int      `json:"week"`
	Suffix          string   `json:"suffix"`
	Verbose         string   `json:"verbose"`
	Exam            bool     `json:"exam"`
	Study           bool     `json:"study"`
	RegWeek         bool     `json:"regWeek"`
	StudyPeriods    []Period `json:"studyPeriods"`
	ExamPeriods     []Period `json:"examPeriods"`
}

func getCurrentWeek(w http.ResponseWriter, r *http.Request) {
	currentDate := time.Now()
	//currentDate := date(2026, 2, 8)

	var weekNum int
	var verbose string
	var suffix string

	isInExamPeriod := false
	isInStudyPeriod := false
	isRegWeek := false
	var firstStudyPeriodStart time.Time

	lang := r.URL.Query().Get("lang")
	if lang != "hu" {
		lang = "en"
	}

	for _, period := range examPeriods {
		if isDateInPeriod(currentDate, period) {
			isInExamPeriod = true
			break
		}
	}

	if !isInExamPeriod {
		for _, period := range studyPeriods {
			if isDateInPeriod(currentDate, period) {
				isInStudyPeriod = true
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
			weekNum = 0
			
			if lang == "hu" {
				
			}
			if lang == "hu" {
				verbose = "Regisztrációs hét"
				suffix = "."
			} else {
				verbose = "Registration week"
				suffix = "th"
			}
		} else if !firstStudyPeriodStart.IsZero() {
			weeksPassed := int(currentDate.Sub(firstStudyPeriodStart).Hours() / (24 * 7))
			suffix = getSuffix(weeksPassed)
			if lang == "hu" {
				suffix = "."
			}
			weekNum = weeksPassed
		} else {
			// Breaks
			weekNum = calculateDaysBetween(currentDate, studyPeriods[0].Start)
			if lang == "hu" {
				verbose = "Szünet"
			} else {
				verbose = "Break"
			}
		}
	} else {
		// Exams
		weekNum = calculateDaysBetween(currentDate, studyPeriods[0].Start)

		if lang == "hu" {
			verbose = "Vizsgaidőszak - szünet"
		} else {
			verbose = "Exams - break"
		}
	}

	w.Header().Set("Content-Type", "application/json")

	response := ResponseData{
		Week:                 weekNum,
		Suffix:                suffix,
		Verbose:              verbose,
		Exam:          isInExamPeriod,
		Study:        isInStudyPeriod,
		RegWeek:            isRegWeek,
		StudyPeriods:    studyPeriods,
		ExamPeriods:      examPeriods,
	}

	json.NewEncoder(w).Encode(response)
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

func (p *Period) SetSemester() {
	year := p.Start.Year()
	var x int

	if p.Start.Month() < time.July {
		year--
		x = 2
	} else {
		x = 1
	}

	yy := (year + 1) % 100
	p.Semester = fmt.Sprintf("%d/%02d/%d", year, yy, x)
}

func main() {
	port := "8080"
	fmt.Println("Starting...")

	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	} else if len(os.Args) > 1 {
		port = os.Args[1]
	}

	for i := range studyPeriods {
		studyPeriods[i].SetSemester()
		examPeriods[i].SetSemester()
	}

	http.HandleFunc("/uwc", getCurrentWeek)

	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()
	fmt.Println("Server is listening on: ", port)

	select {}
}
