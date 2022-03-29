import { CompilerSettings, settingsToTinyGoCommand } from '../compile.helpers';

module.exports = (grunt: IGrunt) => {
    // Get the TinyGo configuration.
    const config = grunt.config('tinygo') as CompilerSettings;
    // Register a task for Gruntfile.
    grunt.registerTask(
        'parse.tinygo',
        () => {
            // When Gruntfile runs this task:
            // Add instructions for 'grunt-exec' to the overall configuration.
            grunt.config.merge({
                exec: {
                    // Build the command for 'grunt-exec'.
                    'compile.tinygo': {
                        cmd: settingsToTinyGoCommand(config),
                        stdout: false
                    }
                }
            });
        }
    );
};
