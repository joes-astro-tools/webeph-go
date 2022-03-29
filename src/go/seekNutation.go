package main

import (
	"math"
	"math/rand"
	"sync"
	"time"

	apparent "webeph/apparent"
	base "webeph/base"
	elliptic "webeph/elliptic"
	julian "webeph/julian"
	moonposition "webeph/moonposition"
	parallax "webeph/parallax"
	pp "webeph/planetposition"
	solar "webeph/solar"
	unit "webeph/unit"
	web "webeph/web"
	zabinski "webeph/zabinski"
)

// Finds the number of days in a month.
// Receives:
//	m: a month represented in Go's time package
//	year: an integer year
// Returns:
// The number of days in a month
func seekDaysInMonth(m time.Month, year int) int {
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
func seekRandomRange(min, max float64, rando *rand.Rand) float64 {
	// 1. Find a random integer between 0 and max
	// 2. Convert the integer to a float (add a .0)
	// 3. Add a random fraction portion
	// 4. Add the minimum.
	return (float64(rando.Int63n(int64(max))) + rando.Float64()) + min
}

// nutation returns nutation in longitude (Δψ) and nutation in obliquity (Δε)
// for a given JDE.
//
// JDE = UT + ΔT, see package deltat.
//
// Computation is by 1980 IAU theory, with terms < .0003″ neglected.
func nutation(jde float64, termsRemoved int) (Δψ, Δε unit.Angle) {
	T := base.J2000Century(jde)
	D := base.Horner(T,
		297.85036, 445267.11148, -0.0019142, 1./189474) * math.Pi / 180
	M := base.Horner(T,
		357.52772, 35999.050340, -0.0001603, -1./300000) * math.Pi / 180
	N := base.Horner(T,
		134.96298, 477198.867398, 0.0086972, 1./5620) * math.Pi / 180
	F := base.Horner(T,
		93.27191, 483202.017538, -0.0036825, 1./327270) * math.Pi / 180
	Ω := base.Horner(T,
		125.04452, -1934.136261, 0.0020708, 1./450000) * math.Pi / 180
	// sum in reverse order to accumulate smaller terms first
	var Δψs, Δεs float64
	for i := len(table22A) - termsRemoved; i >= 0; i-- {
		row := table22A[i]
		arg := row.d*D + row.m*M + row.n*N + row.f*F + row.ω*Ω
		s, c := math.Sincos(arg)
		Δψs += s * (row.s0 + row.s1*T)
		Δεs += c * (row.c0 + row.c1*T)
	}
	Δψ = unit.AngleFromSec(Δψs * .0001)
	Δε = unit.AngleFromSec(Δεs * .0001)
	return
}

func seekMoonPosition(y, m int, t float64, termsRemoved int) (λ unit.Angle) {
	jde := julian.CalendarGregorianToJD(y, m, t)
	geocentricλ, geocentricβ, _ := moonposition.Position(jde)
	Δψ, _ := nutation(jde, termsRemoved)
	Δλ, _ := apparent.EclipticAberration(geocentricλ, geocentricβ, jde)
	λ = (geocentricλ + Δλ + Δψ).Mod1()
	return
}

func q1(lon float64) bool {
	return lon >= 0. && lon < 90.
}

func q4(lon float64) bool {
	return lon >= 270. && lon < 360.
}

func findDiff(result, expected float64) float64 {
	if q1(expected) && q4(result) {
		return (360. - result) + expected
	}
	if q1(result) && q4(expected) {
		return (360. - expected) + result
	}
	if expected > result {
		return expected - result
	}
	return result - expected
}

// Finds topocentric longitude for a planet.
// Receives:
//	y: the year, as an int
//	m: the month, as an int
//	t: the time as a float64. Represents a fractional day.
//	φ: geographic latitude, as a unit.Angle
//	ο: geographic longitude, as a unit.Angle
//	h: the height above mean sea level, in meters
//	planet: the required planet, as a fully spelled out string, ie "saturn" for Saturn, etc.
// Returns:
//	λ: the topocentric ecliptic longitude, as a unit.Angle
//	err: any errors encountered
func seekFindLongitude(y, m int, t float64, φ, ο unit.Angle, h float64, planet, termsRemoved int) (λ unit.Angle) {
	earth := web.LoadPlanet(pp.Earth)
	jd := julian.CalendarGregorianToJD(y, m, t)
	Δψ, Δε := nutation(jd, termsRemoved)
	ε := zabinski.FindObliquity(Δε, jd)
	lst := zabinski.FindSiderealTime(Δψ, Δε, jd, ο)
	plData := &pp.V87Planet{}
	var geocentricλ, geocentricβ, plx unit.Angle
	var geocentricΔ float64
	switch planet {
	case pp.Sun:
		geocentricλ = solar.ApparentLongitude(base.J2000Century(jd))
		geocentricβ = 0.
		geocentricΔ = 1.
		plx = parallax.Horizontal(geocentricΔ)
	case pp.Venus:
		plData.Ibody = planet
		geocentricλ, geocentricβ, geocentricΔ = elliptic.EclipticPosition(plData, earth, jd, true, Δψ)
		plx = parallax.Horizontal(geocentricΔ)
	case pp.Mercury:
		plData.Ibody = planet
		geocentricλ, geocentricβ, geocentricΔ = elliptic.EclipticPosition(plData, earth, jd, true, Δψ)
		plx = parallax.Horizontal(geocentricΔ)
	case pp.Moon:
		geocentricλ, geocentricβ, geocentricΔ = web.MoonPosition(jd)
		plx = moonposition.Parallax(geocentricΔ)
	default:
		plData = web.LoadPlanet(planet)
		geocentricλ, geocentricβ, geocentricΔ = elliptic.EclipticPosition(plData, earth, jd, false, Δψ)
		plx = parallax.Horizontal(geocentricΔ)
	}
	λ = parallax.TopocentricLongitude(geocentricλ, geocentricβ, φ, h, ε, lst, plx)
	return
}

// Checks a month.
// Receives:
//	y: the current year
//	m: the current month
//	c: the outgoing channel to send results
//	wg: the outgoing wait group, waiting for results
// Returns:
//	nothing
// Side effects:
//	Send true or false to the outgoing channel
//	Notifies the wait group when finished
func checkMonth(y, m int, c chan float64, wg *sync.WaitGroup, termsRemoved int) {
	defer wg.Done()
	max := 0.
	rando := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Day, as an integer
	totalDays := seekDaysInMonth(time.Month(m), y)
	for d := 1; d < totalDays; d++ {
		// Time: fractional day out of 24 hours
		dy := float64(d) + rando.Float64()
		// Latitude: random number between -90 and 90
		φ := unit.Angle(seekRandomRange(-90., 180., rando))
		// Longitude: random number between -180 and 180
		ο := unit.Angle(seekRandomRange(-180., 360, rando))
		// Height: random number between -420 (Dead Sea) and 8848 (Mt Everest)
		h := seekRandomRange(-420., 8848, rando)
		// Get the solution for the planet using the full data set.
		// Saturn
		origSaturnλ := seekFindLongitude(y, m, dy, φ, ο, h, pp.Saturn, 1)
		// Jupiter
		origJupiterλ := seekFindLongitude(y, m, dy, φ, ο, h, pp.Jupiter, 1)
		// Mars
		origMarsλ := seekFindLongitude(y, m, dy, φ, ο, h, pp.Mars, 1)
		// Venus
		origVenusλ := seekFindLongitude(y, m, dy, φ, ο, h, pp.Venus, 1)
		// Mercury
		origMercuryλ := seekFindLongitude(y, m, dy, φ, ο, h, pp.Mercury, 1)
		// Moon
		origMoonλ := seekMoonPosition(y, m, dy, 1)
		// Saturn with trimmed terms
		trimmedSaturnλ := seekFindLongitude(y, m, dy, φ, ο, h, pp.Saturn, termsRemoved)
		// Jupiter with trimmed terms
		trimmedJupiterλ := seekFindLongitude(y, m, dy, φ, ο, h, pp.Jupiter, termsRemoved)
		// Mars with trimmed terms
		trimmedMarsλ := seekFindLongitude(y, m, dy, φ, ο, h, pp.Mars, termsRemoved)
		// Venus with trimmed terms
		trimmedVenusλ := seekFindLongitude(y, m, dy, φ, ο, h, pp.Venus, termsRemoved)
		// Mercury with trimmed terms
		trimmedMercuryλ := seekFindLongitude(y, m, dy, φ, ο, h, pp.Mercury, termsRemoved)
		// Moon with trimmed terms
		trimmedMoonλ := seekMoonPosition(y, m, dy, termsRemoved)
		diff := zabinski.FindDiff(origSaturnλ, trimmedSaturnλ)
		if diff > max {
			max = diff
		}
		diff = zabinski.FindDiff(origJupiterλ, trimmedJupiterλ)
		if diff > max {
			max = diff
		}
		diff = zabinski.FindDiff(origMarsλ, trimmedMarsλ)
		if diff > max {
			max = diff
		}
		diff = zabinski.FindDiff(origVenusλ, trimmedVenusλ)
		if diff > max {
			max = diff
		}
		diff = zabinski.FindDiff(origMercuryλ, trimmedMercuryλ)
		if diff > max {
			max = diff
		}
		diff = zabinski.FindDiff(origMoonλ, trimmedMoonλ)
		if diff > max {
			max = diff
		}
	}
	c <- max
}

// Spawns a series of goroutines to check every month in a year.
// Receives:
//	y: the year
//	c: the outgoing channel to send results
//	wg: the outgoing wait group, waiting for results
// Returns:
//	nothing
// Side effects:
//	Send true or false to the outgoing channel
//	Notifies the wait group when finished
func spawnMonths(y int, c chan float64, wg *sync.WaitGroup, termsRemoved int) {
	defer wg.Done()

	max := 0.
	localWG := new(sync.WaitGroup)
	localC := make(chan float64, 12)
	m := 1
	for m < 13 {
		localWG.Add(1)
		go checkMonth(y, m, localC, localWG, termsRemoved)
		m += 1
	}
	// Wait for all results to come in.
	localWG.Wait()
	// Shut down the channel collecting results.
	close(localC)
	// Read results.
	for deviation := range localC {
		if deviation > max {
			max = deviation
		}
	}
	c <- max
}

// Checks a set of times to find the maximum deviation from a variety of times and places.
// Receives:
//	nothing
// Returns:
//	a float64 representing the maximum deviation discovered
// Notes:
//	The configuration itself is stored in plDataCopy, and called deeper in the child functions.
//	Basically, this function starts a chain of asynchronous calls for all the times needed.
//	It eventually checks every day for the next 50 years.
func findMaxDeviation(termsRemoved int) float64 {
	var wg sync.WaitGroup
	c := make(chan float64, 50)
	max := 0.
	y := 2022
	for y < 2023 {
		wg.Add(1)
		// For every year, spawn functions to check each month.
		go spawnMonths(y, c, &wg, termsRemoved)
		y += 1
	}
	// Wait for the 600 spawned functions to resolve.
	wg.Wait()
	// Close off the channel.
	close(c)
	// Check whether any of the months were not within tolerance.
	for deviation := range c {
		if deviation > max {
			max = deviation
		}
	}
	return max
}

func main() {
	tolerance := 1. / 60.
	termsRemoved := 1
	deviation := 0.
	for (termsRemoved < len(table22A)) && (deviation < tolerance) {
		termsRemoved += 1
		deviation = findMaxDeviation(termsRemoved)
	}
	println(termsRemoved)
	println(deviation)
}

var table22A = []struct {
	d, m, n, f, ω  float64
	s0, s1, c0, c1 float64
}{
	{0, 0, 0, 0, 1, -171996, -174.2, 92025, 8.9},
	{-2, 0, 0, 2, 2, -13187, -1.6, 5736, -3.1},
	{0, 0, 0, 2, 2, -2274, -0.2, 977, -0.5},
	{0, 0, 0, 0, 2, 2062, 0.2, -895, 0.5},
	{0, 1, 0, 0, 0, 1426, -3.4, 54, -0.1},
	{0, 0, 1, 0, 0, 712, 0.1, -7, 0},
	{-2, 1, 0, 2, 2, -517, 1.2, 224, -0.6},
	{0, 0, 0, 2, 1, -386, -0.4, 200, 0},
	{0, 0, 1, 2, 2, -301, 0, 129, -0.1},
	{-2, -1, 0, 2, 2, 217, -0.5, -95, 0.3},
	{-2, 0, 1, 0, 0, -158, 0, 0, 0},
	{-2, 0, 0, 2, 1, 129, 0.1, -70, 0},
	{0, 0, -1, 2, 2, 123, 0, -53, 0},
	{2, 0, 0, 0, 0, 63, 0, 0, 0},
	{0, 0, 1, 0, 1, 63, 0.1, -33, 0},
	{2, 0, -1, 2, 2, -59, 0, 26, 0},
	{0, 0, -1, 0, 1, -58, -0.1, 32, 0},
	{0, 0, 1, 2, 1, -51, 0, 27, 0},
	{-2, 0, 2, 0, 0, 48, 0, 0, 0},
	{0, 0, -2, 2, 1, 46, 0, -24, 0},
	{2, 0, 0, 2, 2, -38, 0, 16, 0},
	{0, 0, 2, 2, 2, -31, 0, 13, 0},
	{0, 0, 2, 0, 0, 29, 0, 0, 0},
	{-2, 0, 1, 2, 2, 29, 0, -12, 0},
	{0, 0, 0, 2, 0, 26, 0, 0, 0},
	{-2, 0, 0, 2, 0, -22, 0, 0, 0},
	{0, 0, -1, 2, 1, 21, 0, -10, 0},
	{0, 2, 0, 0, 0, 17, -0.1, 0, 0},
	{2, 0, -1, 0, 1, 16, 0, -8, 0},
	{-2, 2, 0, 2, 2, -16, 0.1, 7, 0},
	{0, 1, 0, 0, 1, -15, 0, 9, 0},
	{-2, 0, 1, 0, 1, -13, 0, 7, 0},
	{0, -1, 0, 0, 1, -12, 0, 6, 0},
	{0, 0, 2, -2, 0, 11, 0, 0, 0},
	{2, 0, -1, 2, 1, -10, 0, 5, 0},
	{2, 0, 1, 2, 2, -8, 0, 3, 0},
	{0, 1, 0, 2, 2, 7, 0, -3, 0},
	{-2, 1, 1, 0, 0, -7, 0, 0, 0},
	{0, -1, 0, 2, 2, -7, 0, 3, 0},
	{2, 0, 0, 2, 1, -7, 0, 3, 0},
	{2, 0, 1, 0, 0, 6, 0, 0, 0},
	{-2, 0, 2, 2, 2, 6, 0, -3, 0},
	{-2, 0, 1, 2, 1, 6, 0, -3, 0},
	{2, 0, -2, 0, 1, -6, 0, 3, 0},
	{2, 0, 0, 0, 1, -6, 0, 3, 0},
	{0, -1, 1, 0, 0, 5, 0, 0, 0},
	{-2, -1, 0, 2, 1, -5, 0, 3, 0},
	{-2, 0, 0, 0, 1, -5, 0, 3, 0},
	{0, 0, 2, 2, 1, -5, 0, 3, 0},
	{-2, 0, 2, 0, 1, 4, 0, 0, 0},
	{-2, 1, 0, 2, 1, 4, 0, 0, 0},
	{0, 0, 1, -2, 0, 4, 0, 0, 0},
	{-1, 0, 1, 0, 0, -4, 0, 0, 0},
	{-2, 1, 0, 0, 0, -4, 0, 0, 0},
	{1, 0, 0, 0, 0, -4, 0, 0, 0},
	{0, 0, 1, 2, 0, 3, 0, 0, 0},
	{0, 0, -2, 2, 2, -3, 0, 0, 0},
	{-1, -1, 1, 0, 0, -3, 0, 0, 0},
	{0, 1, 1, 0, 0, -3, 0, 0, 0},
	{0, -1, 1, 2, 2, -3, 0, 0, 0},
	{2, -1, -1, 2, 2, -3, 0, 0, 0},
	{0, 0, 3, 2, 2, -3, 0, 0, 0},
	{2, -1, 0, 2, 2, -3, 0, 0, 0},
}
