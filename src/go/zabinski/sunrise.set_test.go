package zabinski_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	base "webeph/base"
	julian "webeph/julian"
	testutils "webeph/testutils"
	unit "webeph/unit"
	zabinski "webeph/zabinski"
)

// JDToCalendar returns the calendar date for the given jd.
//
// Note that this function returns a date in either the Julian or Gregorian
// Calendar, as appropriate.
func JDToCalendar(jd float64) (year, month int, day float64) {
	zf, f := math.Modf(jd + .5)
	z := int64(zf)
	a := z
	if z >= 2299151 {
		α := base.FloorDiv64(z*100-186721625, 3652425)
		a = z + 1 + α - base.FloorDiv64(α, 4)
	}
	b := a + 1524
	c := base.FloorDiv64(b*100-12210, 36525)
	d := base.FloorDiv64(36525*c, 100)
	e := int(base.FloorDiv64((b-d)*1e4, 306001))
	// compute return values
	day = float64(int(b-d)-base.FloorDiv(306001*e, 1e4)) + f
	switch e {
	default:
		month = e - 1
	case 14, 15:
		month = e - 13
	}
	switch month {
	default:
		year = int(c) - 4716
	case 1, 2:
		year = int(c) - 4715
	}
	return
}

// Displays a Julian day in readable format.
// Receives:
//	jde: a Julian day
//	tz: the time zone offset
// Returns:
//	a string containing a description of the date
// Notes:
//	Be sure to set tz to positive/negative depending on whether the time zone is positive/negative.
func JDToDateString(jde float64, tz int) string {
	// Example date of 2/10/22@14:00, Eastern Standard Time. JD: 2459621.291666667
	// Get year, month, fractionalD.
	y, m, fractionalD := JDToCalendar(jde)
	// Example: y = 2022, m = 2, fractionalD = 10.791666667
	d := math.Floor(fractionalD)
	// Example: d = 10
	jdTime := fractionalD - math.Floor(fractionalD)
	// Example: jdTime = 0.791666667
	asHours := (jdTime * 24) + float64(tz)
	// Example: asHours = 14.[digits]
	h := math.Floor(asHours)
	// Example: h = 14
	fracMins := (asHours - h) * 60
	// Example: fracMins = 0.[digits]
	mins := math.Floor(fracMins)
	// Example: mins = 0
	secs := math.Floor((fracMins - mins) * 60)
	// Example: secs = 0
	return fmt.Sprintf("%v %v, %v@%v:%v:%v", time.Month(m), d, y, h, mins, secs)
}

// Converts hours, minutes, seconds, timezone to a fractional day.
// Receives:
//	d: the day
//	h: the hour
//	m: the minute
//	s: the second
//	tz: the timezone
// Returns:
//	the fractional day
// Notes:
//	Be sure to set tz to positive/negative depending on whether the time zone is positive/negative.
func HmsTzToFractionalDay(d, h, m, s, tz int) float64 {
	return float64(d) + ((float64(h-tz) / 24) + (float64(m) / 1440) + (float64(s) / 86400))
}

func TestFindSunRiseSet1(t *testing.T) {
	jde := julian.CalendarGregorianToJD(2022, 2, 11)
	φ := unit.AngleFromDeg(42.0028761)
	ο := unit.AngleFromDeg(-71.5147839)
	ptr := zabinski.GetSunRiseSetPtr()
	zabinski.FindSunRiseSet(jde, φ, ο)
	expected := julian.CalendarGregorianToJD(2022, 2, HmsTzToFractionalDay(11, 6, 45, 22, -5))
	// We have no working examples of the spec used.
	// Instead, we have approximations of the spec taken from three locations:
	// https://sunrise-sunset.org/
	// https://gml.noaa.gov/grad/solcalc/
	// A working spreadsheet showing an alternate form of finding sunrise/sunset, taken from https://gml.noaa.gov/grad/solcalc/calcdetails.html
	// None of these examples precisely agree, and together form an ambiguity of about 3 minutes.
	if !testutils.CheckTolerance(ptr[0], expected, testutils.JulianMinuteTolerance*3) {
		t.Errorf("TestFindSunRiseSet1: expected %v to be %v", JDToDateString(ptr[0], -5), JDToDateString(expected, -5))
	}
	// 2/11/22@17:15:8, tz -5 for Eastern Standard Time
	expected = julian.CalendarGregorianToJD(2022, 2, HmsTzToFractionalDay(11, 17, 15, 8, -5))
	if !testutils.CheckTolerance(ptr[1], expected, testutils.JulianMinuteTolerance*3) {
		t.Errorf("TestFindSunRiseSet1: expected %v to be %v", JDToDateString(ptr[1], -5), JDToDateString(expected, -5))
	}
}

func TestFindSunRiseSet2(t *testing.T) {
	jde := julian.CalendarGregorianToJD(2022, 2, 11)
	φ := unit.AngleFromDeg(-42.00287)
	ο := unit.AngleFromDeg(-71.514784)
	ptr := zabinski.GetSunRiseSetPtr()
	zabinski.FindSunRiseSet(jde, φ, ο)
	expected := julian.CalendarGregorianToJD(2022, 2, HmsTzToFractionalDay(11, 7, 4, 0, -3))
	// We have no working examples of the spec used.
	// Instead, we have approximations of the spec taken from three locations:
	// https://sunrise-sunset.org/
	// https://gml.noaa.gov/grad/solcalc/
	// A working spreadsheet showing an alternate form of finding sunrise/sunset, taken from https://gml.noaa.gov/grad/solcalc/calcdetails.html
	// None of these examples precisely agree, and together form an ambiguity of about 3 minutes.
	if !testutils.CheckTolerance(ptr[0], expected, testutils.JulianMinuteTolerance*3) {
		t.Errorf("TestFindSunRiseSet2: expected %v to be %v", JDToDateString(ptr[0], -3), JDToDateString(expected, -3))
	}
	// 2/11/22@17:15:8, tz -5 for Eastern Standard Time
	expected = julian.CalendarGregorianToJD(2022, 2, HmsTzToFractionalDay(11, 20, 56, 0, -3))
	if !testutils.CheckTolerance(ptr[1], expected, testutils.JulianMinuteTolerance*3) {
		t.Errorf("TestFindSunRiseSet2: expected %v to be %v", JDToDateString(ptr[1], -3), JDToDateString(expected, -3))
	}
}

func TestFindSunRiseSet3(t *testing.T) {
	jde := julian.CalendarGregorianToJD(2022, 2, 11)
	φ := unit.AngleFromDeg(-42.00287)
	ο := unit.AngleFromDeg(71.514784)
	ptr := zabinski.GetSunRiseSetPtr()
	zabinski.FindSunRiseSet(jde, φ, ο)
	expected := julian.CalendarGregorianToJD(2022, 2, HmsTzToFractionalDay(11, 5, 31, 0, 5))
	// We have no working examples of the spec used.
	// Instead, we have approximations of the spec taken from three locations:
	// https://sunrise-sunset.org/
	// https://gml.noaa.gov/grad/solcalc/
	// A working spreadsheet showing an alternate form of finding sunrise/sunset, taken from https://gml.noaa.gov/grad/solcalc/calcdetails.html
	// None of these examples precisely agree, and together form an ambiguity of about 3 minutes.
	if !testutils.CheckTolerance(ptr[0], expected, testutils.JulianMinuteTolerance*3) {
		t.Errorf("TestFindSunRiseSet3: expected %v to be %v", JDToDateString(ptr[0], 5), JDToDateString(expected, 5))
	}
	// 2/11/22@17:15:8, tz -5 for Eastern Standard Time
	expected = julian.CalendarGregorianToJD(2022, 2, HmsTzToFractionalDay(11, 19, 24, 0, 5))
	if !testutils.CheckTolerance(ptr[1], expected, testutils.JulianMinuteTolerance*3) {
		t.Errorf("TestFindSunRiseSet3: expected %v to be %v", JDToDateString(ptr[1], 5), JDToDateString(expected, 5))
	}
}
