package julian_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	julian "webeph/julian"
)

func ExampleCalendarGregorianToJD_sputnik() {
	// Example 7.a, p. 61.
	jd := julian.CalendarGregorianToJD(1957, 10, 4.81)
	fmt.Printf("%.2f\n", jd)
	// Output:
	// 2436116.31
}

func ExampleCalendarGregorianToJD_halley() {
	// Example 7.c, p. 64.
	jd1 := julian.CalendarGregorianToJD(1910, 4, 20)
	jd2 := julian.CalendarGregorianToJD(1986, 2, 9)
	fmt.Printf("%.0f days\n", jd2-jd1)
	// Output:
	// 27689 days
}

func TestGreg(t *testing.T) {
	for _, tp := range []struct {
		y, m  int
		d, jd float64
	}{
		{2000, 1, 1.5, 2451545}, // more examples, p. 62
		{1999, 1, 1, 2451179.5},
		{1987, 1, 27, 2446822.5},
		{1987, 6, 19.5, 2446966},
		{1988, 1, 27, 2447187.5},
		{1988, 6, 19.5, 2447332},
		{1900, 1, 1, 2415020.5},
		{1600, 1, 1, 2305447.5},
		{1600, 12, 31, 2305812.5},
	} {
		dt := julian.CalendarGregorianToJD(tp.y, tp.m, tp.d) - tp.jd
		if math.Abs(dt) > .1 {
			t.Logf("%#v", tp)
			t.Fatal("dt:", time.Duration(dt*24*float64(time.Hour)))
		}
	}
}

func ExampleCalendarJulianToJD() {
	// Example 7.b, p. 61.
	jd := julian.CalendarJulianToJD(333, 1, 27.5)
	fmt.Printf("%.1f\n", jd)
	// Output:
	// 1842713.0
}

func ExampleJDToCalendar() {
	// Example 7.c, p. 64.
	cmpnts := julian.GetTimeContainer()
	julian.JDToCalendar(2436116.31)
	y, m, d, hr, mn, sc := cmpnts[0], cmpnts[1], cmpnts[2], cmpnts[3], cmpnts[4], cmpnts[5]
	m += 1 // Convert from Moment format to TinyGo format
	fmt.Printf("%d %s %d %d %d %d\n", y, time.Month(m), d, hr, mn, sc)
	// Output:
	// 1957 October 4 19 26 24
}

func TestYMD(t *testing.T) {
	for _, tp := range []struct {
		jd                  float64
		y, m, d, hr, mn, sc int
	}{
		{1842713, 333, 1, 27, 12, 0, 0},
		{1507900.13, -584, 5, 28, 15, 7, 12},
	} {
		cmpnts := julian.GetTimeContainer()
		julian.JDToCalendar(tp.jd)
		y, m, d, hr, mn, sc := cmpnts[0], cmpnts[1], cmpnts[2], cmpnts[3], cmpnts[4], cmpnts[5]
		m += 1 // Convert from Moment format to TinyGo format
		if y != tp.y || m != tp.m || d != tp.d || hr != tp.hr || mn != tp.mn || sc != tp.sc {
			t.Logf("%#v", tp)
			t.Fatal("JDToYMD", y, m, d, hr, mn, sc)
		}
	}
}
