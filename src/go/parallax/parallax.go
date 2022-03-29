// Copyright 2013 Sonia Keys
// License: MIT

// Parallax: Chapter 40, Correction for Parallax.
package parallax

import (
	"math"

	globe "webeph/globe"
	unit "webeph/unit"
)

// constant for Horizontal.  p. 279.
var hp = unit.AngleFromSec(8.794)

// Horizontal returns equatorial horizontal parallax of a body.
//
// Argument Δ is distance in AU.
//
// Meeus mentions use of this function for the Moon, Sun, planet, or comet.
// That is, for relatively distant objects.  For parallax of the Moon (or
// other relatively close object) see moonposition.Parallax.
func Horizontal(Δ float64) (π unit.Angle) {
	return hp.Div(Δ) // (40.1) p. 279
}

// Finds topocentric longitude.
// Receives:
//	λ: geocentric ecliptic longitude as an Angle
//	β: geocentric ecliptic latitude as an Angle
//	φ: geographic latitude as an Angle
//	h: geographic height from sea level, in meters
//	ε: obliquity as an Angle
//	θ: local sidereal time as a Time
//	π: equatorial horizontal parallax as an Angle
// Returns:
//	λʹ: topocentric longitude as an Angle
func TopocentricLongitude(λ, β, φ unit.Angle, h float64, ε unit.Angle, θ unit.Time, π unit.Angle) (λʹ unit.Angle) {
	S, C := globe.Earth76.ParallaxConstants(φ, h)
	sλ, cλ := λ.Sincos()
	cβ := β.Cos()
	sε, cε := ε.Sincos()
	sθ, cθ := θ.Angle().Sincos()
	sπ := π.Sin()
	N := cλ*cβ - C*sπ*cθ
	λʹ = unit.Angle(math.Atan2(sλ*cβ-sπ*(S*sε+C*cε*sθ), N))
	if λʹ < 0 {
		λʹ += 2 * math.Pi
	}
	return
}
