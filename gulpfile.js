const gulp = require('gulp')
const sass = require('gulp-sass')
const exec = require('child_process').exec
const uglify = require('gulp-uglify')
const concat = require('gulp-concat')

sass.compiler = require('node-sass')

// node_modules js files to be compiled in static/build.js
// NOTE: if you're using `gulp watch` you MUST restart the process after changing anything here
const modules = [
  'jquery/dist/jquery.js',
  'bootstrap/js/dist/util.js',
  'bootstrap/js/dist/collapse.js',
]

for (j in modules) {
  modules[j] = `node_modules/${modules[j]}`
}

gulp.task('go:generate', (done) =>
  exec('env GOARCH=wasm GOOS=js go build -o ../static/main.wasm', {cwd: './frontend'}, (err, stdout, stderr) => {
    if (stdout) {
      console.log(stdout)
    }
    if (stderr) {
      console.error(stderr)
    }

    return done(err)
  })
)

gulp.task('go:watch', () =>
  gulp.watch('./frontend/**/*.go', gulp.series('go:generate'))
)

gulp.task('assets:js', () =>
  gulp.src(modules)
    .pipe(gulp.src('./assets/js/**/*.js'))
    .pipe(concat('build.js'))
    .pipe(uglify())
    .pipe(gulp.dest('./static/js'))
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
  gulp.watch('./assets/js/**/*.js', gulp.series('assets:js'))
})

gulp.task('default', gulp.series('go:generate', 'sass', 'assets:js'))
