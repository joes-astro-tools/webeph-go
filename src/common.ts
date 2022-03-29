import { Moment } from 'moment';

// Geographic coordinates
export interface Geo {
    // Latitude, in degrees
    φ: number;
    // Longitude, in degrees
    ο: number;
    // Height above mean sea level, in meters
    h: number
}

export type DayFn = (mmt: Moment) => number;
export type AngleConversionFn = (measure: number) => number;
export type LongFn = (mmt: Moment, coords: Geo, planet: number, measurePerf?: boolean) => LongitudeResult;

export const degsToRads: (degs: number) => number = (degs: number) => (Math.PI / 180) * degs;
export const radsToDegs: (degs: number) => number = (degs: number) => (180 / Math.PI) * degs;

export interface AstroFns {
    // Finds the equivalent TinyGo unit.Angle for a given degree.
    // Receives:
    //  measure: the angle, in degrees
    // Returns:
    //  the angle, as a TinyGo unit.Angle (always in radians)
    findAngleFromDeg: AngleConversionFn;

    // Finds the topocentric ecliptic longitude of a planet, in degrees.
    // Receives:
    //  mmt: the interesting Moment
    //  coords: coordinates of the observer, as a Geo interface
    //  planet: a number representing the planet, in Chaldean order
    //  measurePerf: measure performance, and include result in return.
    // Returns:
    //  LongitudeResult: an interface containing the ecliptic longitude, and possibly the performance numbers.
    findLongitude: LongFn;

    // Finds the equivalent Julian date for a given moment.
    // Receives:
    //  mmt: a Moment
    // Returns:
    //  the Julian date
    findJD: (mmt: Moment) => number;

    // Finds the ascending geocentric lunar node, in degrees.
    // Receives:
    //  jd: a Julian date
    // Returns:
    //  the ascending node, in degrees.
    // Notes:
    //  Find the opposite angle 180˚ away for the descending node.
    findAscendingNode: (jd: number) => number;

    // Finds the lunar phase, expressed as the difference between solar and lunar geocentric ecliptic longitude.
    // Receives:
    //	jd: the Julian day
    // Returns:
    //	the difference between solar and lunar longitude, in degrees. Rounds to the nearest degree.
    // Notes:
    //	Measures distance by subtracting solar longitude from lunar longitude. This has the effect of expressing waxing/waning:
    //	0: new
    //	0-179: waxing
    //	180: full
    //	181-359 waning
    findPhase: (jd: number) => number;

    // Find geocentric ecliptic longitude for a star, in degrees.
    // Receives:
    //	jd: Julian day
    //	ε: Obliquity, in radians
    //	raH: Right ascension hours
    //	raM: Right ascension minutes
    //	raS: Right ascension seconds
    //	declD: Declination degrees
    //	declM: Declination minutes
    //	declS: Declination seconds
    //	raμ: Change in right ascension, in arc seconds per year
    //	declμ: Change in declination, in arc seconds per year
    // Returns:
    //	geocentric ecliptic longitude, in degrees
    // Notes:
    //	1. The parameters are as granular as they are because star catalogs break it down this way. As we come closer to calculating
    //	 fixed star elections, we may change to decimal degrees for both measures, depending on what we can do with the star catalog.
    //	2. Geocentric and topocentric measure seem to be treated as the same for stellar measure. I find no explanation for this, but
    //   the difference between topocentric and geocentric measures is due to parallax, which is a function of distance. This makes a
    //   big difference for the Moon, and a small difference for planets. (2 degrees for the Moon, and only 8 seconds at maximum for
    //   planets.) Seeing that we are dealing with a fraction of a degree for the planets, and the stars are much farther away, the
    //   difference is likely so small that it is neglected.
    findStar: (
        jd: number,
        ε: number,
        raH: number,
        raM: number,
        raS: number,
        declD: number,
        declM: number,
        declS: number,
        raμ: number,
        declμ: number
    ) => number;

    // Finds sunrise and sunset for a given day.
    // Receives:
    //  jd: a Julian day
    //  coord: a Geo interface representing the observer's geographic coordinates
    // Returns:
    //  an array, where sunrise is provided first, then sunset
    findSunRiseSet: (jd: number, coord: Geo) => Array<number>;

    // Finds obliquity and local sidereal time, in radians.
    // Receives:
    //  jd: a Julian day
    //  coord: a Geo interface representing the observer's geographic coordinates
    // Returns:
    //  an array, where obliquity is provided first, then local sidereal time
    findObliquityLST: (jd: number, coord: Geo) => Array<number>;

    // Finds Regiomontanus houses.
    // Receives:
    //  lst: local sidereal time, in radians
    //  ε: obliquity, in radians
    //  coord: a Geo interface representing the observer's geographic coordinates
    // Returns:
    //  an array of the houses in order, measured in degrees.
    findHouses: (lst: number, ε: number, coord: Geo) => Array<number>;

    // Converts a Julian day to the equivalent Moment.
    // Receives:
    //  jd: Julian day
    //  offset: the offset from UTC, in minutes. Use mmt.utcOffset() to find this value.
    // Returns:
    //  the Moment found
    jdToMoment: (jd: number, offset: number) => Moment;
}

type PlanetNames = 'saturn' | 'jupiter' | 'mars' | 'sun' | 'venus' | 'mercury' | 'moon' | 'earth';

export const planets: { [key in PlanetNames]: number } = {
    saturn: 5,
    jupiter: 4,
    mars: 3,
    sun: 6,
    venus: 1,
    mercury: 0,
    moon: 7,
    earth: 2
};

export interface LongitudeResult {
    eclon: number;
    perfMs?: number;
}

