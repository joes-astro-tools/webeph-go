package schlyter_test

import (
	"testing"

	julian "webeph/julian"
	schlyter "webeph/schlyter"
	utils "webeph/testutils"
)

func TestFindDay(t *testing.T) {
	expected := float64(-3543)
	got := schlyter.FindDay(1990, 4, 19, 0)
	if expected != got {
		t.Errorf("FindDay for 04/19/1990@0000 = Expected %f to be %f", expected, got)
	}
}

func TestMercuryPositionλ(t *testing.T) {
	expected := 170.5709
	jde := julian.CalendarGregorianToJD(1990, 4, 19.)
	got, _, _ := schlyter.HeliocentricMercury(jde)
	if utils.CheckTolerance(got.Deg(), expected, 0.0001) == false {
		t.Errorf("TestMercuryPositionλ for 04/19/1990@0000 = Expected %f to be %f", expected, got)
	}
}

func TestMercuryPositionβ(t *testing.T) {
	expected := 5.9255
	jde := julian.CalendarGregorianToJD(1990, 4, 19.)
	_, got, _ := schlyter.HeliocentricMercury(jde)
	if utils.CheckTolerance(got.Deg(), expected, 0.0001) == false {
		t.Errorf("TestMercuryPositionλ for 04/19/1990@0000 = Expected %f to be %f", expected, got)
	}
}

func TestMercuryPositionΔ(t *testing.T) {
	expected := 0.374862
	jde := julian.CalendarGregorianToJD(1990, 4, 19.)
	_, _, got := schlyter.HeliocentricMercury(jde)
	if utils.CheckTolerance(got, expected, 0.000001) == false {
		t.Errorf("TestMercuryPositionλ for 04/19/1990@0000 = Expected %f to be %f", expected, got)
	}
}
