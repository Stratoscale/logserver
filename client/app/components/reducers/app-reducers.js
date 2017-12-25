import {createReducer} from 'redux-immutablejs'
import {fromJS, Map, List} from 'immutable'
import {ACTIONS} from 'consts'

const INITIAL_STATE = fromJS({
  socket_ready: false,
})

export function socketReady(state, action) {
  return state.set('socket_ready', true)
}

export const entities = createReducer(INITIAL_STATE, {
  [ACTIONS.SOCKET_READY]: socketReady,
})


export default entities