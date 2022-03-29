import moment from 'moment';
import { firstValueFrom, tap } from 'rxjs';
import { AstroFns, Geo, planets } from '../common';
import { EphLoader } from '../loaders/eph.loader';

describe(
    'planet benchmarking',
    () => {
        const iterations = 10;
        const wasmSpeeds: Array<number> = [];
        const loader: EphLoader = new EphLoader();
        let astroLib: AstroFns;
        const time = moment('2022-01-19T15:23:00-05:00');
        const coords: Geo = {
            // Latitude
            φ: 42,
            // Longitude, in degrees
            ο: -71.516667,
            // Height above mean sea level, in meters
            h: 56.0832
        };


        function generatePlanetIteration(plNum: number, iteration: number): void {
            it(
                `venus ${iteration}`,
                () => {
                    const { perfMs } = astroLib.findLongitude(time, coords, plNum, true);
                    if (perfMs) {
                        wasmSpeeds.push(perfMs);
                    }
                    expect(perfMs).toBeDefined();
                }
            );
        }

        function generateVenusIterations(): void {
            for (let i = 0; i < iterations; i++) {
                generatePlanetIteration(planets.venus, i);
            }
        }

        beforeAll(async () => {
            jasmine.getEnv().addReporter({
                suiteDone: () => {
                    const max = Math.max(...wasmSpeeds);
                    const min = Math.min(...wasmSpeeds);
                    const ave = wasmSpeeds.reduce((a, b) => a + b) / wasmSpeeds.length;
                    console.log(`max: ${max}, min: ${min}, average: ${ave}`);
                }
            });

            await firstValueFrom(
                loader.resolve()
                    .pipe(tap(fns => astroLib = fns))
            );
        });

        generateVenusIterations();
    }
);
