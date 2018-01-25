import {Set, List, Map, is} from 'immutable'
import {createSelector, createSelectorCreator, defaultMemoize} from 'reselect'

const FIND_LIMIT = 5000

export const createImmutableSelector = createSelectorCreator(defaultMemoize, is)

export const appStateSelector = (state, props) => state.get('app', Map())

export const isSocketReady = createSelector(
  appStateSelector,
  (app = Map()) => app.get('socket_ready', false)
)

export const isIndexReady = createSelector(
  appStateSelector,
  (app = Map()) => app.get('index_ready', false)
)

export const currentPathSelector = createSelector(
  appStateSelector,
  (app = Map()) => app.get('current_path')
)

const requestsSelector = createSelector(
  appStateSelector,
  (app = Map()) => app.get('requests')
)

export const hasPendingRequest = type => createSelector(
  requestsSelector,
  (requests = Map()) => requests.get(type)
)

export const filesSelector = (state, props) => state.getIn(['files', 'tree'], List())
export const indexSelector = (state, props) => state.getIn(['files', 'index'], Map())

export const fileSystemsSelector = createSelector(
  indexSelector,
  (index) => index.valueSeq().flatMap(file => file.get('instances', List()).map(instance => instance.get('fs'))).toSet()
)

export const filterSelector = createSelector(
  appStateSelector,
  (app = Map()) => app.get('filter')
)

export const levelsSelector = createSelector(
  appStateSelector,
  (app = Map()) => app.get('levels')
)

export const contentSelector = createSelector(
  appStateSelector,
  levelsSelector,
  (app = Map(), levels = Set()) => app.get('content', List()).filter(line => !line.get('level') || levels.includes(line.get('level', '').toLowerCase()))
)

const findMatches = (content = List(), query) => {
  const matches = Map().asMutable()
  if (query) {
    for (let line of content) {
      const lineMatches = []
      const msg         = line.get('msg', '')
      let index         = msg.indexOf(query)
      while (~index) {
        lineMatches.push(index)
        const offset = index + query.length

        index = msg.substr(offset).indexOf(query)
        if (~index) {
          index += offset
        }
      }
      if (lineMatches.length) {
        matches.set(line.get('line'), lineMatches)
      }
      if (matches.size > FIND_LIMIT) {
        return matches.asImmutable()
      }
    }
  }
  return matches.asImmutable()
}

export const findSelector = createSelector(
  appStateSelector,
  (app = Map()) => app.get('find')
)

const querySelector = createSelector(
  findSelector,
  (find = Map()) => find.get('query')
)

export const matchesSelector = createImmutableSelector(
  querySelector,
  contentSelector,
  (query = '', content = List()) => {
    return findMatches(content, query)
  }
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