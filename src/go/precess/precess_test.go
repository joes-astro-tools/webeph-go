package precess_test

import (
	"fmt"

	base "webeph/base"
	coord "webeph/coord"
	julian "webeph/julian"
	precess "webeph/precess"
	sexa "webeph/sexagesimal"
	unit "webeph/unit"
)

func ExampleEclipticPosition() {
	// Example 21.c, p. 137.
	eclFrom := &coord.Ecliptic{
		Lat: unit.AngleFromDeg(1.76549),
		Lon: unit.AngleFromDeg(149.48194),
	}
	eclTo := &coord.Ecliptic{}
	epochFrom := 2000.0
	epochTo := base.JDEToJulianYear(julian.CalendarJulianToJD(-214, 6, 30))
	precess.EclipticPosition(eclFrom, eclTo, epochFrom, epochTo, 0, 0)
	fmt.Printf("%.3f\n", eclTo.Lon.Deg())
	fmt.Printf("%+.3f\n", eclTo.Lat.Deg())
	// Output:
	// 118.704
	// +1.615
}

func ExamplePosition() {
	// Example 21.b, p. 135.
	eq := &coord.Equatorial{
		RA:  unit.NewRA(2, 44, 11.986),
		Dec: unit.NewAngle(' ', 49, 13, 42.48),
	}
	epochFrom := 2000.0
	jdTo := julian.CalendarGregorianToJD(2028, 11, 13.19)
	epochTo := base.JDEToJulianYear(jdTo)
	precess.Position(eq, eq, epochFrom, epochTo,
		unit.HourAngleFromSec(0.03425),
		unit.AngleFromSec(-0.0895))
	fmt.Printf("%0.3d\n", sexa.FmtRA(eq.RA))
	fmt.Printf("%+0.2d\n", sexa.FmtAngle(eq.Dec))
	// Output:
	// 2ʰ46ᵐ11ˢ.331
	// +49°20′54″.54
}
