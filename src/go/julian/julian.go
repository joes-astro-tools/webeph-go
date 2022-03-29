// Julian: Chapter 7, Julian day.
//
// Under "General remarks", Meeus describes the INT function as used in the
// book.  In some contexts, math.Floor might be suitable, but I think more
// often, the functions base.FloorDiv or base.FloorDiv64 will be more
// appropriate.  See documentation in package base.
//
// On p. 63, Modified Julian Day is defined.  See constant JMod in package
// base.
//
// See also related functions JulianToGregorian and GregorianToJulian in
// package jm.
package julian

import (
	"math"
	base "webeph/base"
)

var (
	timeComponents = [6]int{}
)

// CalendarGregorianToJD converts a Gregorian year, month, and day of month
// to Julian day.
//
// Negative years are valid, back to JD 0.  The result is not valid for
// dates before JD 0.
//export calendarGregorianToJD
func CalendarGregorianToJD(y, m int, d float64) float64 {
	switch m {
	case 1, 2:
		y--
		m += 12
	}
	a := base.FloorDiv(y, 100)
	b := 2 - a + base.FloorDiv(a, 4)
	// (7.1) p. 61
	return float64(base.FloorDiv64(36525*(int64(y+4716)), 100)) +
		float64(base.FloorDiv(306*(m+1), 10)+b) + d - 1524.5
}

// CalendarJulianToJD converts a Julian year, month, and day of month to Julian day.
//
// Negative years are valid, back to JD 0.  The result is not valid for
// dates before JD 0.
func CalendarJulianToJD(y, m int, d float64) float64 {
	switch m {
	case 1, 2:
		y--
		m += 12
	}
	return float64(base.FloorDiv64(36525*(int64(y+4716)), 100)) +
		float64(base.FloorDiv(306*(m+1), 10)) + d - 1524.5
}

// Returns the container for calendar date values calculated by JDToCalendar.
// Receives:
//	nothing
// Returns:
//	The pointer in WASM linear memory to the container holding calculation results.
//export getTimeContainer
func GetTimeContainer() *[6]int {
	return &timeComponents
}

// Finds the calendar date for a given Julian day.
// Receives:
//	jd: a Julian day, as a float64
// Returns:
//	nothing
// Notes:
//	Stores results in a package-level container. Use GetTimeContainer to retrieve results.
//export jdToCalendar
func JDToCalendar(jd float64) {
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
	var year, month, day int
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
	floatDay := float64(int(b-d)-base.FloorDiv(306001*e, 1e4)) + f
	day = int(math.Floor(floatDay))
	floatHr := (floatDay - float64(day)) * 24
	floatMn := (floatHr - math.Floor(floatHr)) * 60
	floatSc := math.Round((floatMn - math.Floor(floatMn)) * 60)
	// Months produced by this algorithm fall between 1 and 12.
	// Moment.js needs them to fall between 0 and 11.
	timeComponents = [6]int{year, month - 1, day, int(floatHr), int(floatMn), int(floatSc)}
}
