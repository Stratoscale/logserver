const jsonImporter = require('node-sass-json-importer');

module.exports = (ROOT_PATH, APP_PATH, COMMON_PATH, ENV) => {
  const sourceMap = Boolean(process.env.SOURCE_MAP)

  if (sourceMap) {
    console.log('sass/css source map activated mode')
  }


  let sassLoaders = [
    {
      loader:  "css-loader",
      options: {
        sourceMap,
      }
    },
    {
      loader:  'sass-loader',
      options: {
        sourceMap,
        indentedSyntax: true,
        importer:       jsonImporter,
      }
    }];

  // in prod we use ExtractTextPlugin which doesn't need style-loader
  if (ENV === 'development') {
    sassLoaders.unshift('style-loader')
  }

  return sassLoaders
};