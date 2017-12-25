import {combineReducers} from 'redux-immutablejs'
import files from './file-reducers'
import app from './app-reducers'
import {routerReducer} from 'react-router-redux'

export default combineReducers({
  app,
  files,
  router: routerReducer,
})
