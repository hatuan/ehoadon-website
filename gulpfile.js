var gulp = require('gulp'),
    del = require('del'),
    concat = require('gulp-concat'),
    rename = require('gulp-rename'),
    uglify = require('gulp-uglify');

gulp.task('clean', function (cb) {
    del([
        'dist'
    ], cb);
});

gulp.task('build-js', function() {
    var jsFiles= [
            './scripts/**/*.js'
        ],
        jsDest = './dist/scripts';

    return gulp.src(jsFiles)
        .pipe(uglify())
        .pipe(gulp.dest(jsDest));
});

gulp.task('copy', function () {
    return gulp.src(['assets/**/*', 'content/**/*', 'images/**/*', 'plugins/**/*', 'settings/**/*', 'styles/**/*', 'templates/**/*'], {
        base:"."
    })
    .pipe(gulp.dest('dist'));
});

gulp.task('default', ['clean', 'build-js', 'copy']);