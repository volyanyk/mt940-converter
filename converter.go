package mt940_converter

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/shopspring/decimal"
)

const (
	crlf                  = "\r\n"
	referenceNumber       = ":20:"
	relatedReference      = ":21:"
	accountIdentification = ":25:"
	statementNumber       = ":28C:"
)

type ReferenceNumber struct {
	Value string
}
type RelatedReference struct {
	Value string
}
type StatementNumber struct {
	Value string
}

type AccountIdentification struct {
	CountryIso string
	Iban       string
	Currency   string
}
type TransactionType string

const (
	DEBIT  TransactionType = "D"
	CREDIT                 = "C"
)

type OpeningBalance struct {
	Type     TransactionType
	Date     time.Time
	Currency string
	Amount   decimal.Decimal
}

func GetReferenceNumber(input string) (*ReferenceNumber, error) {

	if !strings.Contains(input, referenceNumber) {
		return nil, fmt.Errorf("no reference number tag found. Expected tag: %s", referenceNumber)
	}
	index := strings.Index(input, crlf)
	result := input[len(referenceNumber):index]
	if len(result) > 16 {
		return nil, fmt.Errorf("the reference number character size is bigger than 16. Size: %v", len(input))
	}
	return &ReferenceNumber{Value: result}, nil
}

func GetRelatedReference(input string) (*RelatedReference, error) {

	if !strings.Contains(input, relatedReference) {
		return nil, fmt.Errorf("no related reference tag found. Expected tag: %s", relatedReference)
	}
	index := strings.Index(input, crlf)
	result := input[len(relatedReference):index]
	if len(result) > 16 {
		return nil, fmt.Errorf("the related reference character size is bigger than 16. Size: %v", len(input))
	}
	return &RelatedReference{Value: result}, nil
}
func GetAccountIdentification(input string) (*AccountIdentification, error) {
	if !strings.Contains(input, accountIdentification) {
		return nil, fmt.Errorf("no account identification tag found. Expected tag: %s", accountIdentification)
	}

	index := strings.Index(input, crlf)
	iban := input[len(accountIdentification):index]
	if len(iban) == 0 {
		return nil, fmt.Errorf("the reference number is empty. Size: %v", len(input))
	}
	if len(iban) > 35 {
		return nil, fmt.Errorf("the reference number character size is bigger than 35. Size: %v", len(input))
	}
	currency := GetLastNChars(iban, 3)
	country := GetFirstNChars(iban, 2)
	if len(country) == 0 {
		return nil, fmt.Errorf("the reference number does not contain country ISO code")
	}
	if len(currency) != 0 {
		iban = iban[len(country) : len(iban)-3]
	} else {
		iban = iban[len(country):]
	}
	return &AccountIdentification{
		CountryIso: country,
		Iban:       iban,
		Currency:   currency,
	}, nil
}

func GetStatementNumber(input string) (*StatementNumber, error) {

	if !strings.Contains(input, statementNumber) {
		return nil, fmt.Errorf("no statement number tag found. Expected tag: %s", statementNumber)
	}
	index := strings.Index(input, crlf)
	result := input[len(statementNumber):index]
	if len(result) > 5 {
		return nil, fmt.Errorf("the statement number character size is bigger than 5. Size: %v", len(input))
	}
	return &StatementNumber{Value: result}, nil
}

func GetLastNChars(input string, number int) string {
	s, done := validateString(input[len(input)-number:])
	if done {
		return s
	}
	return ""
}

func GetFirstNChars(input string, number int) string {
	s, done := validateString(input[:number])
	if done {
		return s
	}
	return ""

}

func validateString(input string) (string, bool) {
	alphabetic := true
	for _, char := range input {
		if !unicode.IsLetter(char) {
			alphabetic = false
			break
		}
	}
	if alphabetic {
		return input, true
	}
	return "", false
}
