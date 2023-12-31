package mt940_converter

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

const (
	crlf                   = "\r\n"
	referenceNumber        = ":20:"
	relatedReference       = ":21:"
	accountIdentification  = ":25:"
	statementNumber        = ":28C:"
	openingBalance         = ":60F:"
	closingBalance         = ":62F:"
	availableBalance       = ":64:"
	transaction            = ":61:"
	transactionDescription = ":86:"
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
type LongDate struct {
	Year  int64
	Month int64
	Day   int64
}
type ShortDate struct {
	Month int64
	Day   int64
}
type TransactionType string
type BalanceType string

const (
	DEBIT  TransactionType = "D"
	CREDIT                 = "C"
)
const (
	OPENING   BalanceType = "O"
	CLOSING               = "C"
	AVAILABLE             = "A"
)

type MyDecimal decimal.Decimal

type Balance struct {
	TransactionType TransactionType
	Date            LongDate
	Currency        string
	Amount          MyDecimal
	BalanceType     BalanceType
}
type TransactionStatement struct {
	LongDate               LongDate
	ShortDate              ShortDate
	TransactionType        TransactionType
	ThirdCurrencyCharacter string
	Amount                 MyDecimal
	Description            string
	DescriptionPrefix      string
}
type TransactionInformation struct {
	Info string
}
type Transaction struct {
	Index       int
	Statement   TransactionStatement
	Information TransactionInformation
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

func GetBalance(input string, balanceType BalanceType) (*Balance, error) {
	var tag string
	if balanceType == OPENING {
		tag = openingBalance
	}
	if balanceType == CLOSING {
		tag = closingBalance
	}
	if balanceType == AVAILABLE {
		tag = availableBalance
	}
	if tag == "" {
		return nil, fmt.Errorf("incorrect tag: %v", tag)

	}

	if !strings.Contains(input, tag) {
		return nil, fmt.Errorf("no proper tag found. Expected tag: %s", tag)
	}
	index := strings.Index(input, crlf)
	result := input[len(tag):index]
	if len(result) > 25 || len(result) < 10 {
		return nil, fmt.Errorf("the balance character size is incorrect. Size: %v", len(input))
	}
	amount, err := GetDecimal(result[10:])
	if err != nil {
		return nil, fmt.Errorf("cannot parse amount. Error: %v", err)
	}
	date, err := GetLongDate(result[1:7])
	if err != nil {
		return nil, fmt.Errorf("cannot parse date. Error: %v", err)
	}

	return &Balance{
		TransactionType: GetTransactionType(GetFirstNChars(result, 1)),
		Date:            *date,
		Currency:        result[7:10],
		Amount:          amount,
		BalanceType:     balanceType,
	}, nil
}

func GetTransactions(input string) (*[]Transaction, error) {
	transactionStrings := strings.Split(input, transaction)[1:]
	var transactions []Transaction
	var err error
	for i, transactionString := range transactionStrings {
		statement, err := GetStatement(transactionString)
		info := GetTransactionInfo(transactionString)
		if err != nil {
			break
		}

		transactions = append(transactions, Transaction{
			Index:       i + 1,
			Statement:   *statement,
			Information: info,
		})
	}
	if err != nil {
		return &[]Transaction{}, err
	}
	return &transactions, nil
}

func GetTransactionInfo(transactionString string) TransactionInformation {
	var info = transactionString[strings.LastIndex(transactionString, transactionDescription)+len(transactionDescription):]
	log.Printf(info)
	return TransactionInformation{Info: info}
}

func GetStatement(transactionString string) (*TransactionStatement, error) {
	var stmt = transactionString[:strings.Index(transactionString, transactionDescription)]
	valueLongDate, err := GetLongDate(stmt[:6])
	valueShortDate, err := GetShortDate(stmt[6:10])
	var transactionType = GetTransactionType(stmt[10:11])
	regex := regexp.MustCompile("^([A-Za-z])?(\\d{1,12},\\d{2}|\\d{1,3},\\d{3},\\d{2}|\\d{1,15})([A-Za-z])(.*?)$")
	matches := regex.FindStringSubmatch(regexp.MustCompile(`\r?\n`).ReplaceAllString(stmt[11:], " "))
	if err != nil {
		return nil, err
	}
	if matches != nil {
		thirdCurrencyCharacter := matches[1]
		amount, err := GetDecimal(matches[2])
		descriptionPrefix := matches[3]
		descr := matches[4]
		if err != nil {
			return nil, err
		}

		return &TransactionStatement{
			LongDate:               *valueLongDate,
			ShortDate:              *valueShortDate,
			TransactionType:        transactionType,
			ThirdCurrencyCharacter: thirdCurrencyCharacter,
			Amount:                 amount,
			DescriptionPrefix:      descriptionPrefix,
			Description:            descr,
		}, nil
	}
	return nil, errors.New("the input statement string is incorrect")

}
