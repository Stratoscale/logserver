var path = require('path');



module.exports = function (ROOT_PATH, COMMON_PATH, babelLoader, fontLoader) {
  return [
    {
      test:    /\.font\.js$/,
      use:     fontLoader,
      exclude: /(node_modules)/
    },
    {
      test:   /jquery|backbone|backbone\.marionette|backbone\.wreqr|backbone\.babysitter|plupload\.full\.min/,
      loader: 'imports-loader?this=>global,define=>false',
    },
    {
      test:    /\.js$/, // include .js files
      exclude: /node_modules|vendor/, // exclude any and all files in the node_modules folder
      use:     [
        {
          loader:  'eslint-loader',
          options: {
            failOnError: process.env.NODE_ENV === 'production',
            failOnWarning: process.env.NODE_ENV === 'production',
          }
        }
      ],
      enforce: 'pre',
    },
    {
      test:    /\.tpl$/,
      include: path.resolve(ROOT_PATH, 'app'),
      use:     ['tpl-loader'],
      exclude: /(node_modules)/
    },
    {
      test:    /\.coffee$/,
      include: [
        path.resolve(ROOT_PATH, 'app'),
        path.resolve(COMMON_PATH),
      ],
      use:     ['coffee-loader'],
      exclude: /(node_modules)/
    },
    {
      test:   /\.(woff|woff2)(\?.*)?$/,
      loader: 'url-loader?limit=10000&mimetype=application/font-woff',
    }, {
      test:   /\.ttf(\?.*)?$/,
      loader: 'url-loader?limit=10000&mimetype=application/octet-stream',
    }, {
      test:   /\.eot(\?.*)?$/,
      loader: 'url-loader?limit=10000&mimetype=application/vnd.ms-fontobject',
    }, {
      test:   /\.svg(\?.*)?$/,
      loader: 'url-loader?limit=10000&mimetype=image/svg+xml',
    },
    {
      test:    /\.(jpe?g|png|gif)$/i,
      loaders: [
        'file-loader?hash=sha512&digest=hex&name=[hash].[name].[ext]',
        'image-webpack-loader?bypassOnDebug&optimizationLevel=7&interlaced=false'
      ],
    },
  ];
}