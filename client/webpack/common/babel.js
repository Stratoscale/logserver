const chalk = require('chalk');

module.exports = ({TEST, DEV}) => {
  const presets = [
    require.resolve('babel-preset-react'),
    [require.resolve('babel-preset-env'), {
      targets:     {
        browsers: DEV && !TEST ? ['last 10 Chrome versions'] : [
          'last 5 versions',
          'not android <= 4.4.3',
          'not ie <= 11'
        ]
      },
      modules:     false,
      useBuiltIns: true,
      debug:       false,
    }]
  ];
  let plugins   = [
    [require.resolve('babel-plugin-transform-object-rest-spread'), {
      "useBuiltIns": true
    }],
    require.resolve('babel-plugin-transform-decorators-legacy'),
    require.resolve('babel-plugin-transform-class-properties'),
  ];

  if (TEST) {
    console.log('Babel in test mode');
    plugins.push(require.resolve('babel-plugin-rewire'));
  } else {
    plugins.push([require.resolve('babel-plugin-import'), {
      libraryName: 'antd',
    }])
  }

  if (DEV && !TEST) {
    plugins.push(require('react-hot-loader/babel'))
    console.log(chalk.yellow('Babel cache enabled'));
  }

  return {
    loader:  'babel-loader',
    options: {
      cacheDirectory: DEV ? true : false,
      presets,
      plugins,
      babelrc:        false,
    },
  }
}
;