package mt940_converter

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	crlf                  = "\r\n"
	referenceNumber       = ":20:"
	relatedReference      = ":21:"
	accountIdentification = ":25:"
)

type ReferenceNumber struct {
	Value string
}
type RelatedReference struct {
	Value string
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

type AccountIdentification struct {
	CountryIso string
	Iban       string
	Currency   string
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
