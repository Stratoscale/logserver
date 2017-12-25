var _ = require('lodash');

var configs = {
  'development': require('./development'),
  'production': require('./production'),
  'common': require('./common')
};

module.exports = function getConfig(env) {
  return _.extend({}, configs.common, configs[env] || {});
};