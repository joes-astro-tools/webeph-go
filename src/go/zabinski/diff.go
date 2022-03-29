package zabinski

import (
	unit "webeph/unit"
)

// Finds whether a longitude is in quadrant 1.
// Receives:
//	λ: longitude, in degrees
// Returns:
//	True if the longitude is in quadrant 1. Otherwise false.
// Notes:
//	Quadrant 1 is anything greater than or equal to 0 degrees, and less than 90 degrees.
func q1(λ float64) bool {
	return λ >= 0. && λ < 90.
}

// Finds whether a longitude is in quadrant 4.
// Receives:
//	λ: longitude, in degrees
// Returns:
//	True if the longitude is in quadrant 4. Otherwise false.
// Notes:
//	Quadrant 4 is anything greater than or equal to 270 degrees, and less than 360 degrees.
func q4(λ float64) bool {
	return λ >= 270. && λ < 360.
}

// Finds the difference between two longitudes.
// Receives:
//	a: longitude, in degrees
//	b: longitude, in degrees
func FindDiff(a, b unit.Angle) float64 {
	aD := a.Deg()
	bD := b.Deg()
	if q1(aD) && q4(bD) {
		return (360. - bD) + aD
	}
	if q1(bD) && q4(aD) {
		return (360. - aD) + bD
	}
	if aD > bD {
		return aD - bD
	}
	return bD - aD
}
