module.exports = function getDefinePlugin(config) {
  return {
    __DEV__: JSON.stringify(JSON.parse(process.env.DEV || 'false')),
    CONFIG:  JSON.stringify(config),
  }
}