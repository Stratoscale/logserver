var path              = require('path');
var _                 = require('lodash');
var webpack           = require('webpack');
var HtmlWebpackPlugin = require('html-webpack-plugin');
var AssetsPlugin      = require('assets-webpack-plugin');
var chalk             = require('chalk');
var common            = require('./common')();

var loaders = common.loaders;

process.traceDeprecation = true;

const devtool = process.env.SOURCE_MAP ? 'cheap-eval-source-map' : 'eval'
console.log('dev tools is:', devtool)

//dev specific loaders
loaders.push({
  test: /\.css$/,
  use:  ['style-loader', 'css-loader'],
});

loaders.push({
  test:    /\.js$/,
  include: [common.APP_PATH, common.LOADER_PATH, common.UI_LOADER_PATH, common.COMMON_PATH],
  exclude: /(node_modules)/,
  use:     [common.babelLoader],
});

loaders.push({
  test:    /\.sass$/,
  use:     common.sassLoaders,
  exclude: /(node_modules)/,
});

module.exports = {
  entry:         common.entry,
  output:        {
    path:              common.DIST_PATH,
    pathinfo:          true,
    publicPath:        common.outputUrl,
    filename:          '[name].js',
    sourceMapFilename: '[name].map'
  },
  module:        {
    rules: loaders,
  },
  stats:         'verbose',
  devtool,
  devServer:     common.devServer,
  resolveLoader: common.resolve.resolveLoaders,
  resolve:       common.resolve.resolve,
  plugins:       [
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NamedModulesPlugin(),
    new webpack.DefinePlugin(common.definePlugin),
    new HtmlWebpackPlugin(_.extend({}, common.htmlWebpackPlugin, {
      chunks:   ['app', 'progress'],
      basePath: process.env.BASEPATH || '',
    })),
    new AssetsPlugin({
      path:     path.resolve(common.ROOT_PATH, 'dist'),
      filename: 'manifest.json',
    }),
  ],
  externals:     {
    'jsdom': 'window',
  },
  node:          {
    fs:            'empty',
    child_process: 'empty',
    net:           'empty',
    tls:           'empty',
  },
};
