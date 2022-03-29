package testutils

import (
	"fmt"
	"math"
)

const StandardTolerance = 1. / 60.     // One arc minute
const SecondTolerance = 1. / 3600.     // One arc second
const JulianMinuteTolerance = 0.000694 // One minute of time, measured in fractional Julian days

func CheckTolerance(result, expected, tolerance float64) bool {
	return math.Abs(result-expected) < tolerance
}

// Transforms a float64 to a string representation in degrees, minutes and seconds.
// Receives:
//	val: a float64
// Returns:
//	a string representing the float as a DMS. For example, 4.5 becomes '4° 30' 0"'
func FloatToDMS(val float64) string {
	degs := math.Floor(val)
	mins := (val - math.Floor(val)) * 60
	secs := (mins - math.Floor(mins)) * 60
	sep := " "
	return fmt.Sprint(int(degs)) + "°" + sep + fmt.Sprint(int(mins)) + "'" + sep + fmt.Sprint(secs) + "\""
}
