import {List, Map} from 'immutable'
import {createSelector} from 'reselect'

export const appStateSelector = (state, props) => state.get('app', Map())

export const isSocketReady = createSelector(
  appStateSelector,
  (app = Map()) => app.get('socket_ready', false)
)

export const currentPathSelector = createSelector(
  appStateSelector,
  (app = Map()) => app.get('current_path')
)

const requestsSelector = createSelector(
  appStateSelector,
  (app = Map()) => app.get('requests')
)

export const hasPendingRequest = createSelector(
  requestsSelector,
  (requests = Map()) => requests.valueSeq().find(Boolean)
)

export const filesSelector = (state, props) => state.getIn(['files', 'tree'], List())
export const indexSelector = (state, props) => state.getIn(['files', 'index'], Map())

export const contentSelector = createSelector(
  appStateSelector,
  (app = Map()) => app.get('content')
)

export const fileSystemsSelector = createSelector(
  indexSelector,
  (index) => index.valueSeq().flatMap(file => file.get('instances', List()).map(instance => instance.get('fs'))).toSet()
)

export const filterSelector = createSelector(
  appStateSelector,
  (app = Map()) => app.get('filter')
)

export const searchSelector        = createSelector(
  appStateSelector,
  (app = Map()) => app.get('search')
)
export const searchResultsSelector = createSelector(
  appStateSelector,
  (app = Map()) => app.get('search_results')
)

export const locationSelect = (state = Map()) => state.get('router').location || {}