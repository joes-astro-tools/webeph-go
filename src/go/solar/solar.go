// Solar: Chapter 25, Solar Coordinates.
//
// Partial implementation:
//
// 1. Higher accuracy positions are not computed with Appendix III but with
// full VSOP87 as implemented in package planetposition.
//
// 2. Higher accuracy correction for aberration (using the formula for
// variation Δλ on p. 168) is not implemented.  Results for example 25.b
// already match the full VSOP87 values on p. 165 even with the low accuracy
// correction for aberration, thus there are no more significant digits that
// would check a more accurate result.  Also the size of the formula presents
// significant chance of typographical error.
package solar

import (
	base "webeph/base"
	unit "webeph/unit"
)

// True returns true geometric longitude and anomaly of the sun referenced to the mean equinox of date.
//
// Argument T is the number of Julian centuries since J2000.
// See base.J2000Century.
//
// Results:
//	s = true geometric longitude, ☉
//	ν = true anomaly
func True(T float64) (s, ν unit.Angle) {
	// (25.2) p. 163
	L0 := unit.AngleFromDeg(base.Horner(T, 280.46646, 36000.76983, 0.0003032))
	M := MeanAnomaly(T)
	C := unit.AngleFromDeg(base.Horner(T, 1.914602, -0.004817, -.000014)*
		M.Sin() +
		(0.019993-.000101*T)*M.Mul(2).Sin() +
		0.000289*M.Mul(3).Sin())
	return (L0 + C).Mod1(), (M + C).Mod1()
}

// MeanAnomaly returns the mean anomaly of Earth at the given T.
//
// Argument T is the number of Julian centuries since J2000.
// See base.J2000Century.
//
// Result is not normalized to the range 0..2π.
func MeanAnomaly(T float64) unit.Angle {
	// (25.3) p. 163
	return unit.AngleFromDeg(base.Horner(T, 357.52911, 35999.05029, -0.0001537))
}

// Eccentricity returns eccentricity of the Earth's orbit around the sun.
//
// Argument T is the number of Julian centuries since J2000.
// See base.J2000Century.
func Eccentricity(T float64) float64 {
	// (25.4) p. 163
	return base.Horner(T, 0.016708634, -0.000042037, -0.0000001267)
}

// ApparentLongitude returns apparent longitude of the Sun referenced
// to the true equinox of date.
//
// Argument T is the number of Julian centuries since J2000.
// See base.J2000Century.
//
// Result includes correction for nutation and aberration.
func ApparentLongitude(T float64) unit.Angle {
	Ω := node(T)
	s, _ := True(T)
	return s - unit.AngleFromDeg(.00569) -
		unit.AngleFromDeg(.00478).Mul(Ω.Sin())
}

func node(T float64) unit.Angle {
	return unit.AngleFromDeg(125.04 - 1934.136*T)
}
