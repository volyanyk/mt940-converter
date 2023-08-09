package mt940_converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLongDateCase(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult *LongDate
		hasError       bool
	}

	testTable := []testCase{
		{name: "Long date is correct", input: "020222", expectedResult: &LongDate{
			Year:  2,
			Month: 2,
			Day:   22,
		}, hasError: false},
		{name: "Long date is empty", input: "032211", expectedResult: &LongDate{
			Year:  3,
			Month: 22,
			Day:   11,
		}, hasError: false},
		{name: "Long date is too long", input: "02010522222", expectedResult: nil, hasError: true},
		{name: "Long date is too short", input: "1111", expectedResult: nil, hasError: true},
	}

	for _, test := range testTable {
		actual, err := GetLongDate(test.input)
		assert.Equal(t, test.expectedResult, actual, test.name)

		if test.hasError {
			assert.NotNil(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}

}
func TestGetShortDateCase(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult *ShortDate
		hasError       bool
	}

	testTable := []testCase{
		{name: "Short date is correct", input: "0222", expectedResult: &ShortDate{
			Month: 2,
			Day:   22,
		}, hasError: false},
		{name: "Short date is empty", input: "2211", expectedResult: &ShortDate{
			Month: 22,
			Day:   11,
		}, hasError: false},
		{name: "Short date is too long", input: "02010522222", expectedResult: nil, hasError: true},
		{name: "Short date is too short", input: "111", expectedResult: nil, hasError: true},
	}

	for _, test := range testTable {
		actual, err := GetShortDate(test.input)
		assert.Equal(t, test.expectedResult, actual, test.name)

		if test.hasError {
			assert.NotNil(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}

}
