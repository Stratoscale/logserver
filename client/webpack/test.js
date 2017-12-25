var path              = require('path');
var webpack           = require('webpack');
var ProgressBarPlugin = require('progress-bar-webpack-plugin');
var common            = require('./common')(true);

var loaders = common.loaders;

loaders.push(
  {
    test:    /\.js$/,
    include: [common.APP_PATH, common.LOADER_PATH, path.resolve(common.ROOT_PATH, 'test'), common.UI_LOADER_PATH, common.COMMON_PATH],
    exclude: /(node_modules)/,
    use:     [common.babelLoader, 'imports-loader?test=test-index'],
  });

loaders.push(
  {
    test:   /(\.css|\.sass)$/,
    loader: 'null-loader'
  });

var resolveTest = common.resolve

// this is used to mock MUI
resolveTest.alias = {
  app: path.resolve(common.ROOT_PATH, 'test')
};


module.exports = {
  devtool:       'eval',
  module:        {
    rules: loaders,
    noParse: [
      /node_modules\/sinon/,
    ],
  },
  resolveLoader: resolveTest.resolveLoaders,
  resolve:       resolveTest.resolve,
  plugins:       [
    new ProgressBarPlugin(),
    new webpack.DefinePlugin(common.definePlugin),
  ],
  externals:     {
    'cheerio':                        'window',
    'react/addons':                   true,
    'react/lib/ExecutionEnvironment': true,
    'react/lib/ReactContext':         true
  },
  node:          {
    fs:            'empty',
    child_process: 'empty',
    net:           'empty',
    tls:           'empty',
  },
};