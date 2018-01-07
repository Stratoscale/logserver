const path                    = require('path');
const _                       = require('lodash');
const webpack                 = require('webpack');
const ExtractTextPlugin       = require('extract-text-webpack-plugin');
const HtmlWebpackPlugin       = require('html-webpack-plugin');
const ProgressBarPlugin       = require('progress-bar-webpack-plugin');
const OptimizeCssAssetsPlugin = require('optimize-css-assets-webpack-plugin');
const common                  = require('./common')();

const loaders = common.loaders;

//dev specific loaders
loaders.push({
  test:    /\.js$/,
  include: [common.APP_PATH, common.LOADER_PATH, common.UI_LOADER_PATH, common.COMMON_PATH],
  exclude: /(node_modules)/,
  use:     [common.babelLoader],
});

loaders.push({
  test:   /\.css$/,
  loader: ExtractTextPlugin.extract({fallback: 'style-loader', use: ['css-loader']}),
});

loaders.push({
  test:   /\.sass$/,
  loader: ExtractTextPlugin.extract({fallback: 'style-loader', use: common.sassLoaders}),
});


module.exports = {
  entry:         common.entry,
  output:        {
    path:              common.DIST_PATH,
    publicPath:        process.env.BASEPATH || '/',
    filename:          '[chunkhash].[name].js',
    sourceMapFilename: '[file].map',
  },
  module:        {
    rules: loaders,
  },
  devtool:       'eval',
  devServer:     _.extend(common.devServer, {hot: false}),
  resolveLoader: common.resolve.resolveLoaders,
  resolve:       common.resolve.resolve,
  externals:     {'jsdom': 'window'},
  node:          {
    fs:            'empty',
    child_process: 'empty',
    net:           'empty',
    tls:           'empty',
  },
  plugins:       [
    new ProgressBarPlugin(),
    new webpack.DefinePlugin(common.definePlugin),
    new webpack.ContextReplacementPlugin(/moment[\/\\]locale$/, /en/),
    new ExtractTextPlugin({filename: '[contenthash].app.css', ignoreOrder: true}),
    new webpack.LoaderOptionsPlugin({
      minimize: true,
      debug:    false,
    }),
    new webpack.optimize.UglifyJsPlugin({
      mangle:    false,
      compress:  true,
      sourceMap: false,
    }),
    new HtmlWebpackPlugin(_.extend({}, common.htmlWebpackPlugin, {
      hash:   true,
      minify: {},
      chunks: ['app'],
    })),
    new OptimizeCssAssetsPlugin(),
  ],
};