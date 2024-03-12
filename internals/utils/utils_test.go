package utils_test

import (
	"multitenant-api-go/internals/utils"
	"testing"
	"time"
)

func Test_GenerateApiKey(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name   string
		prefix string
		input  int8
		want   int
	}{
		// we want to receive a key that is prefix lenght + 1 (it adds char _) + input length
		{name: "generate valid key", prefix: "sb", input: 8, want: 2 + 8 + 1},
		{name: "generate valid key", prefix: "test", input: 14, want: 4 + 14 + 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := utils.GenerateApiKey(tt.input, tt.prefix)
			if err != nil {
				t.Fatalf("got error %v", err)
			}
			if len(result) != tt.want {
				t.Errorf("got %d, want %d", len(result), tt.want)
			}
		})
	}
}

func Test_ToValidDateString(t *testing.T) {
	t.Parallel()
	var tests = []struct {
		name  string
		input time.Time
		want  string
	}{
		{name: "zero date", input: time.Time{}, want: ""},
		{name: "valid date", input: time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), want: "2006-01-02T15:04:05Z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.ToValidDateString(tt.input)
			if result != tt.want {
				t.Errorf("got %s, want %s", result, tt.want)
			}
		})
	}
}
