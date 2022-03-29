package williams_test

import (
	"testing"
	"webeph/testutils"
	"webeph/unit"
	"webeph/williams"
)

func TestFindMediumCoeli1(t *testing.T) {
	// Example, p51
	lst := unit.NewTime('+', 3, 42, 25).Angle()
	ε := unit.AngleFromDeg(23.448333)
	expect := unit.AngleFromDeg(57.868217)
	got := williams.FindMediumCoeli(lst, ε)
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), testutils.SecondTolerance) {
		t.Errorf("FindMC1: expected %v to be %v", got.Deg(), expect.Deg())
	}
}

func TestFindMediumCoeli2(t *testing.T) {
	// Example, p52
	lst := unit.NewTime('+', 10, 33, 14).Angle()
	ε := unit.AngleFromDeg(23.448889)
	expect := unit.AngleFromDeg(156.559047)
	got := williams.FindMediumCoeli(lst, ε)
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), testutils.SecondTolerance) {
		t.Errorf("FindMC2: expected %v to be %v", got.Deg(), expect.Deg())
	}
}

func TestFindMediumCoeli3(t *testing.T) {
	// Example, p53
	lst := unit.NewTime('+', 12, 22, 38).Angle()
	ε := unit.AngleFromDeg(23.439444)
	expect := unit.AngleFromDeg(186.16349)
	got := williams.FindMediumCoeli(lst, ε)
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), testutils.SecondTolerance) {
		t.Errorf("FindMC3: expected %v to be %v", got.Deg(), expect.Deg())
	}
}

func TestFindMediumCoeli4(t *testing.T) {
	// Example, p54
	lst := unit.NewTime('+', 22, 47, 7).Angle()
	ε := unit.AngleFromDeg(23.450278)
	expect := unit.AngleFromDeg(340.260849)
	got := williams.FindMediumCoeli(lst, ε)
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), testutils.SecondTolerance) {
		t.Errorf("FindMC4: expected %v to be %v", got.Deg(), expect.Deg())
	}
}

func TestFindMediumCoeli5(t *testing.T) {
	// Example, p55
	lst := unit.NewTime('+', 16, 13, 6).Angle()
	ε := unit.AngleFromDeg(23.44375)
	expect := unit.AngleFromDeg(245.2063)
	got := williams.FindMediumCoeli(lst, ε)
	if !testutils.CheckTolerance(got.Deg(), expect.Deg(), testutils.SecondTolerance) {
		t.Errorf("FindMC5: expected %v to be %v", got.Deg(), expect.Deg())
	}
}

func TestEPH571(t *testing.T) {
	// April 16, 2063@16:28:01.221 UTC
	lst := unit.Angle(2.9954091191224346)
	ε := unit.Angle(0.4089921762116411)
	got := williams.FindMediumCoeli(lst, ε)
	expected := unit.AngleFromDeg(170.88352662014552)
	if got.Deg() != expected.Deg() {
		t.Errorf("Expected %v to be %v", got.Deg(), expected.Deg())
	}
}

func TestEPH572(t *testing.T) {
	// April 16, 2063@16:29:01.221 UTC
	lst := unit.Angle(2.9997843886306192)
	ε := unit.Angle(0.408992176052117)
	got := williams.FindMediumCoeli(lst, ε)
	expected := unit.AngleFromDeg(171.1556926835057)
	if got.Deg() != expected.Deg() {
		t.Errorf("Expected %v to be %v", got.Deg(), expected.Deg())
	}
}
