// Copyright 2013 Sonia Keys
// License: MIT

package parallax_test

import (
	"fmt"
	"math"
	"testing"

	base "webeph/base"
	julian "webeph/julian"
	moonposition "webeph/moonposition"
	parallax "webeph/parallax"
	sexa "webeph/sexagesimal"
	unit "webeph/unit"
)

func ExampleHorizontal() {
	// Example 40.a, p. 280
	π := parallax.Horizontal(.37276)
	fmt.Printf("%.3s", sexa.FmtAngle(π))
	// Output:
	// 23.592″
}

func TestHorizontal(t *testing.T) {
	// example from moonposition.Parallax, ch 47, p. 342
	_, _, Δ := moonposition.Position(julian.CalendarGregorianToJD(1992, 4, 12))
	π := parallax.Horizontal(Δ / base.AU).Deg()
	want := .99199
	// we don't quite get all the digits here.
	// for close objects we need that Arcsin that's in moonposition.Parallax.
	if math.Abs(π-want) > .0001 {
		t.Fatal("got", π, "want", want)
	}
}

func ExampleTopocentricLongitude() {
	// exercise, p. 282
	λʹ := parallax.TopocentricLongitude(
		unit.NewAngle(' ', 181, 46, 22.5),
		unit.NewAngle(' ', 2, 17, 26.2),
		unit.NewAngle(' ', 50, 5, 7.8),
		0,
		unit.NewAngle(' ', 23, 28, 0.8),
		unit.NewAngle(' ', 209, 46, 7.9).Time(),
		unit.NewAngle(' ', 0, 59, 27.7))
	fmt.Printf("λʹ = %.1s\n", sexa.FmtAngle(λʹ))
	// Output:
	// λʹ = 181°48′5.0″
}
