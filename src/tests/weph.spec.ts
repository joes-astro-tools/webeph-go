import moment, { Moment } from 'moment';
import { EphLoader } from '../loaders/eph.loader';
import { Geo } from '../common';


type MonthDays = 28 | 29 | 30 | 31;

function randomIntBetween(min: number, max: number): number {
    const minInt = Math.ceil(min);
    const maxInt = Math.floor(max);
    return Math.floor(Math.random() * (maxInt - minInt) + minInt);
}

function isLeapYear(y: number): boolean {
    return new Date(y, 1, 29).getMonth() === 1;
}

function findDaysInMonth(m: number, y: number): MonthDays {
    switch (m) {
        case 1:
            return 31;
        case 2:
            return isLeapYear(y) ? 29 : 28;
        case 3:
            return 31;
        case 4:
            return 30;
        case 5:
            return 31;
        case 6:
            return 30;
        case 7:
            return 31;
        case 8:
            return 31;
        case 9:
            return 30;
        case 10:
            return 31;
        case 11:
            return 30;
        default:
            return 31;
    }
}

function formatTime(tm: number): string {
    return tm < 10 ? `0${tm}` : tm.toString();
}

function generateRandomDay(): Moment {
    const y = randomIntBetween(2021, 2073);
    const m = randomIntBetween(1, 13);
    const maxDaysInMonth = findDaysInMonth(m, y);
    const dt = randomIntBetween(1, maxDaysInMonth + 1);
    const hr = randomIntBetween(1, 24);
    const mn = randomIntBetween(1, 60);
    const s = randomIntBetween(1, 60);
    return moment.utc(`${y}-${formatTime(m)}-${formatTime(dt)}T${formatTime(hr)}:${formatTime(mn)}:${formatTime(s)}`);
}

describe('weph testing', () => {
    let loader: EphLoader;

    beforeEach(() => { loader = new EphLoader(); });

    it('#resolve should resolve AstroFns', (done: DoneFn) => {
        loader.resolve()
            .subscribe(astroLib => {
                expect(astroLib).toBeDefined();
                done();
            });
    });

    it(
        'should retrieve twelve houses',
        (done: DoneFn) => {
            loader.resolve()
                .subscribe(astroLib => {
                    const coords: Geo = {
                        // Latitude
                        φ: 42,
                        // Longitude, in degrees
                        ο: -71.516667,
                        // Height above mean sea level, in meters
                        h: 56.0832
                    };
                    const mmt = generateRandomDay();
                    const jd = astroLib.findJD(mmt);
                    const [ε, lst] = astroLib.findObliquityLST(jd, coords);
                    const houses = astroLib.findHouses(lst, ε, coords);
                    houses.map(
                        house => expect(house).toBeGreaterThan(0)
                    );
                    done();
                });
        });

    it(
        'should retrieve sunrise and sunset',
        (done: DoneFn) => {
            loader.resolve()
                .subscribe(astroLib => {
                    const coords: Geo = {
                        // Latitude
                        φ: 42,
                        // Longitude, in degrees
                        ο: -71.516667,
                        // Height above mean sea level, in meters
                        h: 56.0832
                    };
                    const mmt = generateRandomDay();
                    const jd = astroLib.findJD(mmt);
                    const [sunrise, sunset] = astroLib.findSunRiseSet(jd, coords);
                    expect(sunrise).toBeGreaterThan(0);
                    expect(sunset).toBeGreaterThan(0);
                    done();
                });
        }
    );
    it(
        'should convert to a calendar date',
        (done: DoneFn) => {
            loader.resolve()
                .subscribe(astroLib => {
                    const mmt = generateRandomDay();
                    const jd = astroLib.findJD(mmt);
                    const processedMmt = astroLib.jdToMoment(jd, mmt.utcOffset());
                    expect(mmt.isSame(processedMmt, 'second')).toBeTrue();
                    done();
                });
        }
    );

});
