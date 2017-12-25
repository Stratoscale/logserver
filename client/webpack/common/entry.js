module.exports = function getEntry(APP_PATH, LOADER_PATH, UI_LOADER_PATH, PROGRESS_PATH) {
  return {
    'app': ['babel-polyfill', APP_PATH],
  }
};