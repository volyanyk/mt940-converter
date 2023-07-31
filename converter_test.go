package mt940_converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionReferenceCase(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult *ReferenceNumber
		hasError       bool
	}

	testTable := []testCase{
		{name: "Reference number is correct", input: ":20:referenceNumber1\r\n", expectedResult: &ReferenceNumber{Value: "referenceNumber1"}, hasError: false},
		{name: "Reference number is empty", input: ":20:\r\n", expectedResult: &ReferenceNumber{Value: ""}, hasError: false},
		{name: "Reference number is too long", input: ":20:referenceNumber12\r\n", expectedResult: nil, hasError: true},
		{name: "Reference tag not found", input: ":0:testReferenceNumber\r\n", expectedResult: nil, hasError: true},
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

func TestAccountIdentificationCase(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult *AccountIdentification
		hasError       bool
	}

	testTable := []testCase{
		{name: "Account identification is correct", input: ":25:NL17RABO6064103256EUR\r\n", expectedResult: &AccountIdentification{
			CountryIso: "NL",
			Iban:       "17RABO6064103256",
			Currency:   "EUR",
		}, hasError: false},
		{name: "Account identification is correct, even without currency", input: ":25:NL17RABO6064103256\r\n", expectedResult: &AccountIdentification{
			CountryIso: "NL",
			Iban:       "17RABO6064103256",
			Currency:   "",
		}, hasError: false},
		{name: "Account identification is empty", input: ":25:\r\n", expectedResult: nil, hasError: true},
		{name: "Account identification not found", input: "NL17RABO6064103256EUR\r\n", expectedResult: nil, hasError: true},
		{name: "Account identification is too long", input: ":25:NI81CCSF6843126715474931687323111UAH\r\n", expectedResult: nil, hasError: true},
	}

	for _, test := range testTable {
		actual, err := GetAccountIdentification(test.input)
		assert.Equal(t, test.expectedResult, actual, test.name)

		if test.hasError {
			assert.NotNil(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
}

func TestRelatedReferenceCase(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult *RelatedReference
		hasError       bool
	}

	testTable := []testCase{
		{name: "Related reference is correct", input: ":21:relatedReference\r\n", expectedResult: &RelatedReference{Value: "relatedReference"}, hasError: false},
		{name: "Related reference is empty", input: ":21:\r\n", expectedResult: &RelatedReference{Value: ""}, hasError: false},
		{name: "Related reference is too long", input: ":21:relatedReference12\r\n", expectedResult: nil, hasError: true},
		{name: "Related reference tag not found", input: ":0:relatedReference\r\n", expectedResult: nil, hasError: true},
	}

	for _, test := range testTable {
		actual, err := GetRelatedReference(test.input)
		assert.Equal(t, test.expectedResult, actual, test.name)

		if test.hasError {
			assert.NotNil(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
}

func TestStatementNumberCase(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult *StatementNumber
		hasError       bool
	}

	testTable := []testCase{
		{name: "Statement number is correct", input: ":25:44444\r\n", expectedResult: &StatementNumber{Value: "44444"}, hasError: false},
		{name: "Statement number is empty", input: ":25:\r\n", expectedResult: &StatementNumber{Value: ""}, hasError: false},
		{name: "Statement number is too long", input: ":25:555555\r\n", expectedResult: nil, hasError: true},
		{name: "Statement number tag not found", input: ":0:01234\r\n", expectedResult: nil, hasError: true},
	}

	for _, test := range testTable {
		actual, err := GetStatementNumber(test.input)
		assert.Equal(t, test.expectedResult, actual, test.name)

		if test.hasError {
			assert.NotNil(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
}
