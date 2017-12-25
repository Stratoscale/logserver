import {combineReducers} from 'redux-immutablejs';
import files from './file-reducers';
import {routerReducer} from 'react-router-redux'

export default combineReducers({
  files,
  router: routerReducer,
});
