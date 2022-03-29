package parallactic_test

import (
	"testing"
	"webeph/julian"
	"webeph/nutation"
	"webeph/parallactic"
	"webeph/testutils"
	"webeph/unit"
	"webeph/williams"
	"webeph/zabinski"
)

// Tests derived from AstroApp.
// AstroApp (and other software) uses geographic latitude, not geocentric longitude.
// This creates a max variability of about 11 minutes, according to Williams on p15.
// We aim for 15 minutes of accuracy in those cases.

func TestFindAscendant1(t *testing.T) {
	// Q1, northern
	// My birthday
	φ := unit.AngleFromDeg(42.)
	ο := unit.AngleFromDeg(-71.383333)
	jd := julian.CalendarGregorianToJD(1978, 4, 21.265972) // 4/21/1978@01:23, in UTC
	Δψ, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	lst := zabinski.FindSiderealTime(Δψ, Δε, jd, ο)
	expect := unit.AngleFromDeg(302.433333)
	got := parallactic.FindAscendant(ε, φ, lst.Angle())
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), 0.25) {
		t.Errorf("FindAsc1: expected %v to be %v", got.Deg(), expect.Deg())
	}
}

func TestFindAscendant2(t *testing.T) {
	// Q2, northern
	// My birthday
	φ := unit.AngleFromDeg(42.)
	ο := unit.AngleFromDeg(-71.383333)
	jd := julian.CalendarGregorianToJD(1978, 4, 21.515972) // 4/21/1978@07:23, in UTC
	Δψ, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	lst := zabinski.FindSiderealTime(Δψ, Δε, jd, ο)
	expect := unit.AngleFromDeg(76.866667)
	got := parallactic.FindAscendant(ε, φ, lst.Angle())
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), 0.25) {
		t.Errorf("FindAsc1: expected %v to be %v", got.Deg(), expect.Deg())
	}
}

func TestFindAscendant3(t *testing.T) {
	// Q3, northern
	// My birthday
	φ := unit.AngleFromDeg(42.)
	ο := unit.AngleFromDeg(-71.383333)
	jd := julian.CalendarGregorianToJD(1978, 4, 21.765972) // 4/21/1978@13:23, in UTC
	Δψ, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	lst := zabinski.FindSiderealTime(Δψ, Δε, jd, ο)
	expect := unit.AngleFromDeg(151.733333)
	got := parallactic.FindAscendant(ε, φ, lst.Angle())
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), 0.25) {
		t.Errorf("FindAsc1: expected %v to be %v", got.Deg(), expect.Deg())
	}
}

func TestFindAscendant4(t *testing.T) {
	// Q4, northern
	// My birthday
	φ := unit.AngleFromDeg(42.)
	ο := unit.AngleFromDeg(-71.383333)
	jd := julian.CalendarGregorianToJD(1978, 4, 22.015972) // 4/21/1978@19:23, in UTC
	Δψ, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	lst := zabinski.FindSiderealTime(Δψ, Δε, jd, ο)
	expect := unit.AngleFromDeg(222.133333)
	got := parallactic.FindAscendant(ε, φ, lst.Angle())
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), 0.25) {
		t.Errorf("FindAsc1: expected %v to be %v", got.Deg(), expect.Deg())
	}
}

func TestFindAscendant5(t *testing.T) {
	// Q1, southern
	// My birthday
	φ := unit.AngleFromDeg(-42.)
	ο := unit.AngleFromDeg(-71.383333)
	jd := julian.CalendarGregorianToJD(1978, 4, 21.265972) // 4/21/1978@01:23, in UTC
	Δψ, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	lst := zabinski.FindSiderealTime(Δψ, Δε, jd, ο)
	expect := unit.AngleFromDeg(331.35)
	got := parallactic.FindAscendant(ε, φ, lst.Angle())
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), 0.25) {
		t.Errorf("FindAsc1: expected %v to be %v", got.Deg(), expect.Deg())
	}
}

func TestFindAscendant6(t *testing.T) {
	// Example, p52
	// My birthday
	φ := unit.AngleFromDeg(-42.)
	ο := unit.AngleFromDeg(-71.383333)
	jd := julian.CalendarGregorianToJD(1978, 4, 21.515972) // 4/21/1978@07:23, in UTC
	Δψ, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	lst := zabinski.FindSiderealTime(Δψ, Δε, jd, ο)
	expect := unit.AngleFromDeg(41.766667)
	got := parallactic.FindAscendant(ε, φ, lst.Angle())
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), 0.25) {
		t.Errorf("FindAsc1: expected %v to be %v", got.Deg(), expect.Deg())
	}
}

func TestFindAscendant7(t *testing.T) {
	// Q3, southern
	// My birthday
	φ := unit.AngleFromDeg(-42.)
	ο := unit.AngleFromDeg(-71.383333)
	jd := julian.CalendarGregorianToJD(1978, 4, 21.765972) // 4/21/1978@13:23, in UTC
	Δψ, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	lst := zabinski.FindSiderealTime(Δψ, Δε, jd, ο)
	expect := unit.AngleFromDeg(123.05)
	got := parallactic.FindAscendant(ε, φ, lst.Angle())
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), 0.25) {
		t.Errorf("FindAsc1: expected %v to be %v", got.Deg(), expect.Deg())
	}
}

func TestFindAscendant8(t *testing.T) {
	// Q4, southern
	// My birthday
	φ := unit.AngleFromDeg(-42.)
	ο := unit.AngleFromDeg(-71.383333)
	jd := julian.CalendarGregorianToJD(1978, 4, 22.015972) // 4/21/1978@19:23, in UTC
	Δψ, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	lst := zabinski.FindSiderealTime(Δψ, Δε, jd, ο)
	expect := unit.AngleFromDeg(257.366667)
	got := parallactic.FindAscendant(ε, φ, lst.Angle())
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), 0.25) {
		t.Errorf("FindAsc1: expected %v to be %v", got.Deg(), expect.Deg())
	}
}

func TestFindAscendant10(t *testing.T) {
	// Discovered for EPH-53
	φ := unit.AngleFromDeg(9.369047571428567)
	ο := unit.AngleFromDeg(-71.516667)
	jd := julian.CalendarGregorianToJD(2022, 9, 21.447222) // 9/21/2022@6:44:56, time zone -04:00:00
	Δψ, Δε := nutation.Nutation(jd)
	ε := zabinski.FindObliquity(Δε, jd)
	lst := zabinski.FindSiderealTime(Δψ, Δε, jd, ο)
	expect := unit.AngleFromDeg(179.766667)
	got := williams.FindAscendant(lst.Angle(), ε, φ)
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), 0.25) {
		t.Errorf("FindAsc1: expected %v to be %v", got.Deg(), expect.Deg())
	}
}
