package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	pp "webeph/planetposition"
	unit "webeph/unit"
	web "webeph/web"

	"github.com/pkg/profile"
)

// Finds the number of days in a month.
// Receives:
//	m: a month represented in Go's time package
//	year: an integer year
// Returns:
// The number of days in a month
func daysInMonth(m time.Month, year int) int {
	// This is equivalent to time.daysIn(m, year).
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// Finds a random number within a numerical range.
// Receives:
//	min: the minimum number in the range
//	max: the maximum number in the range
//	rando: an instance of Go's rand package
// Returns:
//	the selected random number
// Notes:
//	Make the minimum negative to drop the number to something that includes a negative range.
//	You will need a bit of thinking to use it this way.
//	For example, min: -10, max: 20 selects a number from -10 to 10.
func randomRange(min, max float64, rando *rand.Rand) float64 {
	// 1. Find a random integer between 0 and max
	// 2. Convert the integer to a float (add a .0)
	// 3. Add a random fraction portion
	// 4. Add the minimum.
	return (float64(rando.Int63n(int64(max))) + rando.Float64()) + min
}

func main() {
	rando := rand.New(rand.NewSource(time.Now().UnixNano()))
	profiler := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	for y := 2022; y < 2073; y++ {
		for m := 1; m < 13; m++ {
			// Day, as an integer
			totalDays := daysInMonth(time.Month(m), y)
			for d := 1; d < totalDays; d++ {
				// Time: fractional day out of 24 hours
				dy := float64(d) + rando.Float64()
				// Latitude: random number between -90 and 90
				φ := unit.Angle(randomRange(-90., 180., rando))
				// Longitude: random number between -180 and 180
				ο := unit.Angle(randomRange(-180., 360, rando))
				// Height: random number between -420 (Dead Sea) and 8848 (Mt Everest)
				h := randomRange(-420., 8848, rando)
				findLongitude(y, m, dy, φ, ο, h, pp.Mars)
			}
		}
	}
	profiler.Stop()
}

func findLongitude(y, m int, t float64, φ, ο unit.Angle, h float64, plNum int) (float64, error) {
	λ, err := web.FindLongitude(y, m, t, φ, ο, h, plNum)
	if err != nil {
		return 0., err
	}
	return λ.Deg(), nil
}

// Transforms a float64 to a string representation in degrees, minutes and seconds.
// Receives:
//	val: a float64
// Returns:
//	a string representing the float as a DMS. For example, 4.5 becomes '4° 30' 0"'
func floatToDMS(val float64) string {
	degs := math.Floor(val)
	mins := (val - math.Floor(val)) * 60
	secs := (mins - math.Floor(mins)) * 60
	secString := fmt.Sprintf("%.6f", secs)
	sep := " "
	return fmt.Sprint(int(degs)) + "°" + sep + fmt.Sprint(int(mins)) + "'" + sep + secString + "\""
}
