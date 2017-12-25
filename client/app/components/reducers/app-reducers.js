import {createReducer} from 'redux-immutablejs'
import {fromJS} from 'immutable'
import {ACTIONS} from 'consts'

const INITIAL_STATE = fromJS({
  socket_ready: false,
  current_path: [],
  content:      [],
})

export function socketReady(state, action) {
  return state.set('socket_ready', true)
}

export function setCurrentPath(state, {payload}) {
  return state.set('current_path', payload)
}

export function setContent(state, {payload}) {
  return state.set('content', payload)
}

export const entities = createReducer(INITIAL_STATE, {
  [ACTIONS.SOCKET_READY]:     socketReady,
  [ACTIONS.SET_CURRENT_PATH]: setCurrentPath,
  [ACTIONS.SET_CONTENT]:      setContent,
})


export default entities