const gulp = require('gulp')
const sass = require('gulp-sass')

sass.compiler = require('node-sass')

gulp.task('sass', () =>
  gulp.src('./sass/**/*.scss')
    .pipe(sass().on('error', sass.logError))
    .pipe(gulp.dest('./css'))
)

gulp.task('sass:watch', () =>
  gulp.watch('./sass/**/*.scss', gulp.series('sass'))
)

gulp.task('default', gulp.series('sass'))
