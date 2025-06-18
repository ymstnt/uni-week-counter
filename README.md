# University Week Counter API

This API provides information about study and exam periods for Obuda University. It allows you to check the current week in the academic calendar and retrieve the defined study and exam periods.

## Endpoints
1. Get Current Week

- Endpoint: `/uniWeekCount`
- Method: `GET`
#### Query Parameters:
| Parameter     | Type    | Required? | Description                                                                                             |
|---------------|---------|-----------|---------------------------------------------------------------------------------------------------------|
| lang          | string  | no        | Language for the response. Accepts either `en` (default) or `hu` (Hungarian).                           |
| numberOnly    | boolean | no        | If set to `true`, returns only the week number. Breaks and exams are still returned as text.            |
| daysLeftBreak | boolean | no        | If set to `true`, also returns the number of days left until the next study period during summer break. |
| daysLeftExams | boolean | no        | If set to `true`, also returns the number of days left until the next study period during exam periods. |

#### Example:
`GET https://uwc.ymstnt.com/uniWeekCount?lang=hu&numberOnly=true&daysLeftBreak=true&daysLeftExams=true`

Response:
```JSON
{
  "message": "1st" // or "1", "Exams - break", "Break", "Exams - break (69 days left)", etc.
}
```

2. Get Study Periods

- Endpoint: `/studyPeriods`
- Method: `GET`

#### Example:
`GET https://uwc.ymstnt.com/studyPeriods`

Response:
Returns a JSON array of study periods.
```JSON
[
  {
    "start": "2025-09-08T00:00:00Z",
    "end": "2025-12-13T00:00:00Z"
  },
  ...
]
```

3. Get Exam Periods

- Endpoint: `/examPeriods`
- Method: `GET`

#### Example:
`GET https://uwc.ymstnt.com/examPeriods`

Response:
Returns a JSON array of exam periods.
```JSON
[
  {
    "start": "2025-12-15T00:00:00Z",
    "end": "2026-02-06T00:00:00Z"
  },
  ...
]
```

## Usage
### Hosted
There is a hosted instance available at https://uwc.ymstnt.com

### Running Locally
- Clone the Repository:
  ```bash
  git clone <repository-url>
  cd <repository-directory>
  ```
- Install Go: Make sure you have Go installed.
- Run the API: `go run main.go`
    - Optionally, you can specify a port using a parameter: `go run main.go 9090`
    - ... or using an environment variable:
      ```bash
        export PORT=9090
        go run main.go
      ```
- Access the API: Open your browser or use a tool like Bruno or curl to access the endpoints:
  - Current Week: `http://localhost:8080/uniWeekCount`
  - Study Periods: `http://localhost:8080/studyPeriods`
  - Exam Periods: `http://localhost:8080/examPeriods`
