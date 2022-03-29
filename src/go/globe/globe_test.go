package globe_test

import (
	"fmt"

	globe "webeph/globe"
	unit "webeph/unit"
)

func ExampleEllipsoid_ParallaxConstants() {
	// Example 11.a, p 82.
	// phi = geographic latitude of Palomar
	φ := unit.NewAngle(' ', 33, 21, 22)
	s, c := globe.Earth76.ParallaxConstants(φ, 1706)
	fmt.Printf("ρ sin φ′ = %+.6f\n", s)
	fmt.Printf("ρ cos φ′ = %+.6f\n", c)
	// Output:
	// ρ sin φ′ = +0.546861
	// ρ cos φ′ = +0.836339
}
