//go:build js && wasm

package web

import (
	nutation "webeph/nutation"
	unit "webeph/unit"
	zabinski "webeph/zabinski"
)

var (
	angles = [2]unit.Angle{}
)

// Gets the array containing the obliquity and LST.
// Receives:
//	nothing
// Returns:
//	the address of the storage container for obliquity and LST.
// Notes:
//	Used to send results back to Javascript, in place of the Go runtime's bloated syscall/js functionality.
//export getObliquityLSTContainer
func GetObliquityLSTContainer() *[2]unit.Angle {
	return &angles
}

// Finds obliquity and local sidereal time.
// Receives:
//	jd: the Julian day
//	ο: the longitude, as a unit.Angle
// Returns:
//	nothing.
// Notes:
//	Results are stored in a container internal to the package. Use getObliquityLSTContainer to recover results.
//export findObliquityLST
func FindObliquityAndLST(jd float64, ο unit.Angle) {
	Δψ, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	lst := zabinski.FindSiderealTime(Δψ, Δε, jd, ο)
	angles = [2]unit.Angle{ε, lst.Angle()}
}
