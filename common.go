package mt940_converter

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

func GetLongDate(s string) (*LongDate, error) {
	if len(s) != 6 {
		log.Error().Msg("Incorrect date length")
		return nil, errors.New("incorrect date length")
	}
	year, err := strconv.ParseInt(s[0:2], 10, 8)
	month, err := strconv.ParseInt(s[2:4], 10, 8)
	day, err := strconv.ParseInt(s[4:6], 10, 8)

	if err != nil {
		return nil, err
	}

	return &LongDate{
		Year:  year,
		Month: month,
		Day:   day,
	}, nil
}

func GetShortDate(s string) (*ShortDate, error) {
	if len(s) != 4 {
		return nil, errors.New("incorrect date length")
	}
	month, err := strconv.ParseInt(s[0:2], 10, 8)
	day, err := strconv.ParseInt(s[2:4], 10, 8)

	if err != nil {
		return nil, err
	}

	return &ShortDate{
		Month: month,
		Day:   day,
	}, nil
}

func GetDecimal(s string) (MyDecimal, error) {
	numberWithoutComma := strings.ReplaceAll(s, ",", "")

	decimalNumber, err := decimal.NewFromString(numberWithoutComma)
	if err != nil {
		return MyDecimal{}, err
	}

	return MyDecimal(decimalNumber.Div(decimal.NewFromInt(100))), nil
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

func GetTransactionType(result string) TransactionType {
	return TransactionType(result)
}
