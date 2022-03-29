package nutation_test

import (
	"fmt"

	julian "webeph/julian"
	nutation "webeph/nutation"
	sexa "webeph/sexagesimal"
)

func ExampleNutation() {
	// Example 22.a, p. 148.
	jd := julian.CalendarGregorianToJD(1987, 4, 10)
	Δψ, Δε := nutation.Nutation(jd)
	ε0 := nutation.MeanObliquity(jd)
	ε := ε0 + Δε
	fmt.Printf("%+.3d\n", sexa.FmtAngle(Δψ))
	fmt.Printf("%+.3d\n", sexa.FmtAngle(Δε))
	fmt.Printf("%.3d\n", sexa.FmtAngle(ε0))
	fmt.Printf("%.3d\n", sexa.FmtAngle(ε))
	// Output:
	// -3″.788
	// +9″.443
	// 23°26′27″.407
	// 23°26′36″.850
}
