'use strict';

require("@babel/register");

const sass = require('gulp-sass')(require('sass'));
const { parallel, src, dest, watch, series } = require("gulp");
const babel = require("gulp-babel");
const load_plugins = require("gulp-load-plugins");
const plugins = load_plugins();

let pathto = function (file) {
  return "./public/" + file;
};
let scripts = {
  src: pathto("javascripts/*.js"),
  dest: pathto("dist"),
};
let styles = {
  src: pathto("stylesheets/*.scss"),
  scss: pathto("stylesheets/**/*.scss"),
  dest: pathto("dist"),
};

function javascript(cb) {
  src(scripts.src)
    .pipe(babel())
    .pipe(plugins.concat("app.js"))
    .pipe(plugins.uglify())
    .pipe(dest(scripts.dest));
  cb();
}

function css(cb) {
  src(styles.src)
    .pipe(sass())
    .pipe(plugins.csscomb())
    .pipe(plugins.cleanCss())
    .pipe(dest(styles.dest));
  cb();
}

exports.build = parallel(javascript, css);
exports.default = function () {
  // You can use a single task
  watch(styles.scss, css);
  // Or a composed task
  watch(scripts.src, series(javascript));
};
