import {createReducer} from 'redux-immutablejs'
import {fromJS, Map, List} from 'immutable'
import {ACTIONS} from 'consts'

const INITIAL_STATE = fromJS([{name: 'a'}])

export function setFiles(state, action) {
  return state.set(action.payload)
}

export const entities = createReducer(INITIAL_STATE, {
  [ACTIONS.SET_FILES]: setFiles,
})


export default entities