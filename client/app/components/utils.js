import {store} from 'store/store'

export const runSelector = selector => selector(store.getState());