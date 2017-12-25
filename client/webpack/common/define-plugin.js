module.exports = function getDefinePlugin(config, clusterVersion, REMOTE_URL, uiHash, analyticsKeys) {
  return {
    __DEV__:                     JSON.stringify(JSON.parse(process.env.DEV || 'false')),
    CONFIG:                      JSON.stringify(config),
    CLUSTER_VERSION:             'window[\'' + clusterVersion.var + '\']',
    CLUSTER_VERSION_VAR:         clusterVersion.var,
    CLUSTER_VERSION_PLACEHOLDER: clusterVersion.placeholder,
    REMOTE_URL:                  '\'' + REMOTE_URL + '\'',
    'process.env.NODE_ENV':      JSON.stringify(process.env.DEV || 'production'),
    __UI_VERSION__:              {
      DATE: JSON.stringify(new Date().toISOString()),
      HASH: JSON.stringify(uiHash),
    },
    ANALYTICS_KEYS:              JSON.stringify(analyticsKeys)
  };
};