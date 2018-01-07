var path = require('path')

function cacheBusterQuery() {
  return `cb=${Math.round(Date.now() / 1000)}`
}

module.exports = function htmlWebpackPlugin(config, APP_PATH) {
  return {
    title:            config.title,
    chunks:           ['app'],
    cacheBusterQuery: cacheBusterQuery,
    favicon:          path.resolve(APP_PATH, 'images', 'favico.png'),
    template:         path.resolve(APP_PATH, 'index.html'),
    basePath:         '{{ .BasePath }}', // will be replaced by the server with the actual path
  }
}