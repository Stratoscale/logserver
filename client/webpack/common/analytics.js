const chalk = require('chalk')

const analyticsKeys = {
  production:  {
    woopra:   'symphony-prod.stratoscale.com',
    mixpanel: '80ef9df7f08730702ff8a137fdd7e150',
  },
  development: {
    woopra:   'symphony-dev.stratoscale.com',
    mixpanel: '1c18da8af69ad881faeb667e0c3ab416',
  }
}

module.exports = (env) => {
  if (env === 'development') {
    console.log(
      chalk.yellow(`Analytics (${Object.keys(analyticsKeys[env]).join(',')}) in DEV mode`)
    );
  }
  return {...analyticsKeys[env], env}
}