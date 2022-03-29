package williams

import (
	"math"

	unit "webeph/unit"
)

// Finds geocentric latitude, ie latitude to the equator at the center of the earth.
// Receives:
//	φ: geographic latitude, as a unit.Angle
// Returns:
//	geocentric latitude
func FindGeocentricLat(φ unit.Angle) unit.Angle {
	return unit.Angle(math.Atan(0.9932774 * math.Tan(φ.Rad())))
}
