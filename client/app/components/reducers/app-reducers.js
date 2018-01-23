import {createReducer} from 'redux-immutablejs'
import {List, fromJS} from 'immutable'
import {ACTIONS} from 'consts'

const INITIAL_STATE = fromJS({
  socket_ready:   false,
  index_ready:    false,
  current_path:   [],
  content:        [],
  filter:         null,
  search:         null,
  search_id:      null,
  content_id:     null,
  search_results: [],
  requests:       {},
})

function socketReady(state, action) {
  return state.set('socket_ready', true)
}

function indexReady(state, action) {
  return state.set('index_ready', true)
}

function setCurrentPath(state, {payload}) {
  return state.set('current_path', payload)
}

function setFilter(state, {payload}) {
  return state.set('filter', payload)
}

function setSearch(state, {payload, id}) {
  const result = state.set('search', payload)
  if (id) {
    return result.set('search_id', id)
  }
  return result
}

function setSearchId(state, {payload}) {
  return state.set('search_id', payload)
}

function setContentId(state, {payload}) {
  return state.set('content_id', payload)
}

function setSearchResults(state, {payload, meta: {id, add = false}}) {
  if (add && id === state.get('search_id')) {
    return state.set('search_results', state.get('search_results', List()).concat(fromJS(payload)))
  }
  return state.set('search_results', fromJS(payload)).set('search_id', id)
}

function setContent(state, {payload, meta: {id}}) {
  if (id === state.get('content_id')) {
    return state.set('content', state.get('content', List()).concat(fromJS(payload)))
  }
  return state.set('content', fromJS(payload)).set('content_id', id)
}

function clearContent(state) {
  return state.set('content', List())
}

function sendRequest(state, {payload: {action}}) {
  return state.setIn(['requests', action], true)
}

function receiveRequest(state, {payload: {action}}) {
  return state.setIn(['requests', action], false)
}

export const entities = createReducer(INITIAL_STATE, {
  [ACTIONS.SOCKET_READY]:       socketReady,
  [ACTIONS.INDEX_READY]:        indexReady,
  [ACTIONS.SET_CURRENT_PATH]:   setCurrentPath,
  [ACTIONS.SET_CONTENT]:        setContent,
  [ACTIONS.CLEAR_CONTENT]:      clearContent,
  [ACTIONS.SET_FILTER]:         setFilter,
  [ACTIONS.SET_SEARCH]:         setSearch,
  [ACTIONS.SET_SEARCH_ID]:      setSearchId,
  [ACTIONS.SET_CONTENT_ID]:     setContentId,
  [ACTIONS.SET_SEARCH_RESULTS]: setSearchResults,
  [ACTIONS.SEND_REQUEST]:       sendRequest,
  [ACTIONS.RECEIVE_REQUEST]:    receiveRequest,
})


export default entities