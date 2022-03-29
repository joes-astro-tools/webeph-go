// Apparent: Chapter 23, Apparent Place of a Star
package apparent

import (
	base "webeph/base"
	coord "webeph/coord"
	nutation "webeph/nutation"
	precess "webeph/precess"
	solar "webeph/solar"
	unit "webeph/unit"
)

// κ is the constnt of aberration in radians.
var κ = unit.AngleFromSec(20.49552)

// longitude of perihelian of Earth's orbit.
func perihelion(T float64) unit.Angle {
	return unit.AngleFromDeg(base.Horner(T, 102.93735, 1.71946, .00046))
}

// EclipticAberration returns corrections due to aberration for ecliptic
// coordinates of an object.
func EclipticAberration(λ, β unit.Angle, jd float64) (Δλ, Δβ unit.Angle) {
	T := base.J2000Century(jd)
	s, _ := solar.True(T)
	e := solar.Eccentricity(T)
	π := perihelion(T)
	sβ, cβ := β.Sincos()
	ssλ, csλ := (s - λ).Sincos()
	sπλ, cπλ := (π - λ).Sincos()
	// (23.2) p. 151
	Δλ = κ.Mul((e*cπλ - csλ) / cβ)
	Δβ = -κ.Mul(sβ * (ssλ - e*sπλ))
	return
}

// Aberration returns corrections due to aberration for equatorial
// coordinates of an object.
func Aberration(α unit.RA, δ unit.Angle, jd float64) (Δα2 unit.HourAngle, Δδ2 unit.Angle) {
	ε := nutation.MeanObliquity(jd)
	T := base.J2000Century(jd)
	s, _ := solar.True(T)
	e := solar.Eccentricity(T)
	π := perihelion(T)
	sα, cα := α.Sincos()
	sδ, cδ := δ.Sincos()
	ss, cs := s.Sincos()
	sπ, cπ := π.Sincos()
	cε := ε.Cos()
	tε := ε.Tan()
	q1 := cα * cε
	// (23.3) p. 152
	Δα2 = unit.HourAngle(κ.Rad() * (e*(q1*cπ+sα*sπ) - (q1*cs + sα*ss)) / cδ)
	q2 := cε * (tε*cδ - sα*sδ)
	q3 := cα * sδ
	Δδ2 = κ.Mul(e*(cπ*q2+sπ*q3) - (cs*q2 + ss*q3))
	return
}

// Nutation returns corrections due to nutation for equatorial coordinates
// of an object.
//
// Results are invalid for objects very near the celestial poles.
func Nutation(α unit.RA, δ unit.Angle, jd float64) (Δα1 unit.HourAngle, Δδ1 unit.Angle) {
	ε := nutation.MeanObliquity(jd)
	sε, cε := ε.Sincos()
	Δψ, Δε := nutation.Nutation(jd)
	sα, cα := α.Sincos()
	tδ := δ.Tan()
	// (23.1) p. 151
	Δα1 = unit.HourAngle((cε+sε*sα*tδ)*Δψ.Rad() - cα*tδ*Δε.Rad())
	Δδ1 = Δψ.Mul(sε*cα) + Δε.Mul(sα)
	return
}

// Position computes the apparent position of an object.
//
// Position is computed for equatorial coordinates in eqFrom, considering
// proper motion, precession, nutation, and aberration.  Result is in
// eqTo.  EqFrom and eqTo must be non-nil, but may point to the same struct.
func Position(eqFrom, eqTo *coord.Equatorial, epochFrom, epochTo float64, mα unit.HourAngle, mδ unit.Angle) *coord.Equatorial {
	precess.Position(eqFrom, eqTo, epochFrom, epochTo, mα, mδ)
	jd := base.JulianYearToJDE(epochTo)
	Δα1, Δδ1 := Nutation(eqTo.RA, eqTo.Dec, jd)
	Δα2, Δδ2 := Aberration(eqTo.RA, eqTo.Dec, jd)
	eqTo.RA = eqTo.RA.Add(Δα1 + Δα2)
	eqTo.Dec += Δδ1 + Δδ2
	return eqTo
}
