import {createReducer} from 'redux-immutablejs'
import {fromJS, Map, List} from 'immutable'
import {ACTIONS} from 'consts'

const INITIAL_STATE = fromJS({
  tree:  {},
  index: {},
})

export function setFiles(state = Map(), {payload}) {
  if (!payload) {
    return state
  }
  return state.withMutations(state => {
    payload.forEach(_file => {
      const file = fromJS(_file)
      const path = file.get('path', List())
      state.updateIn(path.unshift('tree').butLast(), (node = Map()) => {
        return node.mergeIn(['files', file.get('key')], file)
      })
      state.setIn(['index', file.get('key')], file)
    })
  })
}

export const files = createReducer(INITIAL_STATE, {
  [ACTIONS.SET_FILES]: setFiles,
})


export default files