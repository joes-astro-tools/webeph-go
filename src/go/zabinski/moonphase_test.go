package zabinski_test

import (
	"testing"

	zabinski "webeph/zabinski"
)

func TestMoonPhase1(t *testing.T) {
	jd := 2459627.340277778
	expected := 182.
	got := zabinski.FindMoonPhase(jd)
	if got != expected {
		t.Errorf("TestMoonPhase1: expected %v to be %v", got, expected)
	}
}

func TestMoonPhase2(t *testing.T) {
	jd := 2459620.340277778
	expected := 104.
	got := zabinski.FindMoonPhase(jd)
	if got != expected {
		t.Errorf("TestMoonPhase2: expected %v to be %v", got, expected)
	}
}

func TestMoonPhase3(t *testing.T) {
	jd := 2459613.340277778
	expected := 21.
	got := zabinski.FindMoonPhase(jd)
	if got != expected {
		t.Errorf("TestMoonPhase3: expected %v to be %v", got, expected)
	}
}

func TestMoonPhase4(t *testing.T) {
	jd := 2459606.340277778
	expected := 286.
	got := zabinski.FindMoonPhase(jd)
	if got != expected {
		t.Errorf("TestMoonPhase4: expected %v to be %v", got, expected)
	}
}
