const gulp = require('gulp')
const sass = require('gulp-sass')
const exec = require('child_process').exec

sass.compiler = require('node-sass')

gulp.task('go:generate', (done) =>
  exec('go generate', {cwd: './frontend'}, done)
)

gulp.task('go:watch', () =>
  gulp.watch('./frontend/**/*.go', gulp.series('go:generate'))
)

gulp.task('sass', () =>
  gulp.src('./assets/sass/**/*.scss')
    .pipe(sass().on('error', sass.logError))
    .pipe(gulp.dest('./static/css'))
)

gulp.task('sass:watch', () =>
  gulp.watch('./assets/sass/**/*.scss', gulp.series('sass'))
)

gulp.task('watch', () => {
  gulp.watch('./frontend/**/*.go', gulp.series('go:generate'))
  gulp.watch('./assets/sass/**/*.scss', gulp.series('sass'))
})

gulp.task('default', gulp.series('sass'))
