package web

import (
	base "webeph/base"
	elliptic "webeph/elliptic"
	julian "webeph/julian"
	moonposition "webeph/moonposition"
	parallax "webeph/parallax"
	pp "webeph/planetposition"
	solar "webeph/solar"
	unit "webeph/unit"
	zabinski "webeph/zabinski"
)

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
func FindLongitude(y, m int, t float64, φ, ο unit.Angle, h float64, planet int) (λ unit.Angle, err error) {
	earth := LoadPlanet(pp.Earth)
	jd := julian.CalendarGregorianToJD(y, m, t)
	// Nutation is expensive: it more than doubles the calculation time.
	// Based on tests in seekNutation, it only improves accuracy by around 0.001 arcseconds.
	// No need to calculate it.
	// The nutation folder was left in place in case it is needed someday.
	Δψ := unit.Angle(0.)
	Δε := unit.Angle(0.)
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
		geocentricλ, geocentricβ, geocentricΔ = MoonPosition(jd)
		plx = moonposition.Parallax(geocentricΔ)
	default:
		plData = LoadPlanet(planet)
		geocentricλ, geocentricβ, geocentricΔ = elliptic.EclipticPosition(plData, earth, jd, false, Δψ)
		plx = parallax.Horizontal(geocentricΔ)
	}
	λ = parallax.TopocentricLongitude(geocentricλ, geocentricβ, φ, h, ε, lst, plx)
	return
}
