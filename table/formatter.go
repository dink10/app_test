package table

import (
	"strconv"
	"strings"
)

// CurrencyFormat formats value as currency.
func CurrencyFormat(val float64, precision int) []byte {
	asString := strconv.FormatFloat(val, 'f', precision, 64)
	sep := strings.Split(asString, ".")
	beforeDecimal := sep[0]
	withSeparator := make([]byte, 0, len(asString)+(len(beforeDecimal)/3))

	initial := len(beforeDecimal) % 3
	if initial > 0 {
		withSeparator = append(withSeparator, beforeDecimal[0:initial]...)
		beforeDecimal = beforeDecimal[initial:]
		if len(beforeDecimal) >= 3 {
			withSeparator = append(withSeparator, byte(' '))
		}
	}

	for len(beforeDecimal) >= 3 {
		withSeparator = append(withSeparator, beforeDecimal[0:3]...)
		beforeDecimal = beforeDecimal[3:]
		if len(beforeDecimal) >= 3 {
			withSeparator = append(withSeparator, byte(' '))
		}
	}
	if precision > 0 {
		withSeparator = append(withSeparator, byte(','))
		withSeparator = append(withSeparator, sep[1]...)
	}

	return withSeparator
}
