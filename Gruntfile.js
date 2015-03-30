module.exports = function (grunt) {
    require('jit-grunt')(grunt, {
      shell: 'grunt-shell-spawn',
    });
    // require('time-grunt')(grunt);

    var gruntConfig = {
        shell: {
            go_install: {
                command: 'go install ./...',
            },
            go_test: {
                command: './test.sh',
            },
            options: {
                failOnError: true
            }
        },
        watch: {
            api: {
                files: ['**/*.go', 'test.sh'],
                tasks: ['clear', 'build_go_code']
            }
        },
    }; 

    grunt.registerTask('build_go_code', ['shell:go_install', 'clear', 'shell:go_test']);
    grunt.registerTask('default', 'watch');
    grunt.initConfig(gruntConfig);
};
