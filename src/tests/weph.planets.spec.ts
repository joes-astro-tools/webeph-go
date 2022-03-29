import moment from 'moment';
import { map } from 'rxjs/operators';
import { Geo, planets } from '../common';
import { EphLoader } from '../loaders/eph.loader';

export const withinToleranceMatcher: jasmine.CustomMatcherFactories = {
    toBeWithinTolerance: () => ({
        compare: (actual: number, expected: number, tolerance: number) => {
            const test = Math.abs(actual - expected) < tolerance;
            return {
                pass: test,
                message: test ? 'Within tolerance' : `Exceeds Tolerance: ${actual}, ${expected}, ${tolerance}`
            };
        }
    })
};

describe(
    'weph planet testing',
    () => {
        const tolerance = 1 / 60;
        let loader: EphLoader;
        const time = moment('2022-01-19T15:23:00-05:00');
        const coords: Geo = {
            // Latitude
            φ: 42,
            // Longitude, in degrees
            ο: -71.516667,
            // Height above mean sea level, in meters
            h: 56.0832
        };

        beforeEach(() => {
            loader = new EphLoader();

            jasmine.addMatchers(withinToleranceMatcher);
        });

        it(
            'saturn',
            (done: DoneFn) => {
                loader.resolve()
                    .pipe(
                        map(astroLib => astroLib.findLongitude(time, coords, planets.saturn))
                    )
                    .subscribe(
                        lon => {
                            const expected = 314.0401;
                            expect(lon.eclon).toBeWithinTolerance(expected, tolerance);
                            done();
                        }
                    );
            }
        );

        it(
            'jupiter',
            (done: DoneFn) => {
                loader.resolve()
                    .pipe(
                        map(astroLib => astroLib.findLongitude(time, coords, planets.jupiter))
                    )
                    .subscribe(
                        lon => {
                            const expected = 334.4736;
                            expect(lon.eclon).toBeWithinTolerance(expected, tolerance);
                            done();
                        }
                    );
            }
        );

        it(
            'mars',
            (done: DoneFn) => {
                loader.resolve()
                    .pipe(
                        map(astroLib => astroLib.findLongitude(time, coords, planets.mars))
                    )
                    .subscribe(
                        lon => {
                            const expected = 266.607;
                            expect(lon.eclon).toBeWithinTolerance(expected, tolerance);
                            done();
                        }
                    );
            }
        );

        it(
            'sun',
            (done: DoneFn) => {
                loader.resolve()
                    .pipe(
                        map(astroLib => astroLib.findLongitude(time, coords, planets.sun))
                    )
                    .subscribe(
                        lon => {
                            const expected = 299.7343;
                            expect(lon.eclon).toBeWithinTolerance(expected, tolerance);
                            done();
                        }
                    );
            }
        );

        it(
            'venus',
            (done: DoneFn) => {
                loader.resolve()
                    .pipe(
                        map(astroLib => astroLib.findLongitude(time, coords, planets.venus))
                    )
                    .subscribe(
                        lon => {
                            const expected = 282.945;
                            expect(lon.eclon).toBeWithinTolerance(expected, tolerance);
                            done();
                        }
                    );
            }
        );

        it(
            'mercury',
            (done: DoneFn) => {
                loader.resolve()
                    .pipe(
                        map(astroLib => astroLib.findLongitude(time, coords, planets.mercury))
                    )
                    .subscribe(
                        lon => {
                            const expected = 307.5887;
                            expect(lon.eclon).toBeWithinTolerance(expected, tolerance);
                            done();
                        }
                    );
            }
        );

        it(
            'moon',
            (done: DoneFn) => {
                loader.resolve()
                    .pipe(
                        map(astroLib => astroLib.findLongitude(time, coords, planets.moon))
                    )
                    .subscribe(
                        lon => {
                            const expected = 141.3216027;
                            expect(lon.eclon).toBeWithinTolerance(expected, tolerance);
                            done();
                        }
                    );
            }
        );

    });
