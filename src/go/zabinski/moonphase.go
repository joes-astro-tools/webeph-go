package zabinski

import (
	"math"

	base "webeph/base"
	moonposition "webeph/moonposition"
	solar "webeph/solar"
)

// Finds the lunar phase, expressed as the difference between solar and lunar geocentric ecliptic longitude.
// Receives:
//	jd: the Julian day
// Returns:
//	the difference between solar and lunar longitude, in degrees. Rounds to the nearest degree.
// Notes:
//	Measures distance by subtracting solar longitude from lunar longitude. This has the effect of expressing waxing/waning:
//	0: new
//	0-179: waxing
//	180: full
//	181-359 waning
//export findMoonPhase
func FindMoonPhase(jd float64) float64 {
	solarGeocentricλ := solar.ApparentLongitude(base.J2000Century(jd))
	lunarGeocentricλ, _, _ := moonposition.Position(jd)
	diff := math.Round(lunarGeocentricλ.Deg() - solarGeocentricλ.Deg())
	if diff < 0 {
		diff += 360
	}
	return diff
}
