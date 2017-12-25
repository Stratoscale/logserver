// vendor-bundles.webpack.config.js
var webpack = require('webpack')
var common  = require('./common')();

let plugins = [];

if (!process.env.DEV) {
  console.log('production build')
  plugins.push(new webpack.DefinePlugin({
    'process.env.NODE_ENV': JSON.stringify('production'),
  }))
  plugins.push(new webpack.optimize.UglifyJsPlugin())
}

plugins.push(new webpack.DllPlugin({
  // The path to the manifest file which maps between
  // modules included in a bundle and the internal IDs
  // within that bundle
  path: 'dist/[name]-manifest.json',
  // The name of the global variable which the library's
  // require function has been assigned to. This must match the
  // output.library option above
  name: '[name]',
}))

module.exports = {
  entry: {
    // create two library bundles, one with jQuery and
    // another with Angular and related libraries
    vendors: common.entry.vendors,
  },

  stats: 'verbose',

  output: {
    filename: '[name].bundle.js',
    path:     'dist/',

    // The name of the global variable which the library's
    // require() function will be assigned to
    library: '[name]',
  },

  plugins: plugins,
}

