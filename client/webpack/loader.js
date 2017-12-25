var path              = require('path');
var _                 = require('lodash');
var webpack           = require('webpack');
var HtmlWebpackPlugin = require('html-webpack-plugin');
var ProgressBarPlugin = require('progress-bar-webpack-plugin');
var AssetsPlugin      = require('assets-webpack-plugin');

var common = require('./common')();

//dev specific loaders
var loaders = [{
  test:    /\.js$/,
  include: [common.UI_LOADER_PATH],
  loaders: common.babelLoader,
},
  {
    test:    /\.js$/, // include .js files
    exclude: /node_modules|vendor/, // exclude any and all files in the node_modules folder
    use:     [
      {
        loader:  'eslint-loader',
        options: {
          failOnError: true,
        }
      }
    ],
    enforce: 'pre',
  }];


module.exports = {
  entry:     {'ui-loader': common.UI_LOADER_PATH},
  output:    {
    path:       common.DIST_PATH + '/ui-loader',
    publicPath: common.outputUrl,
    filename:   '[hash].[name].js',
  },
  module:    {
    rules: loaders,
  },
  devtool:   'source-map',
  devServer: _.extend(common.devServer, {hot: false}),
  resolve:   common.resolve.resolve,
  plugins:   [
    new ProgressBarPlugin(),
    new webpack.DefinePlugin(common.definePlugin),
    new webpack.ContextReplacementPlugin(/moment[\/\\]locale$/, /en/),
    new webpack.optimize.UglifyJsPlugin({
      mangle: false,
    }),
    new HtmlWebpackPlugin(_.extend({}, common.htmlWebpackPlugin, {
      hash:   true,
      minify: {},
      chunks: ['ui-loader'],
    })),
    new AssetsPlugin({
      path:          path.resolve(common.ROOT_PATH, 'dist'),
      filename:      'loader-manifest.json',
      processOutput: function (assets) {
        var result = {
          file: '',
        };
        _.forEach(assets, function (asset, key) {
          if (key === 'ui-loader') {
            result.file = '//app.stratoscale.com/ui-loader' + asset.js;
          }
        });
        return JSON.stringify(result);
      },
    }),
  ],
};
