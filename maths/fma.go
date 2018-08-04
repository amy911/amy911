package maths

func Fma(base, mul, add float64) float64 {
	return base * mul + add
}

func Fma32(base, mul, add float32) float64 {
	return Fma(float64(base), float64(mul), float64(add))
}