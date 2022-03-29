package coord_test

import (
	"fmt"

	base "webeph/base"
	coord "webeph/coord"
	sexa "webeph/sexagesimal"
	unit "webeph/unit"
)

func ExampleEclToEq() {
	// Exercise, end of Example 13.a, p. 95.
	α, δ := coord.EclToEq(
		unit.AngleFromDeg(113.21563),
		unit.AngleFromDeg(6.68417),
		base.SOblJ2000,
		base.COblJ2000)
	fmt.Printf("α = %.3d, δ = %+.2d\n", sexa.FmtRA(α), sexa.FmtAngle(δ))
	// Output:
	// α = 7ʰ45ᵐ18ˢ.946, δ = +28°1′34″.26
}

func ExampleEquatorial_EclToEq() {
	// Exercise, end of Example 13.a, p. 95.
	ecl := &coord.Ecliptic{
		Lon: unit.AngleFromDeg(113.21563),
		Lat: unit.AngleFromDeg(6.68417),
	}
	ε := coord.NewObliquity(unit.AngleFromDeg(23.4392911))
	eq := new(coord.Equatorial).EclToEq(ecl, ε)
	fmt.Printf("α = %.3d, δ = %+.2d\n",
		sexa.FmtRA(eq.RA), sexa.FmtAngle(eq.Dec))
	// Output:
	// α = 7ʰ45ᵐ18ˢ.946, δ = +28°1′34″.26
}

func ExampleEcliptic_EqToEcl() {
	// Example 13.a, p. 95.
	eq := &coord.Equatorial{
		RA:  unit.NewRA(7, 45, 18.946),
		Dec: unit.NewAngle(' ', 28, 1, 34.26),
	}
	obl := coord.NewObliquity(unit.AngleFromDeg(23.4392911))
	ecl := new(coord.Ecliptic).EqToEcl(eq, obl)
	fmt.Printf("λ = %.5j\n", sexa.FmtAngle(ecl.Lon))
	fmt.Printf("β = %+.6j\n", sexa.FmtAngle(ecl.Lat))
	// Output:
	// λ = 113°.21563
	// β = +6°.684170
}

func ExampleEqToEcl() {
	// Example 13.a, p. 95 but using precomputed obliquity sine and cosine.
	λ, β := coord.EqToEcl(
		unit.NewRA(7, 45, 18.946),
		unit.NewAngle(' ', 28, 1, 34.26),
		base.SOblJ2000, base.COblJ2000)
	fmt.Printf("λ = %.5j\n", sexa.FmtAngle(λ))
	fmt.Printf("β = %+.6j\n", sexa.FmtAngle(β))
	// Output:
	// λ = 113°.21563
	// β = +6°.684170
}
