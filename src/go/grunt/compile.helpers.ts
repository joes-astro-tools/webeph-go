export interface CompilerSettings {
    cmd: string;
    source?: string;
    target?: string;
    args?: Array<string | [string, string]>;
}

export interface OptimizerSettings {
    fileLocation: string;
    emitWAT?: boolean;
}

export interface CompilerConfig extends grunt.config.IProjectConfig {
    tinygo: CompilerSettings;
    optimizer: OptimizerSettings;
}

// Concatenates a tuple together, adding '-' and '='.
// Used to build parameter strings for compilers.
// Example:
//  ['prop','value'] becomes '-prop=value'
function concatTuple(m: [string, string]): string {
    return `-${m[0]}=${m[1]}`;
}

// Converts a list of compiler arguments to a long string
// Example:
//  ['p',['bob','larry']] becomes '-p -bob=larry'
function argsToString(args: Array<string | [string, string]>): string {
    return args.map(arg => Array.isArray(arg) ? concatTuple(arg) : `-${arg}`).join(' ');
}

// Builds a compiler command string in the specific format required by TinyGo
// ie 'tinygo build -o [target] [args] [source]'
export function settingsToTinyGoCommand(settings: CompilerSettings, outputFlagReq = true): string {
    const cmpnts = [settings.cmd];
    if (settings.target) {
        cmpnts.push(outputFlagReq ? `-o ${settings.target}` : settings.target);
    }
    if (settings.args && settings.args.length > 0) {
        cmpnts.push(argsToString(settings.args));
    }
    if (settings.source) {
        cmpnts.push(settings.source);
    }
    return cmpnts.join(' ');
}
