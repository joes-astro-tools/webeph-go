package parallactic

import (
	"math"

	unit "webeph/unit"
	williams "webeph/williams"
)

func FindAscendant(ε, φ, θ unit.Angle) (λ unit.Angle) {
	geocentricLat := williams.FindGeocentricLat(φ)
	sε, cε := ε.Sincos()
	sφ, cφ := geocentricLat.Sincos()
	sθ, cθ := θ.Sincos()
	tφ := sφ / cφ
	halfCircle := unit.Angle(math.Pi)
	// (14.2) p. 99
	λ = unit.Angle(math.Atan2(-cθ, (sε*tφ)+(cε*sθ))).Mod1()
	if λ.LessThan(halfCircle) {
		λ = λ.Add(halfCircle)
	} else {
		λ = λ.Subtract(halfCircle)
	}
	return
}
