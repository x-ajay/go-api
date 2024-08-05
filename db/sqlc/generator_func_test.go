package db

import (
	"math/rand"
)

func generateTestOwner() string {
	owner := ""
	for i := 0; i < 6; i++ {
		random := rand.Intn(len(alpha))
		owner = owner + string(alpha[random])
	}
	return owner
}

func generateTestBalance() int32 {
	return rand.Int31()
}

func generateTestCurrency() Currency {
	n := rand.Intn(len(currency))
	return Currency(currency[n])
}
