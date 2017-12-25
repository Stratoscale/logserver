import {createStore as reduxCreateStore, applyMiddleware} from 'redux';
import thunk from 'redux-thunk';
import {fromJS} from 'immutable';
import {composeWithDevTools} from 'redux-devtools-extension/developmentOnly';
import {browserHistory} from 'router/history'
import {routerMiddleware} from 'react-router-redux'

const INITIAL_STATE = fromJS({});

const enhancers = composeWithDevTools(
  applyMiddleware(thunk, routerMiddleware(browserHistory)),
);

export function createStore(reducers, initial_data = INITIAL_STATE) {
  return reduxCreateStore(reducers, initial_data, enhancers);
}