import moment, { Moment } from 'moment';
import { Injectable } from '@angular/core';
import { Resolve } from '@angular/router';
import { from, Observable, of } from 'rxjs';
import { map, switchMap, tap } from 'rxjs/operators';
import { AstroFns, AngleConversionFn, Geo, LongitudeResult } from '../common';
import { makeTinyGoImportObj, goRuntime } from '../tinygo';

const sizeOfFloat64 = 8;

interface TinyGoExport extends WebAssembly.Exports {
    memory: WebAssembly.Memory;
    findDay: (y: number, m: number, d: number, t: number) => number;
    angleFromDeg: AngleConversionFn;
    findDegFromAngle: AngleConversionFn;
    findLongitude: (y: number, m: number, t: number, φ: number, ο: number, h: number, planet: number) => number;
    _start: () => undefined;
    calendarGregorianToJD: (y: number, m: number, d: number) => number;
    findAscendingNode: (jd: number) => number;
    findMoonPhase: (jd: number) => number;
    findStellarLongitude: (
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
    getSunRiseSetPtr: () => number;
    findSunRiseSet: (jde: number, φ: number, ο: number) => void;
    getObliquityLSTContainer: () => number;
    findObliquityLST: (jd: number, ο: number) => void;
    getHouseContainer: () => number;
    findHouses: (lst: number, ε: number, φ: number) => void;
    getTimeContainer: () => number;
    jdToCalendar: (jd: number) => void;
}

@Injectable()
export class EphLoader implements Resolve<AstroFns> {
    private initialized = false;
    private memory: WebAssembly.Memory = new WebAssembly.Memory({ initial: 1 });

    resolve = (): Observable<AstroFns> => {
        if (this.initialized === false) {
            return from(fetch('assets/weph.wasm'))
                .pipe(
                    switchMap(wasmBytes => WebAssembly.instantiateStreaming(wasmBytes, makeTinyGoImportObj())),
                    tap(obj => goRuntime.run(obj.instance)),
                    map(obj => obj.instance.exports as TinyGoExport),
                    tap(exported => {
                        this.memory = exported.memory;
                        this.wasmFindAngleFromDeg = exported.angleFromDeg;
                        this.wasmFindLongitude = exported.findLongitude;
                        this.wasmCalendarGregorianToJD = exported.calendarGregorianToJD;
                        this.wasmFindAscendingNode = exported.findAscendingNode;
                        this.wasmFindMoonPhase = exported.findMoonPhase;
                        this.wasmFindStellarLongitude = exported.findStellarLongitude;
                        this.wasmGetSunRiseSetPtr = exported.getSunRiseSetPtr;
                        this.wasmFindSunRiseSet = exported.findSunRiseSet;
                        this.wasmGetObliquityLSTContainer = exported.getObliquityLSTContainer;
                        this.wasmFindObliquityLST = exported.findObliquityLST;
                        this.wasmGetHouseContainer = exported.getHouseContainer;
                        this.wasmFindHouses = exported.findHouses;
                        this.wasmGetTimeContainer = exported.getTimeContainer;
                        this.wasmJdToCalendar = exported.jdToCalendar;
                    }),
                    tap(() => this.initialized = true),
                    map(() => this.resolveLib())
                );
        }
        return of(this.resolveLib());
    };

    resolveLib(): AstroFns {
        return {
            findAngleFromDeg: this.findAngleFromDeg,
            findLongitude: this.findLongitude,
            findJD: this.findJD,
            findAscendingNode: this.findAscendingNode,
            findPhase: this.findPhase,
            findStar: this.findStar,
            findSunRiseSet: this.findSunRiseSet,
            findObliquityLST: this.findObliquityLST,
            findHouses: this.findHouses,
            jdToMoment: this.jdToMoment
        };
    }

    // Finds the equivalent TinyGo unit.Angle for a given degree.
    // Receives:
    //  measure: the angle, in degrees
    // Returns:
    //  the angle, as a TinyGo unit.Angle (always in radians)
    findAngleFromDeg = (measure: number): number => this.wasmFindAngleFromDeg(measure);

    // Finds the topocentric ecliptic longitude of a planet, in degrees.
    // Receives:
    //  mmt: the interesting Moment
    //  coords: coordinates of the observer, as a Geo interface
    //  planet: a number representing the planet, in Chaldean order
    //  measurePerf: measure performance, and include result in return.
    // Returns:
    //  LongitudeResult: an interface containing the ecliptic longitude, and possibly the performance numbers.
    findLongitude = (mmt: Moment, coords: Geo, planet: number, measurePerf = false): LongitudeResult => {
        const gt = mmt.utc();
        const y = gt.year();
        const m = gt.month() + 1;
        const t = gt.date() + (this.durationSinceMidnight(gt) / 24);
        const lat = this.wasmFindAngleFromDeg(coords.φ);
        const lon = this.wasmFindAngleFromDeg(coords.ο);
        if (measurePerf) {
            const begin = performance.now();
            const eclon = this.wasmFindLongitude(y, m, t, lat, lon, coords.h, planet);
            return {
                eclon,
                perfMs: performance.now() - begin
            };
        }
        else {
            return {
                eclon: this.wasmFindLongitude(y, m, t, lat, lon, coords.h, planet)
            };
        }
    };

    // Finds the equivalent Julian date for a given moment.
    // Receives:
    //  mmt: a Moment
    // Returns:
    //  the Julian date
    findJD = (mmt: Moment): number => {
        const gt = mmt.utc();
        const y = gt.year();
        const m = gt.month() + 1;
        const t = gt.date() + (this.durationSinceMidnight(gt) / 24);
        return this.wasmCalendarGregorianToJD(y, m, t);
    };

    // Finds the ascending geocentric lunar node, in degrees.
    // Receives:
    //  jd: a Julian date
    // Returns:
    //  the ascending node, in degrees.
    // Notes:
    //  Find the opposite angle 180˚ away for the descending node.
    findAscendingNode = (jd: number): number => this.wasmFindAscendingNode(jd);

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
    findPhase = (jd: number): number => this.wasmFindMoonPhase(jd);

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
    findStar = (jd: number,
        ε: number,
        raH: number,
        raM: number,
        raS: number,
        declD: number,
        declM: number,
        declS: number,
        raμ: number,
        declμ: number): number => this.wasmFindStellarLongitude(jd, ε, raH, raM, raS, declD, declM, declS, raμ, declμ);

    // Finds sunrise and sunset for a given day.
    // Receives:
    //  jd: a Julian day
    //  coord: a Geo interface representing the observer's geographic coordinates
    // Returns:
    //  an array, where sunrise is provided first, then sunset
    findSunRiseSet = (jd: number, coord: Geo): Array<number> => {
        const φ = this.wasmFindAngleFromDeg(coord.φ);
        const ο = this.wasmFindAngleFromDeg(coord.ο);
        this.wasmFindSunRiseSet(jd, φ, ο);
        const begin = this.wasmGetSunRiseSetPtr();
        const end = begin + (sizeOfFloat64 * 2);
        const memView = new Float64Array(this.memory.buffer.slice(begin, end));
        return Array.from(memView);
    };

    // Finds obliquity and local sidereal time, in radians.
    // Receives:
    //  jd: a Julian day
    //  coord: a Geo interface representing the observer's geographic coordinates
    // Returns:
    //  an array, where obliquity is provided first, then local sidereal time
    findObliquityLST = (jd: number, coord: Geo): Array<number> => {
        const ο = this.wasmFindAngleFromDeg(coord.ο);
        this.wasmFindObliquityLST(jd, ο);
        const begin = this.wasmGetObliquityLSTContainer();
        // We need to retrieve two numbers from linear memory: (1) obliquity and (2) lst.
        // End address is two 64-bit values from the beginning.
        const end = begin + (sizeOfFloat64 * 2);
        const memView = new Float64Array(this.memory.buffer.slice(begin, end));
        return Array.from(memView);
    };

    // Finds Regiomontanus houses.
    // Receives:
    //  lst: local sidereal time, in radians
    //  ε: obliquity, in radians
    //  coord: a Geo interface representing the observer's geographic coordinates
    // Returns:
    //  an array of the houses in order, measured in degrees.
    findHouses = (lst: number, ε: number, coord: Geo): Array<number> => {
        const φ = this.wasmFindAngleFromDeg(coord.φ);
        this.wasmFindHouses(lst, ε, φ);
        const begin = this.wasmGetHouseContainer();
        const end = begin + (sizeOfFloat64 * 12);
        const memView = new Float64Array(this.memory.buffer.slice(begin, end));
        return Array.from(memView);
    };

    // Converts a Julian day to the equivalent Moment.
    // Receives:
    //  jd: Julian day
    //  offset: the offset from UTC, in minutes. Use mmt.utcOffset() to find this value.
    // Returns:
    //  the Moment found
    jdToMoment = (jd: number, offset: number): Moment => {
        const sizeOfInt32 = 4;
        const begin = this.wasmGetTimeContainer();
        this.wasmJdToCalendar(jd);
        const end = begin + (sizeOfInt32 * 6);
        const memView = new Int32Array(this.memory.buffer.slice(begin, end));
        const [y, m, d, hr, mn, sc] = Array.from(memView);
        let mmt: Moment;
        // Very rarely, sc will be precisely 60. Moment can't handle that.
        if (sc === 60) {
            // Create the time, omitting seconds, then add the seconds to the existing time.
            mmt = moment.utc([y, m, d, hr, mn]);
            mmt.add(sc, 'seconds');
        } else {
            mmt = moment.utc([y, m, d, hr, mn, sc]);
        }
        mmt.utcOffset(offset);
        return mmt;
    };

    fractionalDay(time: Moment): number {
        return time.date() + (this.durationSinceMidnight(time) / 24.0);
    }

    durationSinceMidnight(time: Moment): number {
        const midnight = time.clone().startOf('day');
        return time.diff(midnight, 'hours', true);
    }

    private wasmFindAngleFromDeg: AngleConversionFn = () => 0;
    private wasmFindLongitude: (y: number, m: number, t: number, φ: number, ο: number, h: number, planet: number) => number = () => 0;
    private wasmCalendarGregorianToJD: (y: number, m: number, d: number) => number = () => 0;
    private wasmFindAscendingNode: (jd: number) => number = () => 0;
    private wasmFindMoonPhase: (jd: number) => number = () => 0;
    private wasmFindStellarLongitude: (
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
    ) => number = () => 0;
    private wasmGetSunRiseSetPtr: () => number = () => 0;
    private wasmFindSunRiseSet: (jd: number, φ: number, ο: number) => void = () => 0;
    private wasmGetObliquityLSTContainer: () => number = () => 0;
    private wasmFindObliquityLST: (jd: number, ο: number) => void = () => 0;
    private wasmGetHouseContainer: () => number = () => 0;
    private wasmFindHouses: (lst: number, ε: number, φ: number) => void = () => 0;
    private wasmGetTimeContainer: () => number = () => 0;
    private wasmJdToCalendar: (jd: number) => void = () => 0;
}
