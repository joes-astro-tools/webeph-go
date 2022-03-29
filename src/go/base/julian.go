package base

// Julian and Besselian years described in chapter 21, Precession.
// T, Julian centuries since J2000 described in chapter 22, Nutation.

// JulianYear and other common periods.
const (
	J2000         = 2451545.0
	JulianYear    = 365.25 // days
	JulianCentury = 36525  // days
)

// JulianYearToJDE returns the Julian ephemeris day for a Julian year.
func JulianYearToJDE(jy float64) float64 {
	return J2000 + JulianYear*(jy-2000)
}

// JDEToJulianYear returns a Julian year for a Julian ephemeris day.
func JDEToJulianYear(jde float64) float64 {
	return 2000 + (jde-J2000)/JulianYear
}

// J2000Century returns the number of Julian centuries since J2000.
//
// The quantity appears as T in a number of time series.
func J2000Century(jde float64) float64 {
	// The formula is given in a number of places in the book, for example
	// (12.1) p. 87.
	// (22.1) p. 143.
	// (25.1) p. 163.
	return (jde - J2000) / JulianCentury
}
