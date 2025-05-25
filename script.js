// months are indexed from 0 to 11, 0 being January
let examPeriods = [
  // make sure to format dates correctly
  // exam dates should end on a Saturday
  { start: new Date(2023, 11, 18), end: new Date(2024, 1, 3) },
  { start: new Date(2024, 4, 20), end: new Date(2024, 5, 29) },
  { start: new Date(2024, 11, 16), end: new Date(2025, 1, 8) },
  { start: new Date(2025, 4, 26), end: new Date(2025, 6, 5) },
];

let studyPeriods = [
  // study dates should end on a Saturday
  { start: new Date(2023, 8, 11), end: new Date(2023, 11, 16) },
  { start: new Date(2024, 1, 12), end: new Date(2024, 4, 18) },
  { start: new Date(2024, 8, 9), end: new Date(2024, 11, 14) },
  { start: new Date(2025, 1, 17), end: new Date(2025, 4, 24) },
];

// this must be global
let weekNumberElement;
const currentDate = new Date();
const lang = document.documentElement.lang;

// Helper to check if a date is within a period, including the full end day
function isDateInPeriod(date, period) {
  const periodEnd = new Date(period.end);
  periodEnd.setHours(23, 59, 59, 999);
  return date >= period.start && date <= periodEnd;
}

function checkDateBefore() {
  weekNumberElement = document.querySelector("#week-number");

  let isInExamPeriod = examPeriods.some(
    (period) => isDateInPeriod(currentDate, period)
  );
  let isInStudyPeriod = studyPeriods.find(
    (period) => isDateInPeriod(currentDate, period)
  );

  if (isInExamPeriod) {
    if (lang === "hu") {
      weekNumberElement.innerText = "Vizsgaidőszak - szünet";
    } else {
      weekNumberElement.innerText = "Exams - break";
    }
  } else if (isInStudyPeriod) {
    calculateWeekNumber(isInStudyPeriod.start);
  } else {
    if (lang === "hu") {
      weekNumberElement.innerText = "Szünet";  
    } else {
      weekNumberElement.innerText = "Break";
    }
  }
}

function getFirstStudyPeriodDay() {
  let firstDay = new Date();

  for (const period of studyPeriods) {
    if (period.end > currentDate) {
      firstDay = period.start;
    }
  }

  return firstDay;
}

function calculateWeekNumber(firstWeek) {
  let oneWeekInMilliseconds = 7 * 24 * 60 * 60 * 1000;

  let timeDifference = currentDate.getTime() - firstWeek.getTime();
  let weeksPassed = Math.floor(timeDifference / oneWeekInMilliseconds) + 1;

  let suffix = ".";
  if (lang === "en") {
    switch (weeksPassed) {
      case 1:
        suffix = "st";
        break;
      case 2:
        suffix = "nd";
        break;
      case 3:
        suffix = "rd";
        break;
      default:
        suffix = "th";
    }
  }

  weekNumberElement.innerText = weeksPassed.toString() + suffix;
}

document.addEventListener("DOMContentLoaded", checkDateBefore);
