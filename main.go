package main

import (
	_ "embed"
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

// Define periods in periods.json
// For study periods, start should be the first Monday of study periods (including registration week), end should be the last Saturday
// For exam periods, start should be the first Monday of exam periods, end should be the last Saturday
// Periods should be descending, newest period should be at the top

//go:embed periods.json
var configBytes []byte

type ConfigPeriod struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type ConfigData struct {
	StudyPeriods []ConfigPeriod `json:"studyPeriods"`
	ExamPeriods  []ConfigPeriod `json:"examPeriods"`
}

var studyPeriods []Period
var examPeriods []Period

func loadPeriods() {
	var conf ConfigData
	json.Unmarshal(configBytes, &conf)

	for _, p := range conf.StudyPeriods {
		s, _ := time.Parse("2006-01-02", p.Start)
		e, _ := time.Parse("2006-01-02", p.End)
		studyPeriods = append(studyPeriods, Period{Start: s, End: e})
	}

	for _, p := range conf.ExamPeriods {
		s, _ := time.Parse("2006-01-02", p.Start)
		e, _ := time.Parse("2006-01-02", p.End)
		examPeriods = append(examPeriods, Period{Start: s, End: e})
	}
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
	Week         *int     `json:"week"`
	Suffix       string   `json:"suffix"`
	Verbose      string   `json:"verbose"`
	DaysLeft     int      `json:"daysLeft"`
	Exam         bool     `json:"exam"`
	Study        bool     `json:"study"`
	RegWeek      bool     `json:"regWeek"`
	StudyPeriods []Period `json:"studyPeriods"`
	ExamPeriods  []Period `json:"examPeriods"`
}

func getEverything(w http.ResponseWriter, r *http.Request) {
	currentDate := time.Now()
	//currentDate := date(2026, 2, 8)

	var weekNum *int = nil
	var daysLeft int
	var verbose string
	var suffix string

	isInExamPeriod := false
	isInStudyPeriod := false
	isRegWeek := false
	var currentStudyPeriodStart time.Time
	var currentStudyPeriodEnd time.Time

	countdownToEndOfBreak := r.URL.Query().Has("countdown-breaks")

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
				currentStudyPeriodStart = period.Start
				currentStudyPeriodEnd = period.End
				break
			}
		}

		if !currentStudyPeriodStart.IsZero() {
			regWeekEnd := currentStudyPeriodStart.Add(7 * 24 * time.Hour)
			if currentDate.Before(regWeekEnd) {
				isRegWeek = true
			}
		}

		if isRegWeek {
			val := 0
			weekNum = &val

			daysLeft = calculateDaysBetween(currentDate, currentStudyPeriodEnd)

			if lang == "hu" {
				verbose = "Regisztrációs hét"
				suffix = "."
			} else {
				verbose = "Registration week"
				suffix = "th"
			}
		} else if !currentStudyPeriodStart.IsZero() {
			weeksPassed := int(currentDate.Sub(currentStudyPeriodStart).Hours() / (24 * 7))
			daysLeft = calculateDaysBetween(currentDate, currentStudyPeriodEnd)
			suffix = getSuffix(weeksPassed)
			if lang == "hu" {
				suffix = "."
			}
			weekNum = &weeksPassed
		} else {
			// Breaks
			daysLeft = calculateDaysBetween(currentDate, studyPeriods[0].Start)
			if lang == "hu" {
				verbose = "Szünet"
			} else {
				verbose = "Break"
			}
		}
	} else {
		// Exams
		if countdownToEndOfBreak {
			daysLeft = calculateDaysBetween(currentDate, studyPeriods[0].Start)
		} else {
			daysLeft = calculateDaysBetween(currentDate, examPeriods[1].End)
		}

		if lang == "hu" {
			verbose = "Vizsgaidőszak - szünet"
		} else {
			verbose = "Exams - break"
		}
	}

	w.Header().Set("Content-Type", "application/json")

	response := ResponseData{
		Week:         weekNum,
		Suffix:       suffix,
		Verbose:      verbose,
		DaysLeft:     daysLeft,
		Exam:         isInExamPeriod,
		Study:        isInStudyPeriod,
		RegWeek:      isRegWeek,
		StudyPeriods: studyPeriods,
		ExamPeriods:  examPeriods,
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

	loadPeriods()

	for i := range studyPeriods {
		studyPeriods[i].SetSemester()
		examPeriods[i].SetSemester()
	}

	http.HandleFunc("/uwc", getEverything)

	go func() {
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()
	fmt.Println("Server is listening on: ", port)

	select {}
}
