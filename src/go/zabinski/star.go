package zabinski

import (
	"math"

	"webeph/apparent"
	base "webeph/base"
	coord "webeph/coord"
	unit "webeph/unit"
)

// Find geocentric ecliptic longitude for a star.
// Receives:
//	jd: Julian day
//	ε: Obliquity, as a unit.Angle
//	raH: Right ascension hours
//	raM: Right ascension minutes
//	raS: Right ascension seconds
//	declD: Declination degrees
//	declM: Declination minutes
//	declS: Declination seconds
//	raμ: Change in right ascension, in arc seconds per year
//	declμ: Change in declination, in arc seconds per year
// Returns:
//	geocentric ecliptic longitude
// Notes:
//	1. The parameters are as granular as they are because star catalogs break it down this way. As we come closer to calculating
//	 fixed star elections, we may change to decimal degrees for both measures, depending on what we can do with the star catalog.
//	2. Geocentric and topocentric measure seem to be treated as the same for stellar measure. I find no explanation for this, but
//   the difference between topocentric and geocentric measures is due to parallax, which is a function of distance. This makes a
//   big difference for the Moon, and a small difference for planets. (2 degrees for the Moon, and only 8 seconds at maximum for
//   planets.) Seeing that we are dealing with a fraction of a degree for the planets, and the stars are much farther away, the
//   difference is likely so small that it is neglected.
//export findStellarLongitude
func FindStellarLongitude(jd float64, ε unit.Angle, raH, raM int, raS float64, declD, declM int, declS, raμ, declμ float64) float64 {
	declSign := '-'
	if declD > 0 {
		declSign = '+'
	}
	declD = int(math.Abs(float64(declD)))
	eq := &coord.Equatorial{
		RA:  unit.NewRA(raH, raM, raS),
		Dec: unit.NewAngle(byte(declSign), declD, declM, declS),
	}
	epochFrom := 2000.0
	epochTo := base.JDEToJulianYear(jd)
	// Find the geocentric position of the star. (apparent.Position takes care of the heliocentric position first.)
	eqTo := apparent.Position(
		eq,
		eq,
		epochFrom,
		epochTo,
		unit.HourAngleFromSec(raμ),
		unit.AngleFromSec(declμ),
	)
	// Convert to ecliptic longitude.
	obl := coord.NewObliquity(ε)
	geocentricλ, _ := coord.EqToEcl(eqTo.RA, eqTo.Dec, obl.S, obl.C)
	return geocentricλ.Deg()
}
