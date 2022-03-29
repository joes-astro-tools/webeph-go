package zabinski

import (
	"math"

	base "webeph/base"
	unit "webeph/unit"
)

var (
	sunRiseSet = [2]float64{}
	circle     = unit.Angle(2 * math.Pi)
)

// Retrieves the pointer where sunrise and sunset are stored.
// Receives:
//	nothing
// Returns:
//	a pointer to sunrise/sunset
// Notes:
//	This is done to support the indirection used by TinyGo in WASM when more than one value needs to be produced at a time.
//	This may change when TinyGo starts supporting WASM's multivalue someday.
//export getSunRiseSetPtr
func GetSunRiseSetPtr() *[2]float64 {
	return &sunRiseSet
}

// Calculates the times of sunrise and sunset.
// Receives:
//	jde: the Julian day
//	φ: the geographic latitude, as a unit.Angle
//	ο: the geographic longitude
// Returns:
//	nothing
// Notes:
//	Stores the calculation in a private, package-level variable. Retrieve results by calling GetSunRiseSetPtr.
//	Developed from the spec in https://en.wikipedia.org/wiki/Sunrise_equation.
//export findSunRiseSet
func FindSunRiseSet(jde float64, φ, ο unit.Angle) {
	// 1. Find Julian day since the year 2000.
	currJulianDay := math.Ceil(jde - base.J2000 + 0.0008)
	// 2. Find mean solar time.
	meanSolarTime := currJulianDay - ο.Div(circle.Rad()).Rad()
	// 3. Find solar anomaly.
	m := unit.Angle((0.01720196999454 * meanSolarTime) + 6.2400599667).Mod1()
	m2 := unit.Angle(2 * m).Mod1()
	m3 := unit.Angle(3 * m).Mod1()
	// 4. Find equation of the center.
	eqCenter := (1.9148 * m.Sin()) + (0.02 * m2.Sin()) + (0.0003 * m3.Sin())
	// 5. Find ecliptic longitude of the Sun.
	λ := unit.AngleFromDeg(m.Deg() + eqCenter + 180 + 102.9372).Mod1()
	λ2 := λ.Mul(2).Mod1()
	// 6. Find solar noon.
	transit := base.J2000 + meanSolarTime + (0.0053 * m.Sin()) - (0.0069 * λ2.Sin())
	// 7. Find declination.
	δ := unit.Angle(math.Asin(λ.Sin() * 0.39778850739795))
	// 8. Find solar hour angle.
	hourAngle := unit.Angle(math.Acos((-0.014485726138606 - (φ.Sin() * δ.Sin())) / (φ.Cos() * δ.Cos())))
	// 9. Sunrise and sunset are always equal time from noon. Use the solar hour angle to calculate the number of hours from noon.
	hours := hourAngle.Div(circle.Rad()).Rad()
	sunRiseSet = [2]float64{transit - hours, transit + hours}
}
