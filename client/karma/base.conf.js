var webpackConfig = require('../webpack/test');

module.exports = function () {
  return {
    basePath: '../',
    frameworks: ['mocha', 'chai', 'sinon', 'intl-shim'],
    autoWatchBatchDelay: 100,
    files: [
      'node_modules/babel-polyfill/dist/polyfill.js',
      './app/components/**/**.spec.js',
      './loader/**.spec.js',
      './ui-loader/**.spec.js',
      '../common/**/**.spec.js',
    ],
    browsers: ['Chrome', 'Headless'],
    reporters: ['mocha'],
    preprocessors: {
      './app/components/**/**.spec.js': ['webpack'],
      './loader/**.spec.js': ['webpack'],
      './ui-loader/**.spec.js': ['webpack'],
      '../common/**/**.spec.js': ['webpack'],
    },
    customLaunchers: {
      Headless: {
        base:  'ChromeHeadless',
        flags: ['--no-sandbox'],
      },
    },
    webpack: webpackConfig,
    webpackMiddleware: {
      noInfo: true,
    },
  }
};