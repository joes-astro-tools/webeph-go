// naming-convention shortcoming here:
// https://github.com/typescript-eslint/typescript-eslint/issues/2244
/* eslint-disable @typescript-eslint/naming-convention */

interface GoImports extends WebAssembly.Imports {
    // Apparently various functions to start supporting WASI.
    // WASM: WebAssembly for browsers
    // WASI: WebAssembly outside browsers (somewhere in an OS)
    // https://wasi.dev/
    // (Nothing here related to WASM, so all these are irrelevant to this project.)
    wasi_snapshot_preview1: {
        fd_write: (fd: number, iovs_ptr: number, iovs_len: number, nwritten_ptr: number) => number;
        // The following three functions are labelled 'dummy', and seem to be stubs.
        fd_close: () => 0;
        fd_fdstat_get: () => 0;
        fd_seek: () => 0;
        proc_exit: (code: number) => never;
        random_get: (bufPtr: number, bufLen: number) => number;
    };
    env: {
        // main.main and runtime.alloc do not appear in wasm_exec.js.
        // They must be instantiated at runtime.
        // https://github.com/tinygo-org/tinygo/issues/2495
        'main.main': (a: number, b: number) => void;
        'runtime.alloc': (a: number, b: number, c: number, d: number) => number;
        'time.stopTimer': (a: number, b: number, c: number) => number;
        'time.resetTimer': (a: number, b: number, c: number, d: number) => number;
        'time.startTimer': (a: number, b: number, c: number) => number;
        'sync/atomic.AddInt32': (a: number, b: number, c: number, d: number) => number;
        'runtime.ticks': () => number;
        'runtime.sleepTicks': (timeout: number) => NodeJS.Timeout;
        // https://pkg.go.dev/syscall/js@go1.17.6
        // Manages Javascript access to WebAssembly memory for WASM generated by Go/TinyGo
        // Still in Beta apparently
        //
        // TinyGo does not support finalizers: below function should never be called.
        'syscall/js.finalizeRef': () => void;
        'syscall/js.stringVal': (ret_ptr: number, value_ptr: number, value_len: number) => void;
        'syscall/js.valueGet': (retval: number, v_addr: number, p_ptr: number, p_len: number) => void;
        'syscall/js.valueSet': (v_addr: number, p_ptr: number, p_len: number, x_addr: number) => boolean;
        'syscall/js.valueDelete': (v_addr: number, p_ptr: number, p_len: number) => boolean;
        'syscall/js.valueIndex': (ret_addr: number, v_addr: number, i: number) => boolean;
        'syscall/js.valueSetIndex': (v_addr: number, i: number, x_addr: number) => boolean;
        'syscall/js.valueCall': (
            ret_addr: number, v_addr: number, m_ptr: number, m_len: number,
            args_ptr: number, args_len: number, args_cap: number
        ) => void;
        'syscall/js.valueInvoke': (ret_addr: number, v_addr: number, args_ptr: number, args_len: number, args_cap: number) => void;
        'syscall/js.valueNew': (ret_addr: number, v_addr: number, args_ptr: number, args_len: number, args_cap: number) => void;
        'syscall/js.valueLength': (v_addr: number) => number;
        'syscall/js.valuePrepareString': (ret_addr: number, v_addr: number) => void;
        'syscall/js.valueLoadString': (v_addr: number, slice_ptr: number, slice_len: number, slice_cap: number) => void;
        'syscall/js.valueInstanceOf': (v_addr: number, t_addr: number) => boolean;
        'syscall/js.copyBytesToGo': (ret_addr: number, dest_addr: number, dest_len: number, dest_cap: number, source_addr: number) => void;
        'syscall/js.copyBytesToJS': (
            ret_addr: number, dest_addr: number, source_addr: number,
            source_len: number, source_cap: number
        ) => void;
    };
}

declare class Go {
    importObject: GoImports;
    mem: () => DataView;
    setInt64: (addr: number, v: number) => void;
    getInt64: (addr: number) => number;
    loadValue: (addr: number) => number | undefined;
    storeValue: (addr: number) => void;
    loadSlice: (array: number, len: number) => Uint8Array;
    loadSliceOfValues: (array: number, len: number) => Array<number | undefined>;
    loadString: (ptr: number, len: number) => string;
    timeOrigin: number;
    run(instance: WebAssembly.Instance): Promise<void>;
}
