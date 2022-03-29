package zabinski_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	julian "webeph/julian"
	nutation "webeph/nutation"
	testutils "webeph/testutils"
	unit "webeph/unit"
	zabinski "webeph/zabinski"
)

// Tolerance is one arc minute, or 1/60th of a degree.
const tolerance = 1. / 60.

// Checks for errors, and checks the tolerance between two measures.
// Receives:
//	got: a unit.Angle representing a measure
//	expected: the expected measure
//	err: any existing errors encountered during calculation
//	t: the testing platform
// Returns:
//	nothing
// Side effects:
//	Uses the testing platform to report any errors or tolerance issues.
func reportAnyErrors(got unit.Angle, expected float64, err error, t *testing.T) {
	degs := got.Deg()
	if err != nil {
		t.Errorf(err.Error())
	}
	if testutils.CheckTolerance(degs, expected, tolerance) == false {
		t.Errorf("Longitude: expected %v, got %v", expected, got.Deg())
	}
}

// Tests zabinski.FindLongitude for a variety of functions.
// Receives:
//	t: the testing platform
// Returns:
//	nothing
// Side effects:
//	Uses the testing platform to report success/failure of a few tests.
// Notes:
//	Checks:
//	1. Does the function recognize invalid planet strings?
//	2. Does the function calculate a sample of planet measures correctly?
func TestFindLongitude(t *testing.T) {
	_, err := zabinski.FindLongitude(1, 2, 3., 0., 0., 0., "nibiru", false)
	if err == nil {
		t.Errorf("Nibiru should be detected as imaginary planet")
	}
	// Test data is for January 19, 2022@1523 EST, Woonsocket, RI
	// Year, month as ints
	y, m := 2022, 1
	// Time as fractional day. Be sure to take time zone into account. (Below is for Greenwich Mean Time.)
	dy := 19.849056
	// Latitude
	φ := unit.AngleFromDeg(42.)
	// Longitude
	ο := unit.AngleFromDeg(-71.516667)
	// Height for Woonsocket, RI
	h := 56.0832

	testCases := []struct {
		planet   string
		expected float64
	}{
		{"saturn", 314.0401},
		{"jupiter", 334.4736},
		{"mars", 266.607},
		{"venus", 282.945},
		{"mercury", 307.5887},
		// Make sure that a repeated calculation works.
		{"saturn", 314.0401},
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("Longitude for %s", tc.planet),
			func(t *testing.T) {
				got, err := zabinski.FindLongitude(y, m, dy, φ, ο, h, tc.planet, false)
				reportAnyErrors(got, tc.expected, err, t)
			},
		)
	}

}

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

// Manage the maximum deviation.
// Receives:
//	currMax: the current maximum deviation
//	vsop87: the current VSOP87 measure
//	schlyter: the current Schlyter measure
// Returns:
//	the new maximum, whether the max received or a new max
func manageMaxDeviation(currMax, vsop87, schlyter float64) float64 {
	var diff float64
	if vsop87 > schlyter {
		diff = vsop87 - schlyter
	} else {
		diff = schlyter - vsop87
	}
	if diff > currMax {
		return diff
	}
	return currMax
}

// Tranforms a float64 to a string representation in degrees, minutes and seconds.
// Receives:
//	val: a float64
// Returns:
//	a string representing the float as a DMS. For example, 4.5 becomes '4° 30' 0"'
func floatToDMS(val float64) string {
	degs := math.Floor(val)
	mins := (val - math.Floor(val)) * 60
	secs := (mins - math.Floor(mins)) * 60
	sep := " "
	return fmt.Sprint(int(degs)) + "°" + sep + fmt.Sprint(int(mins)) + "'" + sep + fmt.Sprint(secs) + "\""
}

// Executive function comparing VSOP87 and Schlyter results for multiple planets.
// Receives:
//	t: the testing platform
// Returns:
//	nothing
// Side effects:
//	Runs testing functions for all appropriate planets.
func TestCompareVSOP87Schlyter(t *testing.T) {
	// compareVSOP87SchlyterPlanet("mercury", t)
	// compareVSOP87SchlyterPlanet("venus", t)
	compareVSOP87SchlyterPlanet("mars", t)
}

func displayDate(y, m int, fracDay float64) string {
	d := int(math.Floor(fracDay))
	remD := fracDay - float64(d)

	flH := remD * 24
	flM := (flH - math.Floor(flH)) * 60
	flS := (flM - math.Floor(flM)) * 60
	hr := int(flH)
	mins := int(flM)
	secs := int(flS)

	dt := time.Date(y, time.Month(m), d, hr, mins, secs, 0, time.UTC)
	return dt.String()

}

// Compares VSOP87 and Schlyter results for one planet.
// Receives:
//	planet: the planet being checked
//	t: the testing platform
// Returns:
//	nothing
// Side effects:
// Runs testing functions for one planet.
func compareVSOP87SchlyterPlanet(planet string, t *testing.T) {
	rando := rand.New(rand.NewSource(time.Now().UnixNano()))
	max := 0.

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
				// Get the full VSOP87 solution for the planet.
				vsop87, err := zabinski.FindLongitude(y, m, dy, φ, ο, h, planet, false)
				if err != nil {
					t.Errorf(err.Error())
				}
				// Get the Schlyter solution for the planet.
				schlyter, err := zabinski.FindLongitude(y, m, dy, φ, ο, h, planet, true)
				if err != nil {
					t.Errorf(err.Error())
				}
				vsop87D := vsop87.Deg()
				schlyterD := schlyter.Deg()
				// Maintain the max deviation across the samples.
				max = manageMaxDeviation(max, vsop87D, schlyterD)
				// Check whether the measures are within tolerance.
				if testutils.CheckTolerance(vsop87D, schlyterD, tolerance) == false {
					diff := 0.
					if vsop87D > schlyterD {
						diff = vsop87D - schlyterD
					} else {
						diff = schlyterD - vsop87D
					}
					t.Errorf(
						"VSOP87: %v, Schlyter: %v, diff: %v, dt: %v, latitude: %v, longitude: %v, height: %v",
						vsop87D,
						schlyterD,
						floatToDMS(diff),
						displayDate(y, m, dy),
						φ, ο, h)
				}
			}
		}
	}
	t.Logf("Max deviation: %v", floatToDMS(max))
}

func TestFindLST1(t *testing.T) {
	// Williams, example 1, p8
	ο := unit.AngleFromDeg(-74.283333)
	jd := julian.CalendarGregorianToJD(1951, 3, 5.677083)
	Δψ, Δε := nutation.Nutation(jd)
	expect := 79696.
	got := zabinski.FindSiderealTime(Δψ, Δε, jd, ο)
	if !testutils.CheckTolerance(got.Sec(), expect, 1.) {
		t.Errorf("TestFindLST: expected %v to be %v", got, expect)
	}
}

func TestFindLST2(t *testing.T) {
	// Williams, example 2, p8
	ο := unit.AngleFromDeg(-76.141667)
	jd := julian.CalendarGregorianToJD(1979, 5, 20.826389)
	Δψ, Δε := nutation.Nutation(jd)
	expect := 23814.5076
	got := zabinski.FindSiderealTime(Δψ, Δε, jd, ο)
	if !testutils.CheckTolerance(got.Sec(), expect, 1.) {
		t.Errorf("TestFindLST: expected %v to be %v", got, expect)
	}
}
