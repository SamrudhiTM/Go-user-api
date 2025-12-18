package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name     string
		dob      time.Time
		expected int
	}{
		{
			name:     "Birthday already passed this year",
			dob:      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: time.Now().Year() - 2000,
		},
		{
			name:     "Birthday yet to come this year",
			dob:      time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
			expected: time.Now().Year() - 2000 - 1,
		},
		{
			name:     "Birthday today",
			dob:      time.Date(time.Now().Year()-20, time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC),
			expected: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAge(tt.dob)
			if got != tt.expected {
				t.Errorf("CalculateAge() = %d; want %d", got, tt.expected)
			}
		})
	}
}
