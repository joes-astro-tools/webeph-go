// Elliptic: Chapter 33, Elliptic Motion.
//
// Partial: Various formulas and algorithms are unimplemented for lack of
// examples or test cases.
package elliptic

import (
	"math"

	apparent "webeph/apparent"
	base "webeph/base"
	pp "webeph/planetposition"
	schlyter "webeph/schlyter"
	unit "webeph/unit"
)

// Find a heliocentric position for a planet.
// Receives:
//	p: V87Planet object for the observed planet
//	earth: V87Planet object for the earth
//	jde: Julian day
// Returns:
//	L: heliocentric ecliptic longitude, as an Angle
//	B: heliocentric ecliptic latitude, as an Angle
//	R: distance, in AU
func findHeliocentricPosition(p *pp.V87Planet, jde float64, useSchlyter bool) (L, B unit.Angle, R float64) {
	if useSchlyter {
		if p.Ibody == pp.Mercury {
			return schlyter.HeliocentricMercury(jde)
		}
		if p.Ibody == pp.Venus {
			return schlyter.HeliocentricVenus(jde)
		}
	}
	return p.Position(jde)
}

// EclipticPosition returns observed ecliptic coordinates of a planet at a given time.
// Receives:
//	p: V87Planet object for the observed planet
//	earth: V87Planet object for the earth
//	jde: Julian day
// Returns:
//	λ: ecliptic longitude, as an Angle
// 	β: ecliptic latitude, as an Angle
//	Δ: distance from the planet to the Earth, in AU
func EclipticPosition(p, earth *pp.V87Planet, jde float64, useSchlyter bool, Δψ unit.Angle) (λ, β unit.Angle, Δ float64) {
	L0, B0, R0 := earth.Position(jde)
	var L, B unit.Angle
	var R float64
	L, B, R = findHeliocentricPosition(p, jde, useSchlyter)
	sB0, cB0 := B0.Sincos()
	sL0, cL0 := L0.Sincos()
	sB, cB := B.Sincos()
	sL, cL := L.Sincos()
	x := R*cB*cL - R0*cB0*cL0
	y := R*cB*sL - R0*cB0*sL0
	z := R*sB - R0*sB0
	Δ = math.Sqrt(x*x + y*y + z*z) // (33.4) p. 224
	{
		τ := base.LightTime(Δ)
		// repeating with jde-τ
		L, B, R = findHeliocentricPosition(p, jde-τ, useSchlyter)
		sB, cB = B.Sincos()
		sL, cL = L.Sincos()
		x = R*cB*cL - R0*cB0*cL0
		y = R*cB*sL - R0*cB0*sL0
		z = R*sB - R0*sB0
		Δ = math.Sqrt(x*x + y*y + z*z)
	}
	λ = unit.Angle(math.Atan2(y, x))                // (33.1) p. 223
	β = unit.Angle(math.Atan2(z, math.Hypot(x, y))) // (33.2) p. 223
	Δλ, Δβ := apparent.EclipticAberration(λ, β, jde)
	λ, β = pp.ToFK5(λ+Δλ, β+Δβ, jde)
	λ += Δψ
	λ = λ.Mod1()
	return
}
