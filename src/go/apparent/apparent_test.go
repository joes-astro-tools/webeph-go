package apparent_test

import (
	"fmt"

	apparent "webeph/apparent"
	base "webeph/base"
	coord "webeph/coord"
	julian "webeph/julian"
	sexa "webeph/sexagesimal"
	unit "webeph/unit"
)

func ExampleNutation() {
	// Example 23.a, p. 152
	α := unit.NewRA(2, 46, 11.331)
	δ := unit.NewAngle(' ', 49, 20, 54.54)
	jd := julian.CalendarGregorianToJD(2028, 11, 13.19)
	Δα1, Δδ1 := apparent.Nutation(α, δ, jd)
	fmt.Printf("%.3s  %.3s\n",
		sexa.FmtAngle(unit.Angle(Δα1)), // (Δα1 is in HourAngle)
		sexa.FmtAngle(Δδ1))
	// Output:
	// 15.843″  6.217″
}

func ExampleAberration() {
	// Example 23.a, p. 152
	α := unit.NewRA(2, 46, 11.331)
	δ := unit.NewAngle(' ', 49, 20, 54.54)
	jd := julian.CalendarGregorianToJD(2028, 11, 13.19)
	Δα2, Δδ2 := apparent.Aberration(α, δ, jd)
	fmt.Printf("%.3s  %.3s\n",
		sexa.FmtAngle(unit.Angle(Δα2)), // (Δα2 is in HourAngle)
		sexa.FmtAngle(Δδ2))
	// Output:
	// 30.045″  6.697″
}

func ExamplePosition() {
	// Example 23.a, p. 152
	jd := julian.CalendarGregorianToJD(2028, 11, 13.19)
	eq := &coord.Equatorial{
		RA:  unit.NewRA(2, 44, 11.986),
		Dec: unit.NewAngle(' ', 49, 13, 42.48),
	}
	apparent.Position(eq, eq, 2000, base.JDEToJulianYear(jd),
		unit.HourAngleFromSec(.03425),
		unit.AngleFromSec(-.0895))
	fmt.Printf("α = %0.3d\n", sexa.FmtRA(eq.RA))
	fmt.Printf("δ = %0.2d\n", sexa.FmtAngle(eq.Dec))
	// Output:
	// α = 2ʰ46ᵐ14ˢ.390
	// δ = 49°21′07″.45
}
