import {createReducer} from 'redux-immutablejs'
import {fromJS, Map, List} from 'immutable'
import {ACTIONS} from 'consts'

const INITIAL_STATE = fromJS({})

export function setFiles(state = Map(), {payload}) {
  return state.withMutations(state => {
    payload.forEach(_file => {
      const file = fromJS(_file)
      const path = file.get('path', List())
      state.updateIn(path.butLast(), (node = Map()) => {
        const filename = path.join('/')
        return node.mergeIn(['files', filename], file)
      })
    })
  })
}

export const files = createReducer(INITIAL_STATE, {
  [ACTIONS.SET_FILES]: setFiles,
})


export default files