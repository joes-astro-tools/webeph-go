import * as fs from 'fs';
import * as path from 'path';
import { OptimizerSettings } from '../compile.helpers';

// Binaryen is an ES6 module. Grunt runs in NodeJS, which requires CommonJS.
// CommonJS can load ES6 modules via a dynamic import() expression.
// Typescript can't compile dynamic import()s for CommonJS yet.
// Using a workaround.
// Workaround from https://github.com/microsoft/TypeScript/issues/43329.
// They originally added support for the issue and closed it there, but the code was reverted in
// https://github.com/microsoft/TypeScript/issues/46452. The functionality isn't done yet.
const binaryenPromise = Function('return import("binaryen")')() as Promise<typeof import('binaryen')>;

module.exports = (grunt: IGrunt) => {
    grunt.registerTask(
        'optimize-wasm',
        () => {
            const config = grunt.config('optimizer') as OptimizerSettings;
            const done = grunt.task.current.async();
            // Use the workaround to load binaryen.
            binaryenPromise
                .then(({ default: binaryen }) => binaryen)
                .then(binaryen => {
                    // Load the WASM file. Although readFileSync claims to load a string, it really loads a Uint8Array.
                    // Trying to change the encoding explicitly using a TextEncoder breaks the encoding, so we need to just tell
                    // Typescript to relax.
                    const wasmFile = fs.readFileSync(config.fileLocation);
                    const binaryenModule = binaryen.readBinary(wasmFile as unknown as Uint8Array);
                    // Optimize it.
                    binaryen.setOptimizeLevel(3);
                    binaryenModule.optimize();
                    // Save.
                    const binaryOutput = binaryenModule.emitBinary();
                    fs.writeFileSync(config.fileLocation, binaryOutput);
                    // Emit WAT if necessary.
                    if (config.emitWAT) {
                        const cmpnts = path.parse(config.fileLocation);
                        const pathWAT = cmpnts.dir+'/'+cmpnts.name+'.wat';
                        fs.writeFileSync(pathWAT, binaryenModule.emitText());
                    }
                    // Tell Grunt that we are done.
                    done();
                })
                .catch(err => console.error(err));
        }
    );
};
