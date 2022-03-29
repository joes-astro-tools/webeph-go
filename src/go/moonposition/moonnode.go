package moonposition

import (
	"math"
	base "webeph/base"
	unit "webeph/unit"
)

type AstroTerms struct {
	d   float64
	m   float64
	mP  float64
	f   float64
	d2  float64
	f2  float64
	mP2 float64
	f4  float64
}

type NodeFunc func(AstroTerms) float64

type NodeTerm struct {
	coeff   float64
	sinFunc NodeFunc
}

var (
	nodeTerms = []NodeTerm{
		{
			coeff: -1.4979,
			sinFunc: func(at AstroTerms) float64 {
				return at.d2 - at.f2
			},
		},
		{
			coeff: -.15,
			sinFunc: func(at AstroTerms) float64 {
				return at.m
			},
		},
		{
			coeff: -.1226,
			sinFunc: func(at AstroTerms) float64 {
				return 2 * at.d
			},
		},
		{
			coeff: .1176,
			sinFunc: func(at AstroTerms) float64 {
				return 2 * at.f
			},
		},
		{
			coeff: -0.0801,
			sinFunc: func(at AstroTerms) float64 {
				return at.mP2 - at.f2
			},
		},
		{
			coeff: -0.0616,
			sinFunc: func(at AstroTerms) float64 {
				return at.d2 - at.m - at.f2
			},
		},
		{
			coeff: 0.049,
			sinFunc: func(at AstroTerms) float64 {
				return at.d2 - at.mP
			},
		},
		{
			coeff: 0.0409,
			sinFunc: func(at AstroTerms) float64 {
				return at.mP - at.f2
			},
		},
		{
			coeff: 0.0327,
			sinFunc: func(at AstroTerms) float64 {
				return at.mP
			},
		},
		{
			coeff: 0.0324,
			sinFunc: func(at AstroTerms) float64 {
				return at.d2 + at.m - at.f2
			},
		},
		{
			coeff: 0.0196,
			sinFunc: func(at AstroTerms) float64 {
				d4 := 4 * at.d
				return d4 - at.f4
			},
		},
		{
			coeff: 0.018,
			sinFunc: func(at AstroTerms) float64 {
				return at.d2 - at.mP - at.f2
			},
		},
		{
			coeff: 0.015,
			sinFunc: func(at AstroTerms) float64 {
				return at.d2 - at.mP2
			},
		},
		{
			coeff: -0.015,
			sinFunc: func(at AstroTerms) float64 {
				return at.d2 + at.mP - at.f2
			},
		},
		{
			coeff: -0.0078,
			sinFunc: func(at AstroTerms) float64 {
				return at.d2 - at.m
			},
		},
		{
			coeff: -0.0045,
			sinFunc: func(at AstroTerms) float64 {
				return at.d2 + at.mP
			},
		},
		{
			coeff: 0.0044,
			sinFunc: func(at AstroTerms) float64 {
				return at.mP + at.f2
			},
		},
		{
			coeff: -0.0042,
			sinFunc: func(at AstroTerms) float64 {
				return at.d - at.mP
			},
		},
		{
			coeff: -0.0031,
			sinFunc: func(at AstroTerms) float64 {
				return at.m - at.f2
			},
		},
		{
			coeff: 0.0031,
			sinFunc: func(at AstroTerms) float64 {
				return at.d2 - at.m - at.mP
			},
		},
		{
			coeff: 0.0029,
			sinFunc: func(at AstroTerms) float64 {
				return at.d2 - at.f4
			},
		},
		{
			coeff: 0.0028,
			sinFunc: func(at AstroTerms) float64 {
				return at.m + at.f2
			},
		},
	}
)

// Node returns longitude of the mean ascending node of the lunar orbit.
func Node(jde float64) unit.Angle {
	return unit.AngleFromDeg(base.Horner(base.J2000Century(jde),
		125.0445479, -1934.1362891, .0020754, 1./467441., -1./60616000.)).Mod1()
}

// Finds the instantaneous ascending lunar node.
// Receives:
//	jd: the julian day, as a float64
// Returns:
//	the ecliptic longitude of the instantaneous ascending lunar node.
// Notes:
//	The instantaneous ascending node is the theoretical intersection between the Moon's path, given its current trajectory,
//	and the ecliptic. It is NOT the last place where the Moon had zero ecliptic latitude.
//	A. The last place could have actually been the descending node.
//	B. Since the Moon is so wide, it likely took her the traversal of multiple degrees to actually get across the ecliptic latitude line.
//	Which one is correct in the way that Renaissance astrologers imagined it? Don't know.
//	Swiss Eph uses the instantaneous node. Not that they are correct necessarily, but people probably expect me to match Swiss Eph for now.
//	Also, this algorithm is developed from Chapront's 1991 book, corrected in 1998, partially described by Meeus in 2005.
//	Nobody gives examples, so I am left to rely on AstroApp to approximate the data. Keep in mind that there is no guarantee that AstroApp is correct either.
//	JPL gives dates when they find that the Moon is at zero latitude, but testing against AstroApp, they don't agree,
//	mostly because AstroApp is not attempting to be as precise as JPL. In any case, JPL would be the historical node, not the instantaneous node.
//	Because Chapront and Meeus don't describe their reference frame for this algorithm, there is no rigorous way to reconcile them to JPL.
//export findAscendingNode
func FindAscendingNode(jd float64) float64 {
	d, m, mP, f := dmf(base.J2000Century(jd))
	astroterms := AstroTerms{
		d:   d,
		m:   m,
		mP:  mP,
		f:   f,
		d2:  2 * d,
		f2:  2 * f,
		mP2: 2 * mP,
		f4:  4 * f,
	}
	adjust := 0.
	for _, term := range nodeTerms {
		adjust += term.coeff * math.Sin(term.sinFunc(astroterms))
	}
	return Node(jd).Add(unit.AngleFromDeg(adjust)).Deg()

}
