module.exports = function getEntry(APP_PATH) {
  return {
    'app': ['babel-polyfill', APP_PATH],
  }
}