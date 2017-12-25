var conf = require('./base.conf')();

conf.browsers = ['Headless'];
conf.reporters = ['mocha'];
conf.singleRun = true;
module.exports = function (config) {
  config.set(conf);
};
