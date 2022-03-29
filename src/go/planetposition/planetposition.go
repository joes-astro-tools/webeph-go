// Planetposition: Chapter 32, Positions of the Planets.
//
// Incomplete:
//
// 1.  The package does not implement algorithms that use appendix III,
// but instead implements a full VSOP87 solution.  I do not have a copy
// of the supplimentary disk with appendix III in machine readable form
// and as the appendix is rather large, retyping it by hand is problematic.
// The full VSOP87 data set on the other hand is freely downloadable from
// the internet, so I implement here code that can use that data directly.
//
// 2.  The formula for accuracy of results is not implemented.  It is
// not needed for full VSOP87 solutions.
//
// 3.  Polynomial expressions are not implemented.  Again, implementation
// would involve typing rather large tables of numbers with associated
// risk of typographical errors.
package planetposition

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	base "webeph/base"
	coord "webeph/coord"
	precess "webeph/precess"
	unit "webeph/unit"
)

// Mercury-Neptune planet constants suitable for first argument to LoadPlanet.
const (
	Mercury = iota
	Venus
	Earth
	Mars
	Jupiter
	Saturn
	Sun
	Moon
	nPlanets // sad practicality
)

// parallel arrays, indexed by planet constants.
var (
	// extensions of VSOP87B files
	ext = [nPlanets]string{
		"mer", "ven", "ear", "mar", "jup", "sat", "ura", "nep"}

	// planet names as found in VSOP87B files
	b7 = [nPlanets]string{
		"MERCURY",
		"VENUS  ",
		"EARTH  ",
		"MARS   ",
		"JUPITER",
		"SATURN ",
		"URANUS ",
		"NEPTUNE",
	}
	mars = &V87Planet{
		L: [6][]Abc{
			// L0
			{
				{6.20347711581, 0., 0.},
				{0.18656368093, 5.0503710027, 3340.6124266998},
				{0.01108216816, 5.40099836344, 6681.2248533996},
				{0.00091798406, 5.75478744667, 10021.8372800994},
				{0.00027744987, 5.97049513147, 3.523118349},
				{0.00012315897, 0.84956094002, 2810.9214616052},
				{0.00010610235, 2.93958560338, 2281.2304965106},
			},
			// L1
			{
				{3340.61242700512, 0., 0.},
				{0.01457554523, 3.60433733236, 3340.6124266998},
			},
		},
		B: [6][]Abc{
			// B0
			{
				{0.03197134986, 3.76832042431, 3340.6124266998},
			},
		},
		R: [6][]Abc{
			// R0
			{
				{1.53033488271, 0., 0.},
				{0.1418495316, 3.47971283528, 3340.6124266998},
				{0.00660776362, 3.81783443019, 6681.2248533996},
				{0.00046179117, 4.15595316782, 10021.8372800994},
			},
		},
		Ibody: Mars,
	}
	jupiter = &V87Planet{
		L: [6][]Abc{
			// L0
			{
				{0.59954691494, 0., 0.},
				{0.09695898719, 5.06191793158, 529.6909650946},
				{0.00573610142, 1.44406205629, 7.1135470008},
				{0.00306389205, 5.41734730184, 1059.3819301892},
				{0.00097178296, 4.14264726552, 632.7837393132},
				{0.00072903078, 3.64042916389, 522.5774180938},
				{0.00064263975, 3.41145165351, 103.0927742186},
				{0.00039806064, 2.29376740788, 419.4846438752},
				{0.00038857767, 1.27231755835, 316.3918696566},
			},
			// L1
			{
				{529.69096508814, 0., 0.},
				{0.00489503243, 4.2208293947, 529.6909650946},
			},
		},
		B: [6][]Abc{
			// no terms needed for B
		},
		R: [6][]Abc{
			// R0
			{
				{5.20887429326, 0., 0.},
				{0.25209327119, 3.49108639871, 529.6909650946},
				{0.00610599976, 3.84115365948, 1059.3819301892},
			},
		},
		Ibody: Jupiter,
	}
	saturn = &V87Planet{
		L: [6][]Abc{
			// L0
			{
				{0.87401354025, 0., 0.},
				{0.11107659762, 3.96205090159, 213.299095438},
				{0.01414150957, 4.58581516874, 7.1135470008},
				{0.00398379389, 0.52112032699, 206.1855484372},
				{0.00350769243, 3.30329907896, 426.598190876},
				{0.00206816305, 0.24658372002, 103.0927742186},
				{0.000792713, 3.84007056878, 220.4126424388},
				{0.00023990355, 4.66976924553, 110.2063212194},
			},
			// L1
			{
				{213.2990952169, 0., 0.},
			},
		},
		B: [6][]Abc{
			// no terms needed for B
		},
		R: [6][]Abc{
			// R0
			{
				{9.55758135486, 0., 0.},
				{0.52921382865, 2.39226219573, 213.299095438},
			},
		},
		Ibody: Saturn,
	}
	earth = &V87Planet{
		L: [6][]Abc{
			// L0
			{
				{1.75347045673, 0., 0.},
				{0.03341656453, 4.66925680415, 6283.0758499914},
				{0.00034894275, 4.62610242189, 12566.1516999828},
				{0.00003497056, 2.74411783405, 5753.3848848968},
				{0.00003417572, 2.82886579754, 3.523118349},
				{0.00003135899, 3.62767041756, 77713.7714681205},
			},
			// L1
			{
				{6283.0758499914, 0., 0.},
				{0.00206058863, 2.67823455808, 6283.0758499914},
			},
		},
		B: [6][]Abc{
			// no terms needed for B
		},
		R: [6][]Abc{
			// R0
			{
				{1.00013988784, 0., 0.},
				{0.01670699632, 3.09846350258, 6283.0758499914},
				{0.00013956024, 3.05524609456, 12566.1516999828},
			},
		},
		Ibody: Earth,
	}
)

type Abc struct {
	a, b, c float64
}

type Coeff [6][]Abc

// V87Planet holds VSOP87 coefficients for computing planetary
// positions in spherical coorditates.
type V87Planet struct {
	L, B, R Coeff
	Ibody   int
}

// code tested with version 2.  other versions unknown.
const fileVersion = '2'

// LoadPlanet constructs a V87Planet object from a VSOP87 file.
//
// Argument ibody should be one of the planet constants.
//
// The directory containing the VSOP87 must be indicated by environment
// variable VSOP87.
func LoadPlanet(ibody int) (*V87Planet, error) {
	path := os.Getenv("VSOP87")
	if path == "" {
		// return nil, errors.New("No path assigned to environment variable VSOP87")
		path = "http://localhost:4200/assets/"
	}
	return LoadPlanetPath(ibody, path)
}

// LoadPlanetPath constructs a V87Planet object from a VSOP87 file.
//
// Argument ibody should be one of the planet constants; path should be
// a directory containing the VSOP87 files.
func LoadPlanetPath(ibody int, path string) (*V87Planet, error) {
	if ibody < 0 || ibody >= nPlanets {
		return nil, errors.New("Invalid planet.")
	}
	parsed := path + "VSOP87B." + ext[ibody]
	resp, err := http.Get(parsed)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	v := &V87Planet{}
	v.Ibody = ibody
	lines := strings.Split(string(data), "\n")
	n := 0
	n, err = v.L.Parse('1', ibody, lines, n, false)
	if err != nil {
		return nil, err
	}
	n, err = v.B.Parse('2', ibody, lines, n, false)
	if err != nil {
		return nil, err
	}
	_, err = v.R.Parse('3', ibody, lines, n, true)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c *Coeff) Parse(ic byte, ibody int, lines []string, n int, au bool) (int, error) {
	var cbuf [2047]Abc
	for n < len(lines) {
		line := lines[n]
		if len(line) < 132 {
			break
		}
		if line[41] != ic {
			break
		}
		if iv := line[17]; iv != fileVersion {
			return n, fmt.Errorf("Line %d: expected version %c, "+
				"found %c.", n+1, fileVersion, iv)
		}
		if bo := line[22:29]; bo != b7[ibody] {
			return n, fmt.Errorf("Line %d: expected body %s, "+
				"found %s.", n+1, b7[ibody], bo)
		}
		it := line[59] - '0'
		in, err := strconv.Atoi(strings.TrimSpace(line[60:67]))
		if err != nil {
			return n, fmt.Errorf("Line %d: %v.", n+1, err)
		}
		if in == 0 {
			continue
		}
		if in > len(lines)-n {
			return n, errors.New("Unexpected end of file.")
		}
		n++
		cx := 0
		for _, line := range lines[n : n+in] {
			a := &cbuf[cx]
			a.a, err =
				strconv.ParseFloat(strings.TrimSpace(line[79:97]), 64)
			if err != nil {
				goto parseError
			}
			a.b, err = strconv.ParseFloat(line[98:111], 64)
			if err != nil {
				goto parseError
			}
			a.c, err =
				strconv.ParseFloat(strings.TrimSpace(line[111:131]), 64)
			if err != nil {
				goto parseError
			}
			cx++
			continue
		parseError:
			return n, fmt.Errorf("Line %d: %v.", n+cx+1, err)
		}
		c[it] = append([]Abc{}, cbuf[:cx]...)
		n += in
	}
	return n, nil
}

// Position2000 returns ecliptic position of planets by full VSOP87 theory.
//
// Argument jde is the date for which positions are desired.
//
// Results are for the dynamical equinox and ecliptic J2000.
//
//	L is heliocentric longitude.
//	B is heliocentric latitude.
//	R is heliocentric range in AU.
func (vt *V87Planet) Position2000(jde float64) (L, B unit.Angle, R float64) {
	T := base.J2000Century(jde)
	τ := T * .1
	cf := make([]float64, 6)
	sum := func(series Coeff) float64 {
		for x, terms := range series {
			cf[x] = 0
			// sum terms in reverse order to preserve accuracy
			for y := len(terms) - 1; y >= 0; y-- {
				term := &terms[y]
				cf[x] += term.a * math.Cos(term.b+term.c*τ)
			}
		}
		return base.Horner(τ, cf[:len(series)]...)
	}
	L = unit.Angle(unit.PMod(sum(vt.L), 2*math.Pi))
	B = unit.Angle(sum(vt.B))
	R = sum(vt.R)
	return
}

// Position returns ecliptic position of planets at equinox and ecliptic of date.
//
// Argument jde is the date for which positions are desired.
//
// Results are positions consistent with those from Meeus's Apendix III,
// that is, at equinox and ecliptic of date.
//
//  L is heliocentric longitude.
//  B is heliocentric latitude.
//  R is heliocentric range in AU.
func (vt *V87Planet) Position(jde float64) (L, B unit.Angle, R float64) {
	L, B, R = vt.Position2000(jde)
	eclFrom := &coord.Ecliptic{
		Lat: B,
		Lon: L,
	}
	eclTo := &coord.Ecliptic{}
	epochFrom := 2000.0
	epochTo := base.JDEToJulianYear(jde)
	precess.EclipticPosition(eclFrom, eclTo, epochFrom, epochTo, 0, 0)
	return eclTo.Lon, eclTo.Lat, R
}

// ToFK5 converts ecliptic longitude and latitude from dynamical frame to FK5.
func ToFK5(L, B unit.Angle, jde float64) (L5, B5 unit.Angle) {
	// formula 32.3, p. 219.
	T := base.J2000Century(jde)
	Lp := L - unit.AngleFromDeg(1.397*T+.00031*T*T)
	sLp, cLp := Lp.Sincos()
	// (32.3) p. 219
	L5 = L + unit.AngleFromSec(-.09033+.03916*(cLp+sLp)*B.Tan())
	B5 = B + unit.AngleFromSec(.03916*(cLp-sLp))
	return
}

func GetMars() *V87Planet {
	return mars
}

func GetJupiter() *V87Planet {
	return jupiter
}

func GetSaturn() *V87Planet {
	return saturn
}

func GetEarth() *V87Planet {
	return earth
}
