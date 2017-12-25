const Promise         = require('bluebird');
const fs              = Promise.promisifyAll(require('fs'));
const MESSAGES_FOLDER = '../app/messages';

const en = require('../app/messages/en.json')
const zh = require('../app/messages/zh_HANS.json')

const mapMissingKeys = (compareTo, messageFile) => {
  return Object.keys(compareTo).map((key) => !messageFile[key] ? key : undefined)
}

const reduceToObject = (acc, key) => {
  acc[key] = en[key];
  return acc
};

const readFile = (filename) => fs.readFileAsync(`${MESSAGES_FOLDER}/${filename}`, 'utf8');

const processFiles = (en) => {
  return fs.readdirAsync(MESSAGES_FOLDER).then(files => {
    files.filter(fileName => fileName.endsWith('.json') && fileName !== 'en.json')
      .forEach(fileName => {
        readFile(fileName).then(file => {
          const missingKeys = mapMissingKeys(en, JSON.parse(file)).reduce(reduceToObject, {});
          fs.writeFile(`missing-${fileName}`, JSON.stringify(missingKeys, null, 4))
        })
      })
  })
};

readFile('en.json').then(en => {
  processFiles(JSON.parse(en));
});

return;

const missingZHKeys = mapMissingKeys(en, zh).reduce(reduceToObject, {})
const missingENKeys = mapMissingKeys(zh, en).reduce(reduceToObject, {})

const errorHandler = (err) => console.log(err ? 'error' : 'file saved')

fs.writeFile('missingZH.json', JSON.stringify(missingZHKeys), errorHandler)
fs.writeFile('missingEN.json', JSON.stringify(missingENKeys), errorHandler)