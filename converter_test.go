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
		{name: "Statement number is correct", input: ":28C:44444\r\n", expectedResult: &StatementNumber{Value: "44444"}, hasError: false},
		{name: "Statement number is empty", input: ":28C:\r\n", expectedResult: &StatementNumber{Value: ""}, hasError: false},
		{name: "Statement number is too long", input: ":28C:555555\r\n", expectedResult: nil, hasError: true},
		{name: "Statement number tag not found", input: ":28:01234\r\n", expectedResult: nil, hasError: true},
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

func TestOpeningBalanceCase(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult *Balance
		hasError       bool
	}

	decim1, _ := GetDecimal("73447,91")
	decim2, _ := GetDecimal("734488877,91")
	testTable := []testCase{
		{name: "Opening balance is correct", input: ":60F:C120216UAH73447,91\r\n", expectedResult: &Balance{
			TransactionType: CREDIT,
			Date: InternalDate{
				Year:  12,
				Month: 2,
				Day:   16,
			},
			Currency:    "UAH",
			Amount:      decim1,
			BalanceType: OPENING,
		}, hasError: false},
		{name: "Opening balance is correct", input: ":60F:D110122PLN734488877,91\r\n", expectedResult: &Balance{
			TransactionType: DEBIT,
			Date: InternalDate{
				Year:  11,
				Month: 1,
				Day:   22,
			},
			Currency:    "PLN",
			Amount:      decim2,
			BalanceType: OPENING,
		}, hasError: false},
		{name: "Opening balance is empty", input: ":60F:\r\n", expectedResult: nil, hasError: true},
		{name: "Opening balance is too short", input: ":60F:C\r\n", expectedResult: nil, hasError: true},
		{name: "Opening balance is too long", input: ":60F:C120216UAH73447,9wwww\r\n", expectedResult: nil, hasError: true},
		{name: "Opening balance tag not found", input: ":60:01234\r\n", expectedResult: nil, hasError: true},
	}

	for _, test := range testTable {
		actual, err := GetBalance(test.input, OPENING)
		assert.Equal(t, test.expectedResult, actual, test.name)

		if test.hasError {
			assert.NotNil(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
}

func TestClosingBalanceCase(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult *Balance
		hasError       bool
	}

	decim1, _ := GetDecimal("73447,91")
	decim2, _ := GetDecimal("734488877,91")
	testTable := []testCase{
		{name: "Closing balance is correct", input: ":62F:C120216UAH73447,91\r\n", expectedResult: &Balance{
			TransactionType: CREDIT,
			Date: InternalDate{
				Year:  12,
				Month: 2,
				Day:   16,
			},
			Currency:    "UAH",
			Amount:      decim1,
			BalanceType: CLOSING,
		}, hasError: false},
		{name: "Closing balance is correct", input: ":62F:D110122PLN734488877,91\r\n", expectedResult: &Balance{
			TransactionType: DEBIT,
			Date: InternalDate{
				Year:  11,
				Month: 1,
				Day:   22,
			},
			Currency:    "PLN",
			Amount:      decim2,
			BalanceType: CLOSING,
		}, hasError: false},
		{name: "Closing balance is empty", input: ":62F:\r\n", expectedResult: nil, hasError: true},
		{name: "Closing balance is too short", input: ":62F:C\r\n", expectedResult: nil, hasError: true},
		{name: "Closing balance is too long", input: ":62F:C120216UAH73447,9wwww\r\n", expectedResult: nil, hasError: true},
		{name: "Closing balance tag not found", input: ":62:01234\r\n", expectedResult: nil, hasError: true},
	}

	for _, test := range testTable {
		actual, err := GetBalance(test.input, CLOSING)
		assert.Equal(t, test.expectedResult, actual, test.name)

		if test.hasError {
			assert.NotNil(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
}

func TestAvailableBalanceCase(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedResult *Balance
		hasError       bool
	}

	decim1, _ := GetDecimal("73447,91")
	decim2, _ := GetDecimal("734488877,91")
	testTable := []testCase{
		{name: "Closing balance is correct", input: ":64:C120216UAH73447,91\r\n", expectedResult: &Balance{
			TransactionType: CREDIT,
			Date: InternalDate{
				Year:  12,
				Month: 2,
				Day:   16,
			},
			Currency:    "UAH",
			Amount:      decim1,
			BalanceType: AVAILABLE,
		}, hasError: false},
		{name: "Available balance is correct", input: ":64:D110122PLN734488877,91\r\n", expectedResult: &Balance{
			TransactionType: DEBIT,
			Date: InternalDate{
				Year:  11,
				Month: 1,
				Day:   22,
			},
			Currency:    "PLN",
			Amount:      decim2,
			BalanceType: AVAILABLE,
		}, hasError: false},
		{name: "Available balance is empty", input: ":64:\r\n", expectedResult: nil, hasError: true},
		{name: "Available balance is too short", input: ":64:C\r\n", expectedResult: nil, hasError: true},
		{name: "Available balance is too long", input: ":64:C120216UAH73447,9wwww\r\n", expectedResult: nil, hasError: true},
		{name: "Available balance tag not found", input: ":62F:01234\r\n", expectedResult: nil, hasError: true},
	}

	for _, test := range testTable {
		actual, err := GetBalance(test.input, AVAILABLE)
		assert.Equal(t, test.expectedResult, actual, test.name)

		if test.hasError {
			assert.NotNil(t, err, test.name)
		} else {
			assert.Nil(t, err, test.name)
		}
	}
}
