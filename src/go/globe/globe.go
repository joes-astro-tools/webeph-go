// Globe: Chapter 11, The Earth's Globe.
//
// Globe contains functions concerning the surface of the Earth idealized as
// an ellipsoid of revolution.
package globe

import (
	"math"

	unit "webeph/unit"
)

// Ellipsoid represents an ellipsoid of revolution.
//
// Typical unit for Er is Km.
type Ellipsoid struct {
	Er float64 // equatorial radius
	Fl float64 // flattening
}

// IAU 1976 values.  Radius in Km.
var Earth76 = Ellipsoid{Er: 6378.14, Fl: 1 / 298.257}

// ParallaxConstants computes parallax constants ρ sin φ′ and ρ cos φ′.
//
// Arguments are geographic latitude φ and height h above the ellipsoid.
// For e.Er in Km, h must be in meters.
func (e Ellipsoid) ParallaxConstants(φ unit.Angle, h float64) (s, c float64) {
	boa := 1 - e.Fl
	su, cu := math.Sincos(math.Atan(boa * φ.Tan()))
	s, c = φ.Sincos()
	hoa := h * 1e-3 / e.Er
	return su*boa + hoa*s, cu + hoa*c
}

// Coord represents geographic coordinates on the Earth.
//
// Longitude is measured positively westward from the Greenwich meridian.
type Coord struct {
	Lat unit.Angle // latitude (φ)
	Lon unit.Angle // longitude (ψ, or L)
}
