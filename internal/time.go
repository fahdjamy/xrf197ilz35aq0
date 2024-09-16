package internal

import "time"

func AddMonths(t time.Time, monthsToAdd int) time.Time {
	// Calculate the new month by adding the monthsToAdd to the current month
	newMonth := int(t.Month()) + monthsToAdd

	// newMonth - 1: We subtract 1 from newMonth because months in Go's time package are 1-indexed
	// (January is 1, February is 2, and so on). By subtracting 1, we essentially convert it to a 0-indexed value,
	// which simplifies the subsequent calculation.
	zeroIndexedMonth := newMonth - 1

	// Calculate the new year by adding the number of years that have passed due to the month addition
	newYear := t.Year() + zeroIndexedMonth/12

	// Adjust the newMonth to be within the range of 1-12
	// ensures that the newMonth remains within the valid range of 1 to 12
	newMonth = zeroIndexedMonth%12 + 1

	// Get the maximum number of days in the new month and year
	maxDaysInMonth := daysInMonth(newYear, time.Month(newMonth))

	// Ensure the day doesn't exceed the maximum days in the new month
	newDay := t.Day()
	if newDay > maxDaysInMonth {
		newDay = maxDaysInMonth
	}

	// Create and return a new time.Time object with the adjusted date and time
	hour := t.Hour()
	sec := t.Second()
	minute := t.Minute()
	location := t.Location()
	nanoSec := t.Nanosecond()
	return time.Date(newYear, time.Month(newMonth), newDay, hour, minute, sec, nanoSec, location)
}

// Helper function to get the number of days in a month
func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func IsAfterMonths(date time.Duration, month int) bool {
	now := time.Now()
	dateToCheck := now.Add(date)
	monthsFromNow := now.AddDate(0, month, 0)
	return dateToCheck.After(monthsFromNow)
}
