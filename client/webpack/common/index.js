var execSync       = require('child_process').execSync;
var path           = require('path');
var ROOT_PATH      = path.resolve(__dirname, '..', '..');
var DIST_PATH      = path.resolve(ROOT_PATH, 'dist');
var APP_PATH       = path.resolve(ROOT_PATH, 'app');
var LOADER_PATH    = path.resolve(ROOT_PATH, 'loader');
var UI_LOADER_PATH = path.resolve(ROOT_PATH, 'ui-loader');
var PROGRESS_PATH  = path.resolve(ROOT_PATH, 'progress');
var COMMON_PATH    = path.resolve(ROOT_PATH, '..', 'common')

var env      = process.env.NODE_ENV || 'development';
var apiProxy = process.env.API_PROXY || 'localhost:8008/mock';

var clusterVersionConf = {
  'var':         '__CLUSTER_VERSION__',
  'placeholder': '{%CLUSTER_VERSION%}',
};

var uiHash = execSync('git rev-parse --verify HEAD --short').toString().trim();

var CLUSTER_VERSION = process.env.CLUSTER_VERSION;

var IS_REMOTE = process.env.REMOTE;


var REMOTE_URL = [process.env.REMOTE_URL || '//app.stratoscale.com', clusterVersionConf.placeholder].join('/');

var outputUrl = '/';

const ANALYTICS_DEV = process.env.ANALYTICS_DEV;

const analyticsKeys = require('./analytics')(ANALYTICS_DEV ? 'development' : env)


if (IS_REMOTE) {
  outputUrl = [REMOTE_URL.replace(clusterVersionConf.placeholder, CLUSTER_VERSION), ''].join('/');
}

const fontLoader = ['style-loader', 'css-loader', {loader: 'fontgen-loader', options: {types: 'ttf'}}]

const babelLoader = (isTest) => require('./babel')({TEST: isTest, DEV: env === 'development'});
var devServer     = require('./dev-server')(ROOT_PATH, apiProxy);
var loaders       = require('./loaders')(ROOT_PATH, COMMON_PATH, babelLoader, fontLoader);
var entry         = require('./entry')(APP_PATH, LOADER_PATH, UI_LOADER_PATH, PROGRESS_PATH);

var resolve           = require('./resolve')(ROOT_PATH);
var config            = require('../../config')(env);
var definePlugin      = require('./define-plugin')(config, clusterVersionConf, REMOTE_URL, uiHash, analyticsKeys);
var htmlWebpackPlugin = require('./html-webpack-plugin')(config, APP_PATH, CLUSTER_VERSION);
var sassLoaders       = require('./sass-loader')(ROOT_PATH, APP_PATH, COMMON_PATH, env)

module.exports = function getCommon(isTest = false) {
  return {
    ROOT_PATH:       ROOT_PATH,
    DIST_PATH:       DIST_PATH,
    APP_PATH:        APP_PATH,
    LOADER_PATH:     LOADER_PATH,
    UI_LOADER_PATH:  UI_LOADER_PATH,
    COMMON_PATH:     COMMON_PATH,
    ENV:             env,
    API_PROXY:       apiProxy,
    CLUSTER_VERSION: CLUSTER_VERSION,
    REMOTE_URL:      REMOTE_URL,

    outputUrl:          outputUrl,
    clusterVersionConf: clusterVersionConf,
    definePlugin:       definePlugin,
    htmlWebpackPlugin:  htmlWebpackPlugin,
    config:             config,
    devServer:          devServer,
    sassLoaders:        sassLoaders,
    loaders:            loaders,
    entry:              entry,
    resolve:            resolve,
    fontLoader:         fontLoader,
    babelLoader:        babelLoader(isTest),

  };
};