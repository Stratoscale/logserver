import {createStore} from './create-store';

import reducers from 'reducers/reducers';

console.log('Creating store');

export const store = createStore(reducers);
if (__DEV__) {
  window.__store = store;
}

