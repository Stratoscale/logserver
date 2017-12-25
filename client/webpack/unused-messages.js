const Promise = require('bluebird');
const fs      = Promise.promisifyAll(require('fs'));
const grep1   = require('grep1');

const MESSAGES_FOLDER = '../app/messages';
const FOLDERS         = ['../app', '../../common'];

const readFile  = (filename) => fs.readFileAsync(`${MESSAGES_FOLDER}/${filename}`, 'utf8');
const writeFile = (filename, data) => fs.writeFileAsync(`${MESSAGES_FOLDER}/${filename}`, JSON.stringify(data, null, '  '), 'utf8');

const grep = (args) => new Promise((resolve, reject) => {
  grep1(args, (err, stdout, stderr) => {
    if (err || stderr) {
      reject(err, stderr)
    } else {
      resolve(stdout)
    }
  })
});

const processKey = key => {
  return new Promise((resolve, reject) => {
    grep([`[\'"]${key}[\'"]`, '-r', '--include=*.js', '--include=*.coffee', ...FOLDERS])
      .then(result => resolve(false))
      .catch((err, message) => {
        console.log('unused: ', key);
        resolve(key)
      });
  });
};

readFile('en.json').then(en => {
  const englishMessages = JSON.parse(en);
  const keys            = Object.keys(englishMessages);

  Promise.mapSeries(keys, processKey).then(results => {
    return fs.readdirAsync(MESSAGES_FOLDER).then(files => {
      files.filter(fileName => fileName.endsWith('.json'))
        .forEach(fileName => {
          readFile(fileName).then(file => {
            const messages = JSON.parse(file);
            const updated  = results.filter(Boolean).reduce((result, key) => {
              delete result[key];
              return result
            }, messages);
            writeFile(fileName, updated);
          })
        });
    });
  });
});
