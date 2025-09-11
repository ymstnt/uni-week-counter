# University Week Counter API

This API provides information about study and exam periods for Obuda University. It allows you to check the current week in the academic calendar and retrieve the defined study and exam periods.

## Endpoints
### 1. Get Current Week

- Endpoint: `/uwc`
- Method: `GET`
#### Query Parameters:
| Parameter       | Type    | Required? | Description                                                                                                                                         |
| --------------- | ------- | --------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| lang            | string  | no        | Language for the response. Accepts either `en` (default) or `hu` (Hungarian). It will have no effect if number-only is present.                     |
| append-week     | boolean | no        | Append "week" in the appropriate language to the response. It will have no effect if number-only is preset.                                         |
| number-only     | boolean | no        | If present, returns only the week number. Exam periods are returned as `-1` and breaks are returned as `-2`                                         |
| days-left-break | boolean | no        | If present, also returns the number of days left until the next study period during summer break. It will have no effect if number-only is present. |
| days-left-exam  | boolean | no        | If present, also returns the number of days left until the next study period during exam periods. It will have no effect if number-only is present. |

#### Example:
`GET https://uwc.ymstnt.com/uwc?lang=hu&days-left-break&days-left-exams`

Response:
```JSON
{
  "message": "1." // or "1st", "1", "Exams - break", "Break", "Exams - break (69 days left)", etc.
}
```

### 2. Get Study Periods

- Endpoint: `/study-periods`
- Method: `GET`

#### Example:
`GET https://uwc.ymstnt.com/study-periods`

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

### 3. Get Exam Periods

- Endpoint: `/exam-periods`
- Method: `GET`

#### Example:
`GET https://uwc.ymstnt.com/exam-periods`

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
