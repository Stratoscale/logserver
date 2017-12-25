var path = require('path');

function cacheBusterQuery() {
  return `cb=${Math.round(Date.now() / 1000)}`;
}

module.exports = function htmlWebpackPlugin(config, APP_PATH, CLUSTER_VERSION) {
  return {
    title: config.title,
    chunks: ['loader'],
    cacheBusterQuery: cacheBusterQuery,
    favicon: path.resolve(APP_PATH, 'images', 'favico.png'),
    template: path.resolve(APP_PATH, 'index.html'),
    cluster_version: CLUSTER_VERSION,
  }
};