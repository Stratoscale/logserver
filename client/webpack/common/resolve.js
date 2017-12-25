var path = require('path');

module.exports = function getResolve(ROOT_PATH) {
  return {
    resolve:        {
      modules:    [
        path.resolve(ROOT_PATH, 'app', 'components'),
        path.resolve(ROOT_PATH, 'app', 'messages'),
        path.resolve(ROOT_PATH, '..', 'common'),
        path.resolve(ROOT_PATH, 'node_modules'),
        'node_modules',
      ],
      extensions: ['.js', '.coffee', '.sass', '.json', '.css'],
      alias: {
        'require.resolve': 'resolve',
      },
    },
    resolveLoaders: {
      modules: [
        path.resolve(ROOT_PATH, 'node_modules'),
        'node_modules',
      ],
    },
  }
};
