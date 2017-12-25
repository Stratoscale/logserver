const fs    = require('fs')
const stats = require('../stats.json')

//print total time per module
const outputJson = stats.chunks[0].modules
  .map((module) =>
    ({
      name:    module.name,
      profile: module.profile,
      total:   (module.profile.factory + module.profile.building + (module.profile.dependencies || 0)) / 1000.0
    }))
  .sort((a, b) => (a.total > b.total) ? -1 : 1)

// output to csv

fs.writeFile('outputcsv.text', outputJson.map(item => `"${item.name}" ; ${item.profile.factory}; ${item.profile.building};${item.total}\n`))
