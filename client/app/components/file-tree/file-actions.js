import {fromJS} from 'immutable'
import {ACTIONS} from 'consts'

export const setCurrentPath = (path) => ({
  type:    ACTIONS.SET_CURRENT_PATH,
  payload: fromJS(path),
})