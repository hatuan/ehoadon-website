var gulp = require('gulp'),
    del = require('del'),
    notifier = require('node-notifier'),
    child    = require('child_process'),
    concat = require('gulp-concat'),
    rename = require('gulp-rename'),
    uglify = require('gulp-uglify');

/////////////////////////////////////////////////////////////////////////////////////
//
// cleans the build output
//
/////////////////////////////////////////////////////////////////////////////////////
gulp.task('clean', function () {
    return del(['dist']);
});

/////////////////////////////////////////////////////////////////////////////////////
//
// Build a minified Javascript bundle - the order of the js files is determined
// by browserify
//
/////////////////////////////////////////////////////////////////////////////////////
gulp.task('build-js', function() {
    var jsFiles= [
            './scripts/**/*.js'
        ],
        jsDest = './dist/scripts';

    return gulp.src(jsFiles)
        .pipe(uglify())
        .pipe(gulp.dest(jsDest));
});

/////////////////////////////////////////////////////////////////////////////////////
//
// Copy Other
//
/////////////////////////////////////////////////////////////////////////////////////
gulp.task('copy', function () {
    return gulp.src(['assets/**/*', 'content/**/*', 'images/**/*', 'plugins/**/*', 'settings/**/*', '!settings/**/*.go', 'styles/**/*', 'templates/**/*', 'db/migrations/*', 'goose'], {
        base:"."
    })
    .pipe(gulp.dest('dist'));
});

/////////////////////////////////////////////////////////////////////////////////////
//
// full build (except sprites), applies cache busting to the main page css and js bundles
//
/////////////////////////////////////////////////////////////////////////////////////
gulp.task('build', ['build-js', 'copy'], function () {
});

/////////////////////////////////////////////////////////////////////////////////////
//
// GO Build
//
/////////////////////////////////////////////////////////////////////////////////////
gulp.task('server:build', function() {
    process.env['CGO_ENABLED'] = 0;
    process.env['GOOS'] = 'linux';

    var build =  child.exec('go build -ldflags "-w -s" -a -installsuffix cgo -o ./dist/ehoadon-website', function(err, stdout, stderr) {

        // Something wrong
        if (stderr.length) {
            util.log(util.colors.red('Something wrong with this version :'));
            var lines = stderr.toString()
                .split('\n').filter(function(line) {
                                return line.length
                            });
        
            for (var l in lines)
                util.log(util.colors.red(
                    'Error (go build): ' + lines[l]
                ));
                notifier.notify({
                    title: 'Error (go build)',
                    message: lines
                });
        }
    });
    
    return build;
});

/////////////////////////////////////////////////////////////////////////////////////
//
// installs and builds everything, including sprites
//
/////////////////////////////////////////////////////////////////////////////////////

gulp.task('default', ['clean'], function () {
    gulp.start(['build'], ['server:build']);
});