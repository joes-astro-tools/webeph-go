package zabinski_test

import (
	"testing"

	nutation "webeph/nutation"
	testutils "webeph/testutils"
	zabinski "webeph/zabinski"
)

func TestFindStellarLongitude1(t *testing.T) {
	// 01/26/22@1510 EST
	jd := 2459606.340277778
	_, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	expected := 56.466667
	// Algol
	// data comes from https://vizier.u-strasbg.fr/viz-bin/VizieR-4
	got := zabinski.FindStellarLongitude(jd, ε, 3, 8, 10.131, 40, 57, 20.43, 0.0031, -0.0009)
	if !testutils.CheckTolerance(got, expected, testutils.StandardTolerance) {
		t.Errorf("TestFindStellarLongitude1: expected %v to be %v", got, expected)
	}
}

func TestFindStellarLongitude2(t *testing.T) {
	// 01/26/22@1510 EST
	jd := 2459606.340277778
	_, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	expected := 204.15
	// Spica
	// data comes from https://vizier.u-strasbg.fr/viz-bin/VizieR-4
	got := zabinski.FindStellarLongitude(jd, ε, 13, 25, 11.587, -11, 9, 40.71, -0.00278, -0.00283)
	if !testutils.CheckTolerance(got, expected, testutils.StandardTolerance) {
		t.Errorf("TestFindStellarLongitude2: expected %v to be %v", got, expected)
	}
}

func TestFindStellarLongitude3(t *testing.T) {
	// 01/26/22@1510 EST
	jd := 2459606.340277778
	_, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	expected := 204.533334
	// Arcturus
	// data comes from https://vizier.u-strasbg.fr/viz-bin/VizieR-4
	got := zabinski.FindStellarLongitude(jd, ε, 14, 15, 39.677, 19, 10, 56.71, -0.0714, -1.9984)
	if !testutils.CheckTolerance(got, expected, testutils.StandardTolerance) {
		t.Errorf("TestFindStellarLongitude3: expected %v to be %v", got, expected)
	}
}

func TestFindStellarLongitude4(t *testing.T) {
	// 01/26/22@1510 EST
	jd := 2459606.340277778
	_, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	expected := 14.45
	// Linteum
	// data comes from https://vizier.u-strasbg.fr/viz-bin/VizieR-4
	got := zabinski.FindStellarLongitude(jd, ε, 0, 48, 40.95, 7, 35, 6.14, 0.00572, -0.0524)
	if !testutils.CheckTolerance(got, expected, testutils.StandardTolerance) {
		t.Errorf("TestFindStellarLongitude4: expected %v to be %v", got, expected)
	}
}
