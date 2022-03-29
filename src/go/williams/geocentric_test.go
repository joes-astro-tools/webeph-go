package williams_test

import (
	"testing"
	"webeph/testutils"
	"webeph/unit"
	"webeph/williams"
)

func TestGeocentricLatitude1(t *testing.T) {
	// Example, p15
	expect := 0.
	got := williams.FindGeocentricLat(0.)
	if got.Deg() != expect {
		t.Errorf("GeoLat1: expected %v to be %v", got.Deg(), expect)
	}
}

func TestGeocentricLatitude2(t *testing.T) {
	// Example, p15
	expect := 24.852222
	got := williams.FindGeocentricLat(unit.AngleFromDeg(25.))
	if !testutils.CheckTolerance(got.Deg(), expect, testutils.SecondTolerance) {
		t.Errorf("GeoLat2: expected %v to be %v", got.Deg(), expect)
	}
}

func TestGeocentricLatitude3(t *testing.T) {
	// Example, p15
	expect := 44.806667
	got := williams.FindGeocentricLat(unit.AngleFromDeg(45.))
	if !testutils.CheckTolerance(got.Deg(), expect, testutils.SecondTolerance) {
		t.Errorf("GeoLat3: expected %v to be %v", got.Deg(), expect)
	}
}

func TestGeocentricLatitude4(t *testing.T) {
	// Example, p15
	expect := 69.875556
	got := williams.FindGeocentricLat(unit.AngleFromDeg(70.))
	if !testutils.CheckTolerance(got.Deg(), expect, testutils.SecondTolerance) {
		t.Errorf("GeoLat4: expected %v to be %v", got.Deg(), expect)
	}
}

func TestGeocentricLatitude5(t *testing.T) {
	// Example, p15
	expect := 90.
	got := williams.FindGeocentricLat(unit.AngleFromDeg(90.))
	if got.Deg() != expect {
		t.Errorf("GeoLat4: expected %v to be %v", got.Deg(), expect)
	}
}
