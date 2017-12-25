const mock = () => true

const MUI = {
  vent: {on: mock, trigger: mock},
  addInitializer : mock,
  on: mock,
  module: mock,
  request: mock,
  execute: mock,
  entityManager: {setStore: mock},
  start: mock,
};

export default MUI