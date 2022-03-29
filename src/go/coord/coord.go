// Coord: Chapter 13, Transformation of Coordinates.
//
// Transforms in this package are provided in two forms, function and method.
// The results of the two forms should be identical.
//
// The function forms pass all arguments and results as single values.  These
// forms are best used when you are transforming a single pair of coordinates
// and wish to avoid memory allocation.
//
// The method forms take and return pointers to structs.  These forms are best
// used when you are transforming multiple coordinates and can reuse one or
// more of the structs.  In this case reuse of structs will minimize
// allocations, and the struct pointers will pass more efficiently on the
// stack.  These methods transform their arguments, placing the result in
// the receiver.  The receiver is then returned for convenience.
package coord

import (
	"math"

	unit "webeph/unit"
)

// Obliquity represents the obliquity of the ecliptic.
type Obliquity struct {
	S, C float64 // sine and cosine of obliquity
}

// NewObliquity constructs a new Obliquity.
//
// Struct members are initialized from the given value ε of the obliquity of
// the ecliptic.
func NewObliquity(ε unit.Angle) *Obliquity {
	r := &Obliquity{}
	r.S, r.C = ε.Sincos()
	return r
}

// Ecliptic coordinates are referenced to the plane of the ecliptic.
type Ecliptic struct {
	Lon unit.Angle // Longitude (λ)
	Lat unit.Angle // Latitude (β)
}

// Equatorial coordinates are referenced to the Earth's rotational axis.
type Equatorial struct {
	RA  unit.RA    // Right ascension (α)
	Dec unit.Angle // Declination (δ)
}

// Keep: EclToEq converts ecliptic coordinates to equatorial coordinates.
func (eq *Equatorial) EclToEq(ecl *Ecliptic, ε *Obliquity) *Equatorial {
	eq.RA, eq.Dec = EclToEq(ecl.Lon, ecl.Lat, ε.S, ε.C)
	return eq
}

// EqToEcl converts equatorial coordinates to ecliptic coordinates.
func (ecl *Ecliptic) EqToEcl(eq *Equatorial, ε *Obliquity) *Ecliptic {
	ecl.Lon, ecl.Lat = EqToEcl(eq.RA, eq.Dec, ε.S, ε.C)
	return ecl
}

// EclToEq converts ecliptic coordinates to equatorial coordinates.
//
//	λ: ecliptic longitude coordinate to transform
//	β: ecliptic latitude coordinate to transform
//	sε: sine of obliquity of the ecliptic
//	cε: cosine of obliquity of the ecliptic
//
// Results:
//	α: right ascension
//	δ: declination
func EclToEq(λ, β unit.Angle, sε, cε float64) (α unit.RA, δ unit.Angle) {
	sλ, cλ := λ.Sincos()
	sβ, cβ := β.Sincos()
	α = unit.RAFromRad(math.Atan2(sλ*cε-(sβ/cβ)*sε, cλ)) // (13.3) p. 93
	δ = unit.Angle(math.Asin(sβ*cε + cβ*sε*sλ))          // (13.4) p. 93
	return
}

// EqToEcl converts equatorial coordinates to ecliptic coordinates.
//
//	α: right ascension coordinate to transform
//	δ: declination coordinate to transform
//	sε: sine of obliquity of the ecliptic
//	cε: cosine of obliquity of the ecliptic
//
// Results:
//
//	λ: ecliptic longitude
//	β: ecliptic latitude
func EqToEcl(α unit.RA, δ unit.Angle, sε, cε float64) (λ, β unit.Angle) {
	sα, cα := α.Sincos()
	sδ, cδ := δ.Sincos()
	λ = unit.Angle(math.Atan2(sα*cε+(sδ/cδ)*sε, cα)).Mod1() // (13.1) p. 93
	β = unit.Angle(math.Asin(sδ*cε - cδ*sε*sα)).Mod1()      // (13.2) p. 93
	return
}

var (
	// IAU B1950.0 coordinates of galactic North Pole
	GalacticNorth1950 = &Equatorial{
		RA:  unit.NewRA(12, 49, 0),
		Dec: unit.AngleFromDeg(27.4),
	}
	// Meeus gives 33 as the origin of galactic longitudes relative to the
	// ascending node of of the galactic equator.  33 + 90 = 123, the IAU
	// value for origin relative to the equatorial pole.
	Galactic0Lon1950 = unit.AngleFromDeg(33)
)

// Horizontal coordinates are referenced to the local horizon of an observer
// on the surface of the Earth.
type Horizontal struct {
	Az  unit.Angle // Azimuth (A)
	Alt unit.Angle // Altitude (h)
}

// Galactic coordinates are referenced to the plane of the Milky Way.
type Galactic struct {
	Lat unit.Angle // Latitude (b) in radians
	Lon unit.Angle // Longitude (l) in radians
}
