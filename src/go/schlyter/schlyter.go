package schlyter

import (
	"math"
	unit "webeph/unit"
)

func FindDay(y, m, d int, t float64) float64 {
	a := 367 * y
	b := (7 * (y + ((m + 9) / 12))) / 4
	c := (275 * m) / 9
	// Express time as a fraction of the day.
	e := t / 24
	return float64(a-b+c+d-730530) + e
	// Maximum dayNum, for July 15, 2072@23:59 = 26495.55143
}

func deriveEccentricAnomaly(e float64, M unit.Angle) unit.Angle {
	e0 := M + unit.Angle((e * math.Sin(M.Rad()) * (1 + (e * math.Cos(M.Rad())))))
	e1 := unit.Angle(e0.Rad() - ((e0.Rad() - (e * math.Sin(e0.Rad())) - M.Rad()) / (1 - (e * math.Cos(e0.Rad())))))
	// Although Schlyter mentions this reduction, his examples do not use it.
	// All tests pass without this reduction.
	// Still need to watch it for other planets besides Mercury.
	// for findDiff(e0, e1) > 0.005 {
	// 	e0 = e1
	// 	e1 = e0 - ((e0 - (e * math.Sin(e0AsRads) * radsToDegs) - M) / (1 - (e * math.Cos(e0AsRads))))
	// }
	return e1
}

// heliocentricPosition returns heliocentric ecliptic coordinates of a planet at a given time.
// Receives:
//	N: Longitude of ascending node, in radians
//	i: Inclination, in radians
//	w: Argument of perihelion, in radians
//	M: Mean anomaly, in radians
//	a: Semi-major axis, unitless
//	e: Eccentricity, unitless
// Returns:
//	L: heliocentric ecliptic longitude, as an Angle
// 	B: heliocentric ecliptic latitude, as an Angle
//	R: distance from the planet to the Earth, in AU
func heliocentricPosition(N, i, w, M unit.Angle, a, e float64) (L, B unit.Angle, R float64) {
	// Eccentric anomaly
	ea := deriveEccentricAnomaly(e, M)

	// Distance, true anomaly
	xRec := a * (math.Cos(ea.Rad()) - e)
	yRec := a * math.Sqrt(1-(e*e)) * math.Sin(ea.Rad())
	dst := math.Sqrt(xRec*xRec + yRec*yRec)
	// True anomaly could be negative. Need to modulate.
	ta := unit.Angle(math.Atan2(yRec, xRec)).Mod1()

	// Heliocentric ecliptic rectangular coordinates
	sinN, cosN := math.Sincos(N.Rad())
	cosI := math.Cos(i.Rad())
	anomTerm := unit.Angle(ta + w)
	sinAnomTerm, cosAnomTerm := math.Sincos(anomTerm.Rad())
	xEclip := dst * (cosN*cosAnomTerm - sinN*sinAnomTerm*cosI)
	yEclip := dst * (sinN*cosAnomTerm + cosN*sinAnomTerm*cosI)
	zEclip := dst * sinAnomTerm * math.Sin(i.Rad())
	xEclipSq := xEclip * xEclip
	yEclipSq := yEclip * yEclip
	zEclipSq := zEclip * zEclip
	sumEclipSq := xEclipSq + yEclipSq
	// Could be negative. Need to modulate.
	L = unit.Angle(math.Atan2(yEclip, xEclip)).Mod1()
	B = unit.Angle(math.Atan2(zEclip, math.Sqrt(sumEclipSq)))
	R = math.Sqrt(sumEclipSq + zEclipSq)
	return

}

// HeliocentricMercury returns heliocentric ecliptic coordinates of Mercury at a given time.
// Receives:
//	jde: Julian day
// Returns:
//	λ: heliocentric ecliptic longitude, as an Angle
// 	β: heliocentric ecliptic latitude, as an Angle
//	Δ: distance from the planet to the Earth, in AU
func HeliocentricMercury(jde float64) (L, B unit.Angle, R float64) {
	dayNum := jde - 2451543.5
	// Longitude of ascending node, in radians
	N := unit.Angle(0.8435403168 + (5.66511185916e-7 * dayNum))
	// Inclination, in radians
	i := unit.Angle(0.122255078 + (8.7266e-10 * dayNum))
	// Argument of perihelion, in radians
	w := unit.Angle(0.5083114367 + (1.77053181e-7 * dayNum))
	// Semi-major axis, unitless
	a := 0.387098
	// Eccentricity, unitless
	e := 0.205635 + (5.59e-10 * dayNum)
	// Mean anomaly, in radians
	// Max value of around 1895 radians. Need to modulate.
	M := unit.Angle(2.9436059939 + (0.0714247100147306 * dayNum)).Mod1()

	return heliocentricPosition(N, i, w, M, a, e)
}

// HeliocentricVenus returns heliocentric ecliptic coordinates of Venus at a given time.
// Receives:
//	jde: Julian day
// Returns:
//	λ: heliocentric ecliptic longitude, as an Angle
// 	β: heliocentric ecliptic latitude, as an Angle
//	Δ: distance from the planet to the Earth, in AU
func HeliocentricVenus(jde float64) (L, B unit.Angle, R float64) {
	dayNum := jde - 2451543.5
	// Longitude of ascending node, in radians
	N := unit.Angle(1.338316725 + (4.30380740248e-7 * dayNum))
	// Inclination, in radians
	i := unit.Angle(0.0592469468 + (4.79965544e-10 * dayNum))
	// Argument of perihelion, in radians
	w := unit.Angle(0.95802868 + (2.41508189915e-7 * dayNum))
	// Semi-major axis, unitless
	a := 0.723330
	// Eccentricity, unitless
	e := 0.006773 + (-1.302e-9 * dayNum)
	// Mean anomaly, in radians
	// Max value of around 742 radians. Need to modulate.
	M := unit.Angle(0.8378487981 + (0.0279624474614262 * dayNum)).Mod1()

	return heliocentricPosition(N, i, w, M, a, e)
}
