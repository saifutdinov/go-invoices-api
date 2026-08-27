package domain

import (
	"strconv"
	"strings"
)

// Important amount parsing function. Have to use for every amount
func ParseAmount(amount string) int64 {

	amountParts := strings.Split(amount, ".")

	switch len(amountParts) {
	case 1:
		amountInt, err := strconv.ParseInt(amountParts[0], 10, 64)
		if err != nil {
			return 0
		}
		// euro - cent multiplier
		return amountInt * 100
	case 2:
		amountInt, err := strconv.ParseInt(amountParts[0], 10, 64)
		if err != nil {
			return 0
		}

		// euro - cent multiplier
		amountInt *= 100

		amountCents, err := strconv.ParseInt(amountParts[1], 10, 64)
		if err != nil {
			return 0
		}

		if amountCents > 0 {
			amountInt += amountCents
		}

		return amountInt
	}
	return 0

}
