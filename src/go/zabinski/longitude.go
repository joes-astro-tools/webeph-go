package zabinski

import (
	"errors"

	elliptic "webeph/elliptic"
	julian "webeph/julian"
	nutation "webeph/nutation"
	parallax "webeph/parallax"
	pp "webeph/planetposition"
	sidereal "webeph/sidereal"
	unit "webeph/unit"
)

// Map of full planet names to abbreviations. Abbreviations used below for more mapping.
var Planets = map[string]string{
	"saturn":  "sa",
	"jupiter": "ju",
	"mars":    "ma",
	"sun":     "su",
	"venus":   "ve",
	"mercury": "me",
	"moon":    "mo",
	"earth":   "ea"}

// Map of abbreviations to the ints used by planetposition to track planets.
var PPConsts = map[string]int{
	"sa": pp.Saturn,
	"ju": pp.Jupiter,
	"ma": pp.Mars,
	"ve": pp.Venus,
	"me": pp.Mercury,
	"ea": pp.Earth}

// Map of ints to planet data. These should be instantiated at runtime once.
var PPData = map[int]*pp.V87Planet{}

// Gets data for planetary positioning.
// Receives:
//	planet: a planet as a string. The fully spelled out name, ie "saturn" for Saturn, "venus" for Venus, etc.
// Returns:
//	planetData: data for the planet
//	err: any errors encountered
// Notes:
//	Also maintains the underlying map of data. If we haven't obtained the data, it loads the appropriate file.
//	Makes it so that all planetary data is only loaded into memory once.
func getPPData(pl string) (planetData *pp.V87Planet, err error) {
	plNum := PPConsts[pl]
	if PPData[plNum] == nil {
		PPData[PPConsts[pl]], err = pp.LoadPlanet(plNum)
		if err != nil {
			return nil, err
		}
	}
	return PPData[plNum], nil
}

// Finds obliquity.
// Receives:
//	jd: a Julian day
// Returns:
//	the obliquity, as a unit.Angle
func FindObliquity(Δε unit.Angle, jd float64) unit.Angle {
	ε0 := nutation.MeanObliquity(jd)
	return ε0 + Δε
}

// Finds sidereal time.
// Receives:
//	jd: a Julian day
//	ο: a geographic longitude, as an Angle
// Returns:
//	the sidereal time, as a unit.Time
func FindSiderealTime(Δψ, Δε unit.Angle, jd float64, ο unit.Angle) unit.Time {
	gmt := sidereal.Apparent(Δψ, Δε, jd)
	return (gmt + ο.Time()).Mod1()
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
func FindLongitude(y, m int, t float64, φ, ο unit.Angle, h float64, planet string, test bool) (λ unit.Angle, err error) {
	pl := Planets[planet]
	if pl == "" {
		return 0., errors.New("Invalid planet requested: get planet from zabinski.Planets")
	}
	plData, err := getPPData(pl)
	if err != nil {
		return 0., err
	}
	pl = Planets["earth"]
	earth, err := getPPData(pl)
	if err != nil {
		return 0., err
	}
	jd := julian.CalendarGregorianToJD(y, m, t)
	Δψ, Δε := nutation.Nutation(jd)
	ε := FindObliquity(Δε, jd)
	lst := FindSiderealTime(Δψ, Δε, jd, ο)
	geocentricλ, geocentricβ, geocentricΔ := elliptic.EclipticPosition(plData, earth, jd, test, Δψ)
	horizPx := parallax.Horizontal(geocentricΔ)
	λ = parallax.TopocentricLongitude(geocentricλ, geocentricβ, φ, h, ε, lst, horizPx)
	return
}
