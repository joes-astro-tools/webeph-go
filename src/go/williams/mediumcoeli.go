package williams

import (
	"math"

	unit "webeph/unit"
)

// Finds the midheaven, or Medium Coeli.
// Receives:
//	lst: local sidereal time, as a unit.Angle
//	ε: obliquity, as a unit.Angle
// Returns:
//	medium coeli, as a unit.Angle
func FindMediumCoeli(lst, ε unit.Angle) unit.Angle {
	// Right ascension of the Medium Coeli
	ramc := lst

	cosε := ε.Cos()
	if q1(lst) {
		a := ramc.Tan() / cosε
		return unit.Angle(math.Atan(a))
	}
	if q2(lst) {
		tangent := ramc.Subtract(angle90).Tan()
		a := unit.Angle(math.Atan(cosε * tangent))
		return angle90.Add(a)
	}
	if q3(lst) {
		tangent := ramc.Subtract(angle180).Tan()
		return unit.Angle(q2C + math.Atan2(tangent, cosε))
	}
	tangent := ramc.Subtract(angle270).Tan()
	return angle270.Add(unit.Angle(math.Atan(cosε * tangent)))
}
