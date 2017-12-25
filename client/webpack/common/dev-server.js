const path = require('path')

const devServerPort      = process.env.PORT || '8080'
const devServerHost      = process.env.HOST || 'localhost'
const isUiBackendLocally = process.env.UI_BACKEND === 'LOCAL'
const uiBackendServer    = isUiBackendLocally ? 'localhost:4172' : process.env.API_PROXY

module.exports = function getDevServer(ROOT_PATH, apiProxy) {
  const urlPrefix = isUiBackendLocally ? 'http://' : 'https://'
  return {
    overlay:            {
      warnings: false,
      errors:   true
    },
    historyApiFallback: true,
    inline:             true,
    hot:                true,
    contentBase:        path.resolve(ROOT_PATH, 'dist'),
    host:               devServerHost,
    port:               devServerPort,
    proxy:              {
      '/ws/*': {
        target: 'http://localhost:8888',
        secure: false,
        // pathRewrite: {
        //   '^/ui': isUiBackendLocally ? '' : '/ui',     // rewrite path
        // },
        ws:     true,
      },
      // '/ui/*': {
      //   target:      urlPrefix + uiBackendServer,
      //   secure:      false,
      //   pathRewrite: {
      //     '^/ui': isUiBackendLocally ? '' : '/ui',     // rewrite path
      //   },
      // },
    },
  }
}
