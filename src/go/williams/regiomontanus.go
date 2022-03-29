package williams

import (
	"math"

	unit "webeph/unit"
)

// Finds Regiomontanus houses for a northern latitude (above the equator).
// Receives:
//	ε: obliquity, as an Angle
//	oahc: oblique ascension of the house cusp, as a unit.Angle
//	qhc: quasi-latitude of the house circle, as a unit.Angle
// Returns:
//	the geocentric (and topocentric) ecliptic longitude of the house cusp, as a unit.Angle
func FindNorthernHouse(ε, oahc, qhc unit.Angle) unit.Angle {
	if q1(oahc) {
		a := unit.Angle(math.Acos((ε.Cos() * qhc.Sin()) + (ε.Sin() * qhc.Cos() * oahc.Cos())))
		coahc := unit.Angle(math.Acos(ε.Tan() * qhc.Tan()))
		if oahc.LessThan(coahc) {
			return unit.Angle(math.Asin((qhc.Cos() * oahc.Sin()) / a.Sin()))
		}
		if oahc.GreaterThan(coahc) {
			return angle180.Subtract(unit.Angle(math.Asin((qhc.Cos() * oahc.Sin()) / a.Sin())))
		}
		return angle90
	}
	if q2(oahc) {
		a := unit.Angle(math.Acos((-ε.Cos() * qhc.Sin()) - (ε.Sin() * qhc.Cos() * oahc.Cos())))
		return angle90.Add(unit.Angle(math.Acos((qhc.Cos() * oahc.Subtract(angle90).Cos()) / a.Sin())))
	}
	if q3(oahc) {
		a := unit.Angle(math.Acos((-ε.Cos() * qhc.Sin()) - (ε.Sin() * qhc.Cos() * oahc.Cos())))
		return angle180.Add(unit.Angle(math.Asin((qhc.Cos() * oahc.Subtract(angle180).Sin()) / a.Sin())))
	}
	a := unit.Angle(math.Acos((ε.Cos() * qhc.Sin()) + (ε.Sin() * qhc.Cos() * oahc.Cos())))
	coahc := angle360.Subtract(unit.Angle(math.Acos(ε.Tan() * qhc.Tan())))
	if oahc.LessThan(coahc) {
		return angle270.Subtract(unit.Angle(math.Acos((qhc.Cos() * oahc.Subtract(angle270).Cos()) / a.Sin())))
	}
	if oahc.GreaterThan(coahc) {
		return angle270.Add(unit.Angle(math.Acos((qhc.Cos() * oahc.Subtract(angle270).Cos()) / a.Sin())))
	}
	return angle270
}

// Finds Regiomontanus houses for a southern latitude (below the equator).
// Receives:
//	ε: obliquity, as an Angle
//	oahc: oblique ascension of the house cusp, as a unit.Angle
//	qhc: quasi-latitude of the house circle, as a unit.Angle
// Returns:
//	the geocentric (and topocentric) ecliptic longitude of the house cusp, as a unit.Angle
func FindSouthernHouse(ε, oahc, qhc unit.Angle) unit.Angle {
	if q1(oahc) {
		a := unit.Angle(math.Acos((-ε.Cos() * qhc.Sin()) + (ε.Sin() * qhc.Cos() * oahc.Cos())))
		return unit.Angle(math.Asin((qhc.Cos() * oahc.Sin()) / a.Sin()))
	}
	if q2(oahc) {
		a := unit.Angle(math.Acos((ε.Cos() * qhc.Sin()) - (ε.Sin() * qhc.Cos() * oahc.Cos())))
		coahc := unit.Angle(math.Acos(-ε.Tan() * qhc.Tan()))
		if oahc.LessThan(coahc) {
			return angle90.Subtract(unit.Angle(math.Acos((qhc.Cos() * oahc.Subtract(angle90).Cos()) / a.Sin())))
		}
		if oahc.GreaterThan(coahc) {
			return angle90.Add(unit.Angle(math.Acos((qhc.Cos() * oahc.Subtract(angle90).Cos()) / a.Sin())))
		}
		return angle90
	}
	if q3(oahc) {
		a := unit.Angle(math.Acos((ε.Cos() * qhc.Sin()) - (ε.Sin() * qhc.Cos() * oahc.Cos())))
		coahc := angle360.Subtract(unit.Angle(math.Acos(-ε.Tan() * qhc.Tan())))
		if oahc.LessThan(coahc) {
			return angle180.Add(unit.Angle(math.Asin((qhc.Cos() * oahc.Subtract(angle180).Sin()) / a.Sin())))
		}
		if oahc.GreaterThan(coahc) {
			return angle360.Subtract(unit.Angle(math.Asin((qhc.Cos() * oahc.Subtract(angle180).Sin()) / a.Sin())))
		}
		return angle270
	}
	a := unit.Angle(math.Acos((-ε.Cos() * qhc.Sin()) + (ε.Sin() * qhc.Cos() * oahc.Cos())))
	return angle270.Add(unit.Angle(math.Acos((qhc.Cos() * oahc.Subtract(angle270).Cos()) / a.Sin())))
}

// Finds Regiomontanus houses.
// Receives:
//	lst: local sidereal time, as a unit.Angle
//	ε: obliquity, as a unit.Angle
//	φ: latitude, as a unit.Angle
// Returns:
//	an array of 8 houses
// Notes:
//	Houses are returned in the following order: house 2, house 3, house 5, house 6, house 8, house 9, house 11, house 12
//	Ascendant, Midheaven and their opposites are calculated elsewhere.
func FindHouses(lst, ε, φ unit.Angle) [8]unit.Angle {
	lat := φ.Abs()
	geocentricLat := FindGeocentricLat(lat)
	if φ > angle90.Subtract(ε) {
		return FindArcticHouses(geocentricLat, lst, ε)
	}
	ramc := lst

	// Oblique ascensions of house cusps
	oa11 := ramc.Add(angle30).Mod1()
	oa12 := ramc.Add(angle60).Mod1()
	oa2 := ramc.Add(angle120).Mod1()
	oa3 := ramc.Add(angle150).Mod1()

	// Quasi-latitude of house circle: see p97 for a sketch of this concept.
	tanLat := geocentricLat.Tan()
	q11 := unit.Angle(math.Atan(tanLat * angle30.Sin()))
	q3 := q11
	q12 := unit.Angle(math.Atan(tanLat * angle60.Sin()))
	q2 := q12

	var (
		h2, h3, h11, h12 unit.Angle
	)

	if φ > 0 {
		h2 = FindNorthernHouse(ε, oa2, q2)
		h3 = FindNorthernHouse(ε, oa3, q3)
		h11 = FindNorthernHouse(ε, oa11, q11)
		h12 = FindNorthernHouse(ε, oa12, q12)

	} else {
		h2 = FindSouthernHouse(ε, oa2, q2)
		h3 = FindSouthernHouse(ε, oa3, q3)
		h11 = FindSouthernHouse(ε, oa11, q11)
		h12 = FindSouthernHouse(ε, oa12, q12)
	}

	h5 := h11.Add(angle180).Mod1()
	h6 := h12.Add(angle180).Mod1()
	h8 := h2.Add(angle180).Mod1()
	h9 := h3.Add(angle180).Mod1()

	return [8]unit.Angle{h2, h3, h5, h6, h8, h9, h11, h12}
}

// Finds Regiomonatanus house for arctic latitudes.
// Receives:
//	ε: obliquity, as a unit.Angle
//	oahc: oblique ascension of the house cusp, as a unit.Angle
//	qhc: quasi-latitude of the house circle, as a unit.Angle
// Returns:
//	the geocentric (and topocentric) ecliptic longitude of the house cusp, as a unit.Angle
func FindArcticHouse(ε, oahc, qhc unit.Angle) unit.Angle {
	q1oahc := func() bool { return oahc.GreaterThanOrEqual(unit.Angle(0)) && (oahc.LessThan(angle90)) }
	q2oahc := func() bool { return oahc.GreaterThanOrEqual(angle90) && oahc.LessThan(angle180) }
	q3oahc := func() bool { return oahc.GreaterThanOrEqual(angle180) && oahc.LessThan(angle270) }

	if q1oahc() {
		a := unit.Angle(math.Acos((-ε.Cos() * qhc.Sin()) + (ε.Sin() * qhc.Cos() * oahc.Cos())))
		return unit.Angle(math.Asin((qhc.Cos() * oahc.Sin()) / a.Sin()))
	}
	if q2oahc() {
		a := unit.Angle(math.Acos((ε.Cos() * qhc.Sin()) - (ε.Sin() * qhc.Cos() * oahc.Cos())))
		if qhc > angle90.Subtract(ε) {
			return angle90.Subtract(unit.Angle(math.Acos((qhc.Cos() * oahc.Subtract(angle90).Cos()) / a.Sin())))
		}
		coahc := unit.Angle(math.Acos(-ε.Tan() * qhc.Tan()))
		if oahc < coahc {
			return angle90.Subtract(unit.Angle(math.Acos((qhc.Cos() * oahc.Subtract(angle90).Cos()) / a.Sin())))
		}
		if oahc > coahc {
			return angle90.Add(unit.Angle(math.Acos((qhc.Cos() * oahc.Subtract(angle90).Cos()) / a.Sin())))
		}
		return angle90
	}
	if q3oahc() {
		a := unit.Angle(math.Acos((ε.Cos() * qhc.Sin()) - (ε.Sin() * qhc.Cos() * oahc.Cos())))
		if qhc > angle90.Subtract(ε) {
			return angle360.Subtract(unit.Angle(math.Asin((qhc.Cos() * oahc.Subtract(angle180).Sin()) / a.Sin())))
		}
		coahc := angle360.Subtract(unit.Angle(math.Acos(-ε.Tan() * qhc.Tan())))
		if oahc < coahc {
			return angle180.Add(unit.Angle(math.Asin((qhc.Cos() * oahc.Subtract(angle180).Sin()) / a.Sin())))
		}
		if oahc > coahc {
			return angle360.Subtract(unit.Angle(math.Asin((qhc.Cos() * oahc.Subtract(angle180).Sin()) / a.Sin())))
		}
		return angle270
	}
	a := unit.Angle(math.Acos((-ε.Cos() * qhc.Sin()) + (ε.Sin() * qhc.Cos() * oahc.Cos())))
	return angle270.Add(unit.Angle(math.Acos((qhc.Cos() * oahc.Subtract(angle270).Cos()) / a.Sin())))
}

// Finds Regiomontanus houses for arctic latitudes.
// Receives:
//	geocentricLat: geocentric latitude, as a unit.Angle
//	lst: local sidereal time, as a unit.Angle
//	ε: obliquity, as a unit.Angle
// Returns:
//	an array of 8 houses
// Notes:
//	Houses are returned in the following order: house 2, house 3, house 5, house 6, house 8, house 11, house 12
//	Ascendant, Midheaven and their opposites are calculated elsewhere.
func FindArcticHouses(geocentricLat, lst, ε unit.Angle) [8]unit.Angle {
	horizonDecl := angle90.Subtract(geocentricLat)
	// 'pced' = 'points of critical ecliptic declination'
	tan := unit.Angle(math.Asin(horizonDecl.Tan() / ε.Tan()))
	pced1 := angle180.Add(tan)
	pced2 := angle360.Subtract(tan)
	ramc := unit.Angle(lst.Rad())

	// Quasi-latitude of house circle: see p97 for a sketch of this concept.
	tanLat := geocentricLat.Tan()
	q11 := unit.Angle(math.Atan(tanLat * angle30.Sin()))
	q3 := q11
	q12 := unit.Angle(math.Atan(tanLat * angle60.Sin()))
	q2 := q12

	var (
		h2, h3, h11, h12 unit.Angle
	)
	lowerLatOkay := func() bool { return (pced1 > ramc) && (ramc > pced2) }
	if lowerLatOkay() {
		// Calculate as if we were in a lower latitude.

		// Oblique ascensions of house cusps
		oa11 := ramc.Add(angle30).Mod1()
		oa12 := ramc.Add(angle60).Mod1()
		oa2 := ramc.Add(angle120).Mod1()
		oa3 := ramc.Add(angle150).Mod1()

		h2 = FindNorthernHouse(ε, oa2, q2)
		h3 = FindNorthernHouse(ε, oa3, q3)
		h11 = FindNorthernHouse(ε, oa11, q11)
		h12 = FindNorthernHouse(ε, oa12, q12)

		h5 := h11.Add(angle180).Mod1()
		h6 := h12.Add(angle180).Mod1()
		h8 := h2.Add(angle180).Mod1()
		h9 := h3.Add(angle180).Mod1()

		return [8]unit.Angle{h2, h3, h5, h6, h8, h9, h11, h12}
	}

	// Oblique ascensions of house cusps
	oa11 := ramc.Subtract(angle30).Mod1()
	oa12 := ramc.Subtract(angle60).Mod1()
	oa2 := ramc.Subtract(angle120).Mod1()
	oa3 := ramc.Subtract(angle150).Mod1()

	h2 = FindArcticHouse(ε, oa2, q2)
	h3 = FindArcticHouse(ε, oa3, q3)
	h11 = FindArcticHouse(ε, oa11, q11)
	h12 = FindArcticHouse(ε, oa12, q12)

	h5 := h11.Add(angle180).Mod1()
	h6 := h12.Add(angle180).Mod1()
	h8 := h2.Add(angle180).Mod1()
	h9 := h3.Add(angle180).Mod1()

	return [8]unit.Angle{h2, h3, h5, h6, h8, h9, h11, h12}
}
