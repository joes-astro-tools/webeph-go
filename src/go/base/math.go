package base

import unit "webeph/unit"

// SmallAngle is threshold used by various routines for switching between
// trigonometric functions and Pythagorean approximations.
//
// In chapter 17, p. 109, Meeus recommends 10′.
var (
	SmallAngle    = unit.AngleFromMin(10) // about .003 radians
	CosSmallAngle = SmallAngle.Cos()      // about .999996
)

// Horner evaluates a polynomal with coefficients c at x.  The constant
// term is c[0].  The function panics with an empty coefficient list.
func Horner(x float64, c ...float64) float64 {
	i := len(c) - 1
	y := c[i]
	for i > 0 {
		i--
		y = y*x + c[i] // sorry, no fused multiply-add in Go
	}
	return y
}

// FloorDiv returns the integer floor of the fractional value (x / y).
//
// It uses integer math only, so is more efficient than using floating point
// intermediate values.  This function can be used in many places where INT()
// appears in AA.  As with built in integer division, it panics with y == 0.
func FloorDiv(x, y int) (q int) {
	q = x / y
	if (x < 0) != (y < 0) && x%y != 0 {
		q--
	}
	return
}

// FloorDiv64 returns the integer floor of the fractional value (x / y).
//
// It uses integer math only, so is more efficient than using floating point
// intermediate values.  This function can be used in many places where INT()
// appears in AA.  As with built in integer division, it panics with y == 0.
func FloorDiv64(x, y int64) (q int64) {
	q = x / y
	if (x < 0) != (y < 0) && x%y != 0 {
		q--
	}
	return
}
