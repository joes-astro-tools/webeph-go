package main

import (
	unit "webeph/unit"
	web "webeph/web"
)

func main() {}

//export findLongitude
func FindLongitude(y, m int, t float64, φ, ο unit.Angle, h float64, plNum int) float64 {
	λ, err := web.FindLongitude(y, m, t, φ, ο, h, plNum)
	if err != nil {
		web.ErrMsg = err.Error()
		return 0
	}
	return λ.Deg()
}
