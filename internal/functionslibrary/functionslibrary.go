package functionslibrary

import (
	"math/rand"
	"time"
)

// генерация рандомного флоата64 в заданых границах
func GenerateRandomValue(min int, max int, precision int) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	intPart := rng.Intn(max-min+1) + min

	decimalPart := float64(int(rand.Float64()*float64((precision*10)))) / float64((precision * 10))
	mixedValue := float64(intPart) + decimalPart

	if mixedValue <= float64(min) || mixedValue >= float64(max) {
		return float64(intPart)
	}

	return mixedValue
}
