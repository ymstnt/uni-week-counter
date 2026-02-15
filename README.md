# University Week Counter API

This API provides information about study and exam periods for Obuda University. It allows you to check the current week in the academic calendar and retrieve the defined study and exam periods.

## Endpoints

### 1. Get everything

- Endpoint: `/uwc`
- Method: `GET`

#### Query Parameters:

| Parameter | Type   | Required? | Description                                                                                                             |
| --------- | ------ | --------- | ----------------------------------------------------------------------------------------------------------------------- |
| lang      | string | no        | Language for the response. Affects the suffix and verbose responses. Accepts either `en` (default) or `hu` (Hungarian). |

#### Example:

`GET https://api.ymstnt.com/uwc?lang=hu`

Response:

Returns the number of the week (or remaining days if it's a break or exam period), suffix (language dependent), verbose name (if it's not a study period), if it's an exam period, if it's a study period, if it's registration week, the study periods and the exam periods.

```JSON
{
  "week": 1,
  "suffix": ".",
  "verbose": "", // can be: "Registration week" or "Regisztrációs hét", "Exams - break" or "Vizsgaidőszak - szünet" and "Break" or "Szünet"
  "exam": false,
  "study": true,
  "regWeek": false,
  "studyPeriods": [
    {
      "start": "2026-02-09T00:00:00Z",
      "end": "2026-05-23T00:00:00Z",
      "semester": "2025/26/2"
    },
    {
      "start": "2025-09-01T00:00:00Z",
      "end": "2025-12-13T00:00:00Z",
      "semester": "2025/26/1"
    },
    ...
  ],
  "examPeriods": [
    {
      "start": "2026-05-26T00:00:00Z",
      "end": "2026-07-04T00:00:00Z",
      "semester": "2025/26/2"
    },
    {
      "start": "2025-12-15T00:00:00Z",
      "end": "2026-02-06T00:00:00Z",
      "semester": "2025/26/1"
    },
    ...
  ]
}
```

| Key          | Type                                                 | Description                                                                                                                   |
| ------------ | ---------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| week         | int                                                  | The current numnber of the week of the semester OR the remaining days of the break/exam period.                               |
| suffix       | string                                               | Suffix of the number of the week, depending on the language.                                                                  |
| verbose      | string                                               | The verbose name of the current period. Only has effect if it's registration week/exam period/break. Depends on the language. |
| exam         | bool                                                 | True if an exam period is active. Mutually exclusive with `study`.                                                            |
| study        | bool                                                 | True if a study period is active. Mutually exclusive with `exam`.                                                             |
| regWeek      | bool                                                 | True if it's registration week. (the 0th week of a semester)                                                                  |
| studyPeriods | array[ { start(time), end(time), semester(string)} ] | A JSON array of study periods in descending order.                                                                            |
| examPeriods  | array[ { start(time), end(time), semester(string)}   | A JSON array of exam periods in descending order.                                                                             |

## Usage

### Hosted

There is a hosted instance available at https://api.ymstnt.com

### Running Locally

- Clone the Repository:
  ```bash
  git clone https://github.com/ymstnt/uni-week-counter
  cd uni-week-counter
  ```
- Install Go: Make sure you have Go installed. Alternatively, if you're using Nix, you can do `nix develop` instead.
- Run the API: `go run main.go`
  - Optionally, you can specify a port using a parameter: `go run main.go 9090`
  - ... or using an environment variable:
    ```bash
      export PORT=9090
      go run main.go
    ```
- Access the API: Open your browser or use a tool like Bruno or curl to access the endpoints:
  - Current Week: `http://localhost:8080/uwc`
  - Study Periods: `http://localhost:8080/study-periods`
  - Exam Periods: `http://localhost:8080/exam-periods`
