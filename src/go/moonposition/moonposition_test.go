// Copyright 2013 Sonia Keys
// License: MIT

package moonposition_test

import (
	"fmt"
	"math"
	"testing"

	base "webeph/base"
	julian "webeph/julian"
	moonposition "webeph/moonposition"
)

func ExamplePosition() {
	// Example 47.a, p. 342.
	λ, β, Δ := moonposition.Position(julian.CalendarGregorianToJD(1992, 4, 12))
	fmt.Printf("λ = %.6f\n", λ.Deg())
	fmt.Printf("β = %.6f\n", β.Deg())
	fmt.Printf("Δ = %.1f\n", Δ)
	// Output:
	// λ = 133.162655
	// β = -3.229126
	// Δ = 368409.7
}

func ExampleParallax() {
	// Example 47.a, p. 342.
	_, _, Δ := moonposition.Position(julian.CalendarGregorianToJD(1992, 4, 12))
	π := moonposition.Parallax(Δ)
	fmt.Printf("π = %.6f\n", π.Deg())
	// Output:
	// π = 0.991990
}

func TestParallax(t *testing.T) {
	// test case from ch 40, p. 280
	got := moonposition.Parallax(.37276 * base.AU).Sec()
	want := 23.592
	if math.Abs(got-want) > .001 {
		t.Error(got, want)
	}
}
