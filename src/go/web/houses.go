//go:build js && wasm

package web

import (
	parallactic "webeph/parallactic"
	unit "webeph/unit"
	williams "webeph/williams"
)

var (
	houseContainer = [12]float64{}
)

// Gets the array containing the houses.
// Receives:
//	nothing
// Returns:
//	the address of the storage container for houses.
// Notes:
//	Used to send results back to Javascript, in place of the Go runtime's bloated syscall/js functionality.
//export getHouseContainer
func getHouseContainer() *[12]float64 {
	return &houseContainer
}

// Finds all houses.
// Receives:
//	lst: local sidereal time, as a unit.Angle
//	ε: obliquity, as a unit.Angle
//	φ: latitude, as a unit.Angle
// Returns:
//	nothing
// Notes:
//	Stores results in a private variable. Use getHouseContainer() to retrieve results.
//export findHouses
func FindHouses(lst, ε, φ unit.Angle) {
	angle180 := unit.AngleFromDeg(180)
	h1 := parallactic.FindAscendant(ε, φ, lst)
	h7 := h1.Add(angle180)
	h10 := williams.FindMediumCoeli(lst, ε)
	h4 := h10.Add(angle180)
	// house 2, house 3, house 5, house 6, house 8, house 9, house 11, house 12
	regHouses := williams.FindHouses(lst, ε, φ)
	h2 := regHouses[0]
	h3 := regHouses[1]
	h5 := regHouses[2]
	h6 := regHouses[3]
	h8 := regHouses[4]
	h9 := regHouses[5]
	h11 := regHouses[6]
	h12 := regHouses[7]
	houseAngles := [12]unit.Angle{h1, h2, h3, h4, h5, h6, h7, h8, h9, h10, h11, h12}
	houseDegs := [12]float64{}
	for i, v := range houseAngles {
		houseDegs[i] = v.Deg()
	}
	houseContainer = houseDegs
}
