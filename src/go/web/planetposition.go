package web

import (
	pp "webeph/planetposition"
)

// LoadPlanetPath constructs a V87Planet object from a VSOP87 file.
//
// Argument ibody should be one of the planet constants; path should be
// a directory containing the VSOP87 files.
func LoadPlanet(ibody int) *pp.V87Planet {
	switch ibody {
	case pp.Mars:
		return pp.GetMars()
	case pp.Jupiter:
		return pp.GetJupiter()
	case pp.Saturn:
		return pp.GetSaturn()
	default:
		return pp.GetEarth()
	}
}
