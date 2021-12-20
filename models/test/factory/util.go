package factory

import (
	"fmt"
	"math"
	"math/rand"
)

func RandInt() int {
	seed := 1000000
	r := int(math.Floor(rand.Float64() * float64(seed)))
	return r
}

func RandIntString() string {
	return fmt.Sprintf("something%v", RandInt())
}
