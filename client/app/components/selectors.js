import {List, Map} from 'immutable'
import {createSelector} from 'reselect'

export const appStateSelector = (state, props) => state.get('app', Map())

export const isSocketReady = createSelector(
  appStateSelector,
  (app = Map()) => app.get('socket_ready', false)
)

export const filesSelector = (state, props) => state.get('files', List())