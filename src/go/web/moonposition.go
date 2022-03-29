package web

import (
	apparent "webeph/apparent"
	moonposition "webeph/moonposition"
	nutation "webeph/nutation"
	unit "webeph/unit"
)

func MoonPosition(jde float64) (λ, β unit.Angle, Δ float64) {
	geocentricλ, geocentricβ, geocentricΔ := moonposition.Position(jde)
	Δψ, _ := nutation.Nutation(jde)
	Δλ, Δβ := apparent.EclipticAberration(geocentricλ, geocentricβ, jde)
	λ = geocentricλ + Δλ + Δψ
	β = geocentricβ + Δβ
	Δ = geocentricΔ
	return
}
