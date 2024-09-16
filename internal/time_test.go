package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAddMonths(t *testing.T) {
	// Test cases with various combinations of input dates and months to add
	testCases := []struct {
		name         string
		inputDate    time.Time
		monthsToAdd  int
		expectedDate time.Time
	}{
		{
			monthsToAdd:  3,
			name:         "Add months within the same year",
			inputDate:    time.Date(2024, 9, 15, 19, 49, 0, 0, time.UTC),
			expectedDate: time.Date(2024, 12, 15, 19, 49, 0, 0, time.UTC),
		},
		{
			monthsToAdd:  4,
			name:         "Add months crossing into the next year",
			inputDate:    time.Date(2024, 9, 15, 19, 49, 0, 0, time.UTC),
			expectedDate: time.Date(2025, 1, 15, 19, 49, 0, 0, time.UTC),
		},
		{
			monthsToAdd:  1,
			name:         "Add months with day overflow",
			inputDate:    time.Date(2024, 1, 31, 19, 49, 0, 0, time.UTC),
			expectedDate: time.Date(2024, 2, 29, 19, 49, 0, 0, time.UTC), // Leap year
		},
		{
			monthsToAdd:  1,
			name:         "Add months with day overflow (non-leap year)",
			inputDate:    time.Date(2023, 1, 31, 19, 49, 0, 0, time.UTC),
			expectedDate: time.Date(2023, 2, 28, 19, 49, 0, 0, time.UTC),
		},
		{
			monthsToAdd:  36, // 3 years
			name:         "Add a large number of months",
			inputDate:    time.Date(2024, 9, 15, 19, 49, 0, 0, time.UTC),
			expectedDate: time.Date(2027, 9, 15, 19, 49, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := AddMonths(tc.inputDate, tc.monthsToAdd)
			if !result.Equal(tc.expectedDate) {
				t.Errorf("For input date %v and months to add %d, expected %v but got %v", tc.inputDate, tc.monthsToAdd, tc.expectedDate, result)
			}
		})
	}
}

func TestIsAfterMonths(t *testing.T) {
	currentDate := time.Now()
	testCases := []struct {
		name        string
		afterMonths int
		expected    bool
		inputDate   time.Time
	}{
		{name: "returns true for input date after 2 months", afterMonths: 2, inputDate: AddMonths(currentDate, 3), expected: true},
		{name: "returns false for input date before 2 months", afterMonths: 2, inputDate: AddMonths(currentDate, 1), expected: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dateAfter := tc.afterMonths
			sub := tc.inputDate.Sub(currentDate)
			assert.Equal(t, tc.expected, IsAfterMonths(sub, dateAfter))
		})
	}
}
