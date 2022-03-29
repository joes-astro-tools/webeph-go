/* eslint-disable @typescript-eslint/no-unused-vars */

export const goRuntime = new Go();

export const makeTinyGoImportObj = (): GoImports => {
    // Go class instantiated in wasm_exec.js
    // wasm_exec.js executed during bootstrap: setting in angular.json
    const importObj = goRuntime.importObject;
    // main.main and runtime.alloc do not appear in wasm_exec.js.
    // They must be instantiated at runtime.
    // https://github.com/tinygo-org/tinygo/issues/2495
    importObj.env['runtime.alloc'] = (a: number, b: number, c: number, d: number) => 0;
    importObj.env['main.main'] = (a: number, b: number) => void 0;
    // Timers are not implemented yet.
    // https://github.com/tinygo-org/tinygo/issues/1037
    importObj.env['time.stopTimer'] = (a: number, b: number, c: number) => 0;
    importObj.env['time.resetTimer'] = (a: number, b: number, c: number, d: number) => 0;
    importObj.env['time.startTimer'] = (a: number, b: number, c: number) => 0;
    // Don't find any reference to this one. Probably need to submit an issue for it.
    importObj.env['sync/atomic.AddInt32'] = (a: number, b: number, c: number, d: number) => 0;
    return importObj;
};
