package mt940_converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionReferenceCase(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult ReferenceNumber
		hasError       bool
	}

	testTable := []testCase{
		{name: "Reference number is correct", input: ":20:referenceNumber1\r\n", expectedResult: ReferenceNumber{Value: "referenceNumber1"}, hasError: false},
		{name: "Reference number is empty", input: ":20:\r\n", expectedResult: ReferenceNumber{Value: ""}, hasError: false},
		{name: "Reference number is too long", input: ":20:referenceNumber12\r\n", expectedResult: ReferenceNumber{Value: ""}, hasError: true},
		{name: "Reference tag not found", input: ":0:testReferenceNumber\r\n", expectedResult: ReferenceNumber{Value: ""}, hasError: true},
	}

	for _, test := range testTable {
		actual, err := GetReferenceNumber(test.input)
		assert.Equal(t, test.expectedResult, actual, test.name)

		if test.hasError {
			assert.NotNil(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
}
