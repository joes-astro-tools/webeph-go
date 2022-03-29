package elliptic_test

import (
	"fmt"

	elliptic "webeph/elliptic"
	pp "webeph/planetposition"
	sexa "webeph/sexagesimal"
)

func ExampleEclipticPosition() {
	// Example 33.a, p. 225.  VSOP87 result p. 227.
	earth, err := pp.LoadPlanet(pp.Earth)
	if err != nil {
		fmt.Println(err)
		return
	}
	venus, err := pp.LoadPlanet(pp.Venus)
	if err != nil {
		fmt.Println(err)
		return
	}
	λ, β, Δ := elliptic.EclipticPosition(venus, earth, 2448976.5, false)

	fmt.Printf("λ = %.3d\n", sexa.FmtAngle(λ))
	fmt.Printf("β = %.2d\n", sexa.FmtAngle(β))
	fmt.Printf("Δ = %.6f\n", Δ)
	// Output:
	// λ = 313°4′52″.838
	// β = -2°5′5″.36
	// Δ = 0.910948
}
