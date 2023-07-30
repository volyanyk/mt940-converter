package mt940_converter

import (
	"fmt"
	"strings"
)

const (
	crlf            = "\r\n"
	referenceNumber = ":20:"
)

type ReferenceNumber struct {
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
