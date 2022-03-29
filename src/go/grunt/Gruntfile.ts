import { CompilerConfig } from './compile.helpers';

const fileLocation = '../../../assets/weph.wasm';

// All the settings for the various commands live here.
const compilerConfig: CompilerConfig = {
    tinygo: {
        cmd: 'tinygo build',
        source: '../../eph.go',
        target: fileLocation,
        args: [
            'no-debug',
            ['scheduler', 'none'],
            ['panic', 'trap']
        ]
    },
    optimizer: {
        fileLocation
    }
};

module.exports = (grunt: IGrunt) => {
    // Set up the configuration for each command's parsing subtask.
    grunt.initConfig(compilerConfig);
    // Load all of the subtasks in the 'grunt-tasks' folder.
    // (One or more tasks will configure commands for grunt-exec.)
    grunt.loadTasks('grunt-tasks');
    // Load grunt-exec: this will execute commands from the command line using the configuration.
    grunt.loadNpmTasks('grunt-exec');
    // Run tasks in this order:
    //  parse.tinygo: configure a command for 'exec' to use the TinyGo compiler to compile a WASM file.
    //  exec: run the above command
    //  optimize-wasm: use Binaryen to optimize the WASM file.
    grunt.registerTask(
        'default',
        [
            'parse.tinygo',
            'exec',
            'optimize-wasm'
        ]
    );
};
