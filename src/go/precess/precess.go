// Precession: Chapter 21, Precession.

package precess

import (
	"math"

	base "webeph/base"
	coord "webeph/coord"
	nutation "webeph/nutation"
	unit "webeph/unit"
)

// Precessor represents precession from one epoch to another.
//
// Construct with NewPrecessor, then call method Precess.
// After construction, Precess may be called multiple times to precess
// different coordinates with the same initial and final epochs.
type Precessor struct {
	ζ      unit.RA
	z      unit.Angle
	sθ, cθ float64
}

const d = math.Pi / 180
const s = d / 3600

// Package variables allow these slices to be reused.  (As composite
// literals inside of NewPrecessor they would be reallocated on every
// function call.)
var (
	// coefficients from (21.2) p. 134
	ζT = []float64{2306.2181 * s, 1.39656 * s, -0.000139 * s}
	zT = []float64{2306.2181 * s, 1.39656 * s, -0.000139 * s}
	θT = []float64{2004.3109 * s, -0.8533 * s, -0.000217 * s}

	// coefficients from (21.3) p. 134
	ζt = []float64{2306.2181 * s, 0.30188 * s, 0.017998 * s}
	zt = []float64{2306.2181 * s, 1.09468 * s, 0.018203 * s}
	θt = []float64{2004.3109 * s, -0.42665 * s, -0.041833 * s}
)

// EclipticPrecessor represents precession from one epoch to another.
//
// Construct with NewEclipticPrecessor, then call method Precess.
// After construction, Precess may be called multiple times to precess
// different coordinates with the same initial and final epochs.
type EclipticPrecessor struct {
	sη, cη float64
	π, p   unit.Angle
}

var (
	// coefficients from (21.5) p. 136, scaled to radians
	ηT = []float64{47.0029 * s, -0.06603 * s, 0.000598 * s}
	πT = []float64{174.876384 * d, 3289.4789 * s, 0.60622 * s}
	pT = []float64{5029.0966 * s, 2.22226 * s, -0.000042 * s}

	// coefficients from (21.6) p. 136, scaled to radians
	ηt = []float64{47.0029 * s, -0.03302 * s, 0.000060 * s}
	πt = []float64{174.876384 * d, -869.8089 * s, 0.03536 * s}
	pt = []float64{5029.0966 * s, 1.11113 * s, -0.000006 * s}
)

// NewPrecessor constructs a Precessor object and initializes it to precess
// coordinates from epochFrom to epochTo.
func NewPrecessor(epochFrom, epochTo float64) *Precessor {
	// (21.2) p. 134
	ζCoeff := ζt
	zCoeff := zt
	θCoeff := θt
	if epochFrom != 2000 {
		T := (epochFrom - 2000) * .01
		ζCoeff = []float64{
			base.Horner(T, ζT...),
			0.30188*s - 0.000344*s*T,
			0.017998 * s}
		zCoeff = []float64{
			base.Horner(T, zT...),
			1.09468*s + 0.000066*s*T,
			0.018203 * s}
		θCoeff = []float64{
			base.Horner(T, θT...),
			-0.42665*s - 0.000217*s*T,
			-0.041833 * s}
	}
	t := (epochTo - epochFrom) * .01
	p := &Precessor{
		ζ: unit.RA(base.Horner(t, ζCoeff...) * t),
		z: unit.Angle(base.Horner(t, zCoeff...) * t),
	}
	θ := base.Horner(t, θCoeff...) * t
	p.sθ, p.cθ = math.Sincos(θ)
	return p
}

// Precess precesses coordinates eqFrom, leaving result in eqTo.
//
// The same struct may be used for eqFrom and eqTo.
// EqTo is returned for convenience.
func (p *Precessor) Precess(eqFrom, eqTo *coord.Equatorial) *coord.Equatorial {
	// (21.4) p. 134
	sδ, cδ := eqFrom.Dec.Sincos()
	sαζ, cαζ := (eqFrom.RA + p.ζ).Sincos()
	A := cδ * sαζ
	B := p.cθ*cδ*cαζ - p.sθ*sδ
	C := p.sθ*cδ*cαζ + p.cθ*sδ
	eqTo.RA = unit.RAFromRad(math.Atan2(A, B) + p.z.Rad())
	if math.Abs(C) < base.CosSmallAngle {
		eqTo.Dec = unit.Angle(math.Asin(C))
	} else {
		eqTo.Dec = unit.Angle(math.Acos(math.Hypot(A, B))) // near pole
		if C < 0 {
			eqTo.Dec = -eqTo.Dec
		}
	}
	return eqTo
}

// NewEclipticPrecessor constructs an EclipticPrecessor object and initializes
// it to precess coordinates from epochFrom to epochTo.
func NewEclipticPrecessor(epochFrom, epochTo float64) *EclipticPrecessor {
	// (21.5) p. 136
	ηCoeff := ηt
	πCoeff := πt
	pCoeff := pt
	if epochFrom != 2000 {
		T := (epochFrom - 2000) * .01
		ηCoeff = []float64{
			base.Horner(T, ηT...),
			-0.03302*s + 0.000598*s*T,
			0.000060 * s}
		πCoeff = []float64{
			base.Horner(T, πT...),
			-869.8089*s - 0.50491*s*T,
			0.03536 * s}
		pCoeff = []float64{
			base.Horner(T, pT...),
			1.11113*s - 0.000042*s*T,
			-0.000006 * s}
	}
	t := (epochTo - epochFrom) * .01
	p := &EclipticPrecessor{
		π: unit.Angle(base.Horner(t, πCoeff...)),
		p: unit.Angle(base.Horner(t, pCoeff...) * t),
	}
	η := unit.Angle(base.Horner(t, ηCoeff...) * t)
	p.sη, p.cη = η.Sincos()
	return p
}

// EclipticPrecess precesses coordinates eclFrom, leaving result in eclTo.
//
// The same struct may be used for eclFrom and eclTo.
// EclTo is returned for convenience.
func (p *EclipticPrecessor) Precess(eclFrom, eclTo *coord.Ecliptic) *coord.Ecliptic {
	// (21.7) p. 137
	sβ, cβ := eclFrom.Lat.Sincos()
	sd, cd := (p.π - eclFrom.Lon).Sincos()
	A := p.cη*cβ*sd - p.sη*sβ
	B := cβ * cd
	C := p.cη*sβ + p.sη*cβ*sd
	eclTo.Lon = p.p + p.π - unit.Angle(math.Atan2(A, B))
	if math.Abs(C) < base.CosSmallAngle {
		eclTo.Lat = unit.Angle(math.Asin(C))
	} else {
		eclTo.Lat = unit.Angle(math.Acos(math.Hypot(A, B))) // near pole
		if C < 0 {
			eclTo.Lat = -eclTo.Lat
		}
	}
	return eclTo
}

// EclipticPosition precesses ecliptic coordinates from one epoch to another,
// including proper motions.
//
// While eclFrom is given as ecliptic coordinates, proper motions mα, mδ are
// still expected to be equatorial.  If proper motions are not to be considered
// or are not applicable, pass 0, 0.
//
// Both eclFrom and eclTo must be non-nil, although they may point to the same
// struct.  EclTo is returned for convenience.
func EclipticPosition(eclFrom, eclTo *coord.Ecliptic, epochFrom, epochTo float64, mα unit.HourAngle, mδ unit.Angle) *coord.Ecliptic {
	p := NewEclipticPrecessor(epochFrom, epochTo)
	*eclTo = *eclFrom
	if mα != 0 || mδ != 0 {
		mλ, mβ := eqProperMotionToEcl(mα, mδ, epochFrom, eclFrom)
		t := epochTo - epochFrom
		eclTo.Lon += mλ.Mul(t)
		eclTo.Lat += mβ.Mul(t)
	}
	return p.Precess(eclTo, eclTo)
}

// Position precesses equatorial coordinates from one epoch to another,
// including proper motions.
//
// If proper motions are not to be considered or are not applicable, pass 0, 0
// for mα, mδ
//
// Both eqFrom and eqTo must be non-nil, although they may point to the same
// struct.  EqTo is returned for convenience.
func Position(eqFrom, eqTo *coord.Equatorial, epochFrom, epochTo float64, mα unit.HourAngle, mδ unit.Angle) *coord.Equatorial {
	p := NewPrecessor(epochFrom, epochTo)
	t := epochTo - epochFrom
	eqTo.RA = unit.RAFromRad(eqFrom.RA.Rad() + mα.Rad()*t)
	eqTo.Dec = eqFrom.Dec + mδ*unit.Angle(t)
	return p.Precess(eqTo, eqTo)
}

func eqProperMotionToEcl(mα unit.HourAngle, mδ unit.Angle, epoch float64, pos *coord.Ecliptic) (mλ, mβ unit.Angle) {
	ε := nutation.MeanObliquity(base.JulianYearToJDE(epoch))
	sε, cε := ε.Sincos()
	α, δ := coord.EclToEq(pos.Lon, pos.Lat, sε, cε)
	sα, cα := α.Sincos()
	sδ, cδ := δ.Sincos()
	cβ := pos.Lat.Cos()
	mλ = (mδ.Mul(sε*cα) + unit.Angle(mα).Mul(cδ*(cε*cδ+sε*sδ*sα))).Div(cβ * cβ)
	mβ = (mδ.Mul(cε*cδ+sε*sδ*sα) - unit.Angle(mα).Mul(sε*cα*cδ)).Div(cβ)
	return
}
